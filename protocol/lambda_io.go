package protocol

import "time"

// RequestSSHCertLambdaPayload is used to pass the required information to the lambda function
type RequestSSHCertLambdaPayload struct {
	// Type of SSH-cert being requested.
	//
	// The following are accepted:
	//    * "host"
	//    * "user
	CertificateType CertType `json:"certificate_type"`
	// Specify the key identity when signing a public key.
	Identity string `json:"certificate_identity"`
	// Specify one or more principals (user or host names) to be included in a certificate when signing a key.
	Principals []string `json:"certificate_principals"`
	// Length of time the Signed Certificate will be valid for.
	ValidityInterval time.Duration `json:"validity_interval"`
	// Specify a certificate option when signing a key.
	UserKeyOptions []string `json:"user_key_options,omitempty"`
	// Public Key to submit to the CA for signing.
	// Supported types: (Others may work, ymmv)
	//
	//    * "ed25519"
	//    * "rsa"
	PublicKey string `json:"public_key"`
}

// RequestSSHCertLambdaResponse is used to return pertinent information from the lambda function
type RequestSSHCertLambdaResponse struct {
	// Type of SSH-cert that was generated.
	CertificateType CertType `json:"certificate_type"`
	// 64-character key used with the CertType to fetch the Certificate bundle from S3
	LookupKey string `json:"lookup_key"`
}
