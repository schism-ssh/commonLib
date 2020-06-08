package protocol

import (
	"fmt"
	"time"
)

// CertType: Schism supports two types of certificates: "user" and "host".
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

// RequestSSHCertLambdaPayload is used to pass the required information to the lambda function
type RequestSSHCertLambdaPayload struct {
	// Type of SSH-cert being requested.
	//
	// The following are accepted:
	//    * "host"
	//    * "user
	CertificateType  CertType      `json:"certificate_type"`
	// Specify the key identity when signing a public key.
	Identity         string        `json:"certificate_identity"`
	// Specify one or more principals (user or host names) to be included in a certificate when signing a key.
	Principals       []string      `json:"certificate_principals"`
	// Length of time the Signed Certificate will be valid for.
	ValidityInterval time.Duration `json:"validity_interval"`
	// Specify a certificate option when signing a key.
	UserKeyOptions   []string      `json:"user_key_options,omitempty"`
	// Public Key to submit to the CA for signing.
	// Supported types: (Others may work, ymmv)
	//
	//    * "ed25519"
	//    * "rsa"
	PublicKey        string        `json:"public_key"`
}

type RequestSSHCertLambdaResponse struct {
	LookupKey string `json:"lookup_key"`
}

type S3Object interface {
	ObjectKey(prefix string) string
}

type SignedCertificateS3Object struct {
	CertificateType             CertType          `json:"certificate_type"`
	IssuedOn                    time.Time         `json:"issued_on"`
	Identity                    string            `json:"identity"`
	Principals                  []string          `json:"certificate_principals"`
	ValidityInterval            time.Duration     `json:"validity_interval"`
	RawSignedCertificate        []byte            `json:"signed_certificate"`
	OppositePublicCA            string            `json:"opposite_public_ca"`
	SignedCertificateEncryption map[string]string `json:"signed_certificate_encryption,omitempty"`
}

func (c *SignedCertificateS3Object) ObjectKey(prefix string) string {
	lookupKey := GenerateLookupKey(c.Identity, c.Principals)
	return fmt.Sprintf("%s%ss/%s.json", prefix, c.CertificateType, lookupKey)
}

type CAPublicKeyS3Object struct {
	CertificateType    CertType `json:"certificate_type"`
	AuthorizedKey      []byte   `json:"authorized_key"`
	HostCertAuthDomain string   `json:"host_cert_auth_domain,omitempty"`
}

func (c *CAPublicKeyS3Object) ObjectKey(prefix string) string {
	objectPrefix := "CA-Certs/"
	var subKey string
	if c.HostCertAuthDomain != "" {
		subKey = fmt.Sprintf("%s-%s", c.HostCertAuthDomain, c.CertificateType)
	} else {
		subKey = string(c.CertificateType)
	}
	return fmt.Sprintf("%s%s%s.json", prefix, objectPrefix, subKey)
}
