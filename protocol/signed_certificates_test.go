package protocol

import (
	"testing"
)

var prefix = "schism-test/"

func TestCAPublicKeyS3Object_ObjectKey(t *testing.T) {
	sampleAuthKey := []byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN6gR4rRcthrCNDgBdOHhJQD/7bS+RTt/+BtUqAZGMEa")
	kFPrint := "SHA256:yCYTo2nP5zUcJuLWlHEJKj0jEElUE2wZvEMuh82UMQM"
	type fields struct {
		CertificateType    CertType
		HostCertAuthDomain string
		AuthorizedKey      []byte
		KeyFingerprint     string
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "User Public CA S3Object Key Without Fingerprint",
			fields: fields{CertificateType: UserCertificate},
			args:   args{prefix: prefix},
			want:   "schism-test/CA-Pubkeys/user.json",
		},
		{
			name: "User Public CA S3Object Key With Fingerprint",
			fields: fields{
				CertificateType: UserCertificate,
				KeyFingerprint:  kFPrint,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/user-yCYTo2nP5zUcJuLWlHEJKj0jEElUE2wZvEMuh82UMQM.json",
		},
		{
			name: "Host Public CA S3Object Key Without AuthDomain Without Fingerprint",
			fields: fields{
				CertificateType: HostCertificate,
				AuthorizedKey:   sampleAuthKey,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/host.json",
		},
		{
			name: "Host Public CA S3Object Key With AuthDomain Without Fingerprint",
			fields: fields{
				CertificateType:    HostCertificate,
				HostCertAuthDomain: "example.com",
				AuthorizedKey:      sampleAuthKey,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/host-example.com.json",
		},
		{
			name: "Host Public CA S3Object Key With AuthDomain With Fingerprint",
			fields: fields{
				CertificateType:    HostCertificate,
				HostCertAuthDomain: "example.com",
				AuthorizedKey:      sampleAuthKey,
				KeyFingerprint:     kFPrint,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/host-example.com-yCYTo2nP5zUcJuLWlHEJKj0jEElUE2wZvEMuh82UMQM.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CAPublicKeyS3Object{
				CertificateType:    tt.fields.CertificateType,
				HostCertAuthDomain: tt.fields.HostCertAuthDomain,
				AuthorizedKey:      tt.fields.AuthorizedKey,
				KeyFingerprint:     tt.fields.KeyFingerprint,
			}
			if got := c.ObjectKey(tt.args.prefix); got != tt.want {
				t.Errorf("ObjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignedCertificateS3Object_ObjectKey(t *testing.T) {
	type fields struct {
		CertificateType CertType
		Identity        string
		Principals      []string
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "HostSignedCertS3Object",
			fields: fields{
				CertificateType: HostCertificate,
				Identity:        "test.example.com",
				Principals:      []string{"test.example.com"},
			},
			args: args{prefix: prefix},
			want: "schism-test/host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json",
		},
		{
			name: "UserSignedCertS3Object",
			fields: fields{
				CertificateType: UserCertificate,
				Identity:        "user@test.example.com",
				Principals:      []string{"user", "admin"},
			},
			args: args{prefix: prefix},
			want: "schism-test/user:69206403b2f940935765c084335bcd2d9caed2fbd86a7056ddab98ce698e4ce1.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SignedCertificateS3Object{
				CertificateType: tt.fields.CertificateType,
				Identity:        tt.fields.Identity,
				Principals:      tt.fields.Principals,
			}
			if got := c.ObjectKey(tt.args.prefix); got != tt.want {
				t.Errorf("ObjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCertType_OppositeCA(t *testing.T) {
	tests := []struct {
		name string
		ct   CertType
		want CertType
	}{
		{
			name: "Host yields User",
			ct:   UserCertificate,
			want: HostCertificate,
		},
		{
			name: "User yields Host",
			ct:   HostCertificate,
			want: UserCertificate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.OppositeCA(); got != tt.want {
				t.Errorf("OppositeCA() = %v, want %v", got, tt.want)
			}
		})
	}
}
