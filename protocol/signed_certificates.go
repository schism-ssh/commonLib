package protocol

import (
	"fmt"
	"time"
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
)

// OppositeCA returns "host" for "user" and "user" for "host"
func (ct CertType) OppositeCA() CertType {
	return (map[CertType]CertType{
		HostCertificate: UserCertificate,
		UserCertificate: HostCertificate,
	})[ct]
}

// S3Object provides an interface for saving mostly
// any struct to S3 as a Marshaled JSON object
type S3Object interface {
	// ObjectKey should take a prefix and
	// return a full object key for where the object should be saved
	ObjectKey(prefix string) string
}

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

// ObjectKey, given a prefix, return a key for S3 by invoking GenerateLookupKey()
// and adding an `s' to the CertType
//
//  Format:
//   {prefix}{CertType}s/{LookupKey}.json
func (c *SignedCertificateS3Object) ObjectKey(prefix string) string {
	lookupKey := GenerateLookupKey(c.Identity, c.Principals)
	return fmt.Sprintf("%s%ss/%s.json", prefix, c.CertificateType, lookupKey)
}

// CAPublicKeyS3Object represents all the information
// that will be saved to S3 for a given CA PublicKey
type CAPublicKeyS3Object struct {
	// Type of public key to be saved
	CertificateType CertType `json:"certificate_type"`
	// The raw representation of PublicKey after Marshaling to an AuthorizedKey format
	AuthorizedKey []byte `json:"authorized_key"`
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
//   {prefix}CA-Certs/{host|user}.json
//   {prefix}CA-Certs/{host|user}-{HostCertAuthDomain}.json
func (c *CAPublicKeyS3Object) ObjectKey(prefix string) string {
	objectPrefix := "CA-Certs/"
	var subKey string
	if c.HostCertAuthDomain != "" {
		subKey = fmt.Sprintf("%s-%s", c.CertificateType, c.HostCertAuthDomain)
	} else {
		subKey = string(c.CertificateType)
	}
	return fmt.Sprintf("%s%s%s.json", prefix, objectPrefix, subKey)
}
