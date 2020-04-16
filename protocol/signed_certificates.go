package protocol

type RequestSignedCertificateLambdaPayload struct {
	CertificateType  string   `json:"certificate_type"`
	Identity         string   `json:"certificate_identity"`
	Principals       []string `json:"certificate_principals"`
	ValidityInterval string   `json:"validity_interval"`
	PublicKey        string   `json:"public_key"`
}

type RequestSignedCertificateLambdaResponse struct {
	LookupKey string `json:"lookup_key"`
}

type SignedCertificateS3Object struct {
	CertificateType             string            `json:"certificate_type"`
	IssuedOn                    string            `json:"issued_on"`
	Identity                    string            `json:"identity"`
	Principals                  []string          `json:"certificate_principals"`
	ValidityInterval            string            `json:"validity_interval"`
	RawSignedCertificate        []string          `json:"signed_certificate"`
	SignedCertificateEncryption map[string]string `json:"signed_certificate_encryption"`
}
