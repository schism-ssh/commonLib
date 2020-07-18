package protocol

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// CertType: Schism supports two types of certificates: "user" and "host"
//
// Used to manage different aspects of the certification process
type CertType string

// Valid options for CertType
const (
	// Authenticates server hosts to users
	HostCertificate CertType = "host"
	// Authenticates users to servers
	UserCertificate CertType = "user"
	// Used in `schism list -ca` primarily but maybe I'll find more places this is useful
	// cakp ~> CertAuthority KeyPair
	CaKeyPair CertType = "cakp"
)

// OppositeCA returns "host" for "user" and "user" for "host"
func (ct CertType) OppositeCA() CertType {
	return (map[CertType]CertType{
		HostCertificate: UserCertificate,
		UserCertificate: HostCertificate,
	})[ct]
}

// Expand returns the full CertType given a single character type
//  h => host | u => user
//
// This is a no-op for already expanded `CertType`s
func (ct CertType) Expand() CertType {
	switch ct {
	case "h", HostCertificate:
		return HostCertificate
	case "u", UserCertificate:
		return UserCertificate
	case "c", CaKeyPair:
		return CaKeyPair
	default:
		return ""
	}
}

// S3Object provides an interface for saving mostly
// any struct to S3 as a Marshaled JSON object
type S3Object interface {
	// ObjectKey should take a prefix and
	// return a full object key for where the object should be saved
	ObjectKey(prefix string) string

	// LoadObject takes an S3 connection, S3 Bucket, and S3 ObjectKey
	// and populates the wantFields of the struct with the saved S3 object's data
	LoadObject(s3Svc s3iface.S3API, s3Bucket string, s3ObjectKey string) error
}

// rawLoadS3Object takes a bucket and key and returns the raw bytes.
// an error is returned if s3 has issues or if the object body cannot be read for some reason
func rawLoadS3Object(s3Svc s3iface.S3API, s3Bucket string, s3ObjectKey string) ([]byte, error) {
	object, err := s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3ObjectKey),
	})
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(object.Body)
	return body, err
}

// S3CaPubkeyPrefix The subprefix for storing the Public CA keys
//   Full Object path will follow this template
//    {profile.S3Prefix}{S3CaPubkeyPrefix}/{CertType}-{bundle_key_extras}.json
const S3CaPubkeyPrefix = "CA-Pubkeys/"

// S3CertStoragePrefix The subprefix for storing the Signed Certificate
//   Full Object path will follow this template
//    {profile.S3Prefix}{S3CertStoragePrefix}{LookupKey}.json
const S3CertStoragePrefix = "Signed-Certs/"

// SignedCertificateS3Object represents all the information
// that will be saved to S3 for a Signed SSH Certificate
type SignedCertificateS3Object struct {
	// Type of SSH-cert to be saved
	CertificateType CertType `json:"certificate_type"`
	// Timestamp of when we minted the cert
	IssuedOn time.Time `json:"issued_on"`
	// Originally requested Identity for this Certificate
	Identity string `json:"identity"`
	// Originally requested Principals for this Certificate
	Principals []string `json:"certificate_principals"`
	// How long will this Certificate be valid for?
	ValidityInterval time.Duration `json:"validity_interval"`
	// The raw representation of this Certificate after Marshaling
	// TODO: Work on KMS encryption for this section
	RawSignedCertificate []byte `json:"signed_certificate"`
	// The S3 ObjectKey for the AuthorizedKey half of the CA
	//   In theory this is here because when you have both halves working for Schism,
	//   Hosts need the Public half of the User CA to authenticate UserCertificates,
	//   and the reverse for the Users' side
	OppositePublicCA string `json:"opposite_public_ca"`
	// TODO: To be implemented later.
	SignedCertificateEncryption map[string]string `json:"signed_certificate_encryption,omitempty"`
}

// ObjectKey  given a prefix, return a key for S3 by invoking GenerateLookupKey()
// and calling .String() on the result
//
//  Format:
//   {prefix}{LookupKey.String()}.json
func (c *SignedCertificateS3Object) ObjectKey(prefix string) string {
	lookupKey := GenerateLookupKey(c.Identity, c.Principals, c.CertificateType)
	return fmt.Sprintf("%s%s%s.json", prefix, S3CertStoragePrefix, lookupKey)
}

// LoadObject loads an object from the given s3://{s3Bucket}/{s3ObjectKey} and un-marshals it into a SignedCertificateS3Object
func (c *SignedCertificateS3Object) LoadObject(s3Svc s3iface.S3API, s3Bucket string, s3ObjectKey string) error {
	body, err := rawLoadS3Object(s3Svc, s3Bucket, s3ObjectKey)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, c); err != nil {
		return fmt.Errorf("unable to unmarshal object (%s): %w", s3ObjectKey, err)
	}
	return nil
}

// CAPublicKeyS3Object represents all the information
// that will be saved to S3 for a given CA PublicKey
type CAPublicKeyS3Object struct {
	// Type of public key to be saved
	CertificateType CertType `json:"certificate_type"`
	// The raw representation of PublicKey after Marshaling to an AuthorizedKey format
	AuthorizedKey []byte `json:"authorized_key"`
	// The Fingerprint of the PublicKey as returned by ssh.FingerprintSHA256
	KeyFingerprint string `json:"fingerprint"`
	// If this is from a Host CA, the AuthDomain is the domain (or subdomain)
	// that the Host CA is authorized to sign certificates for
	//
	// In theory this can be a comma separated list but I haven't tested that yet
	HostCertAuthDomain string `json:"host_cert_auth_domain,omitempty"`
}

// ObjectKey, given a prefix, return a key for S3 based on CA type.
//
// If the HostCertAuthDomain is set, this will be added to the ObjectKey
//
//  Format:
//   {prefix}{S3CaPubkeyPrefix}/{host|user}.json
//   {prefix}{S3CaPubkeyPrefix}/{host|user}-{fingerprint}.json
//   {prefix}{S3CaPubkeyPrefix}/{host|user}-{HostCertAuthDomain}.json
//   {prefix}{S3CaPubkeyPrefix}/{host|user}-{HostCertAuthDomain}-{fingerprint}.json
func (c *CAPublicKeyS3Object) ObjectKey(prefix string) string {
	var subKey string
	if c.HostCertAuthDomain != "" {
		subKey = fmt.Sprintf("%s-%s", c.CertificateType, c.HostCertAuthDomain)
	} else {
		subKey = string(c.CertificateType)
	}
	if c.KeyFingerprint != "" {
		// KeyFingerprint is in format "SHA256:{fingerprint}"
		fingerprint := strings.Split(c.KeyFingerprint, ":")[1]
		subKey = fmt.Sprintf("%s-%s", subKey, fingerprint)
	}
	return fmt.Sprintf("%s%s%s.json", prefix, S3CaPubkeyPrefix, subKey)
}

// LoadObject loads an object from the given s3://{s3Bucket}/{s3ObjectKey} and un-marshals it into a CAPublicKeyS3Object
func (c *CAPublicKeyS3Object) LoadObject(s3Svc s3iface.S3API, s3Bucket string, s3ObjectKey string) error {
	body, err := rawLoadS3Object(s3Svc, s3Bucket, s3ObjectKey)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, c); err != nil {
		return fmt.Errorf("unable to unmarshal object (%s): %w", s3ObjectKey, err)
	}
	return nil
}
