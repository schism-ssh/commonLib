package protocol

import "testing"

var prefix = "schism-test/"

func TestCAPublicKeyS3Object_ObjectKey(t *testing.T) {
	type fields struct {
		CertificateType    string
		HostCertAuthDomain string
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
			name: "UserPublicCAS3ObjectKey",
			fields: fields{
				CertificateType:    UserCertificate,
				HostCertAuthDomain: "*.example.com",
			},
			args: args{
				prefix: prefix,
			},
			want: "schism-test/CA-Certs/*.example.com-user.json",
		},
		{
			name: "HostPublicCAS3ObjectKey",
			fields: fields{
				CertificateType: HostCertificate,
			},
			args: args{
				prefix: prefix,
			},
			want: "schism-test/CA-Certs/host.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CAPublicKeyS3Object{
				CertificateType:    tt.fields.CertificateType,
				HostCertAuthDomain: tt.fields.HostCertAuthDomain,
			}
			if got := c.ObjectKey(tt.args.prefix); got != tt.want {
				t.Errorf("ObjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignedCertificateS3Object_ObjectKey(t *testing.T) {
	type fields struct {
		CertificateType             string
		Identity                    string
		Principals                  []string
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
			args: args{
				prefix: prefix,
			},
			want: "schism-test/hosts/55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json",
		},
		{
			name: "UserSignedCertS3Object",
			fields: fields{
				CertificateType: UserCertificate,
				Identity:        "user@test.example.com",
				Principals:      []string{"user", "admin"},
			},
			args: args{
				prefix: prefix,
			},
			want: "schism-test/users/69206403b2f940935765c084335bcd2d9caed2fbd86a7056ddab98ce698e4ce1.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &SignedCertificateS3Object{
				CertificateType:             tt.fields.CertificateType,
				Identity:                    tt.fields.Identity,
				Principals:                  tt.fields.Principals,
			}
			if got := c.ObjectKey(tt.args.prefix); got != tt.want {
				t.Errorf("ObjectKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
