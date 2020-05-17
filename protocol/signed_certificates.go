package protocol

import (
	"fmt"
	"time"
)

const (
	HostCertificate = "host"
	UserCertificate = "user"
)

type RequestSSHCertLambdaPayload struct {
	CertificateType  string        `json:"certificate_type"`
	Identity         string        `json:"certificate_identity"`
	Principals       []string      `json:"certificate_principals"`
	ValidityInterval time.Duration `json:"validity_interval"`
	UserKeyOptions   []string      `json:"user_key_options,omitempty"`
	PublicKey        string        `json:"public_key"`
}

type RequestSSHCertLambdaResponse struct {
	LookupKey string `json:"lookup_key"`
}

type S3Object interface {
	ObjectKey(prefix string) string
}

type SignedCertificateS3Object struct {
	CertificateType             string            `json:"certificate_type"`
	IssuedOn                    string            `json:"issued_on"`
	Identity                    string            `json:"identity"`
	Principals                  []string          `json:"certificate_principals"`
	ValidityInterval            string            `json:"validity_interval"`
	RawSignedCertificate        string            `json:"signed_certificate"`
	OppositePublicCA            string            `json:"opposite_public_ca"`
	SignedCertificateEncryption map[string]string `json:"signed_certificate_encryption,omitempty"`
}

func (c *SignedCertificateS3Object) ObjectKey(prefix string) string {
	lookupKey := GenerateLookupKey(c.Identity, c.Principals)
	return fmt.Sprintf("%s%ss/%s.json", prefix, c.CertificateType, lookupKey)
}

type CAPublicKeyS3Object struct {
	CertificateType    string `json:"certificate_type"`
	AuthorizedKey      []byte `json:"authorized_key"`
	HostCertAuthDomain string `json:"host_cert_auth_domain,omitempty"`
}

func (c *CAPublicKeyS3Object) ObjectKey(prefix string) string {
	objectPrefix := "CA-Certs/"
	var subKey string
	if c.HostCertAuthDomain != "" {
		subKey = fmt.Sprintf("%s-%s", c.HostCertAuthDomain, c.CertificateType)
	} else {
		subKey = c.CertificateType
	}
	return fmt.Sprintf("%s%s%s.json", prefix, objectPrefix, subKey)
}
