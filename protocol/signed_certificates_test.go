package protocol_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"code.agarg.me/schism/commonLib/protocol"
)

var prefix = "schism-test/"

func TestCAPublicKeyS3Object_ObjectKey(t *testing.T) {
	sampleAuthKey := []byte("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN6gR4rRcthrCNDgBdOHhJQD/7bS+RTt/+BtUqAZGMEa")
	kFPrint := "SHA256:yCYTo2nP5zUcJuLWlHEJKj0jEElUE2wZvEMuh82UMQM"
	type fields struct {
		CertificateType    protocol.CertType
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
			fields: fields{CertificateType: protocol.UserCertificate},
			args:   args{prefix: prefix},
			want:   "schism-test/CA-Pubkeys/user.json",
		},
		{
			name: "User Public CA S3Object Key With Fingerprint",
			fields: fields{
				CertificateType: protocol.UserCertificate,
				KeyFingerprint:  kFPrint,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/user-yCYTo2nP5zUcJuLWlHEJKj0jEElUE2wZvEMuh82UMQM.json",
		},
		{
			name: "Host Public CA S3Object Key Without AuthDomain Without Fingerprint",
			fields: fields{
				CertificateType: protocol.HostCertificate,
				AuthorizedKey:   sampleAuthKey,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/host.json",
		},
		{
			name: "Host Public CA S3Object Key With AuthDomain Without Fingerprint",
			fields: fields{
				CertificateType:    protocol.HostCertificate,
				HostCertAuthDomain: "example.com",
				AuthorizedKey:      sampleAuthKey,
			},
			args: args{prefix: prefix},
			want: "schism-test/CA-Pubkeys/host-example.com.json",
		},
		{
			name: "Host Public CA S3Object Key With AuthDomain With Fingerprint",
			fields: fields{
				CertificateType:    protocol.HostCertificate,
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
			c := &protocol.CAPublicKeyS3Object{
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

func TestCAPublicKeyS3Object_LoadObject(t *testing.T) {
	type fields struct {
		CertificateType    protocol.CertType
		AuthorizedKey      []byte
		KeyFingerprint     string
		HostCertAuthDomain string
	}
	type args struct {
		s3Svc       s3iface.S3API
		s3Bucket    string
		s3ObjectKey string
	}
	tests := []struct {
		name       string
		wantFields fields
		args       args
		wantErr    bool
	}{
		{
			name: "Loads a ca public key from S3",
			wantFields: fields{
				CertificateType: protocol.UserCertificate,
				KeyFingerprint:  "SHA256:Gyc2MeVs5jZsL2lnDQj8C0FA6qOdZavwl+aY6APh7TM",
			},
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: protocol.S3CaPubkeyPrefix + "user.json",
			},
			wantErr: false,
		},
		{
			name: "Handles errors from s3 gracefully",
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: "invalid-key",
			},
			wantErr: true,
		},
		{
			name: "Handles corrupt ca pubkeys gracefully",
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: "empty-objects",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &protocol.CAPublicKeyS3Object{}
			want := &protocol.CAPublicKeyS3Object{
				CertificateType:    tt.wantFields.CertificateType,
				AuthorizedKey:      tt.wantFields.AuthorizedKey,
				KeyFingerprint:     tt.wantFields.KeyFingerprint,
				HostCertAuthDomain: tt.wantFields.HostCertAuthDomain,
			}
			if err := c.LoadObject(tt.args.s3Svc, tt.args.s3Bucket, tt.args.s3ObjectKey); (err != nil) != tt.wantErr {
				t.Errorf("LoadObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c, want) {
				t.Errorf("LoadObject() got = %+v, want %+v", c, want)
			}
		})
	}
}

func TestSignedCertificateS3Object_ObjectKey(t *testing.T) {
	type fields struct {
		CertificateType protocol.CertType
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
				CertificateType: protocol.HostCertificate,
				Identity:        "test.example.com",
				Principals:      []string{"test.example.com"},
			},
			args: args{prefix: prefix},
			want: "schism-test/Signed-Certs/host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json",
		},
		{
			name: "UserSignedCertS3Object",
			fields: fields{
				CertificateType: protocol.UserCertificate,
				Identity:        "user@test.example.com",
				Principals:      []string{"user", "admin"},
			},
			args: args{prefix: prefix},
			want: "schism-test/Signed-Certs/user:69206403b2f940935765c084335bcd2d9caed2fbd86a7056ddab98ce698e4ce1.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &protocol.SignedCertificateS3Object{
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

func TestSignedCertificateS3Object_LoadObject(t *testing.T) {
	type fields struct {
		CertificateType             protocol.CertType
		IssuedOn                    time.Time
		Identity                    string
		Principals                  []string
		ValidityInterval            time.Duration
		RawSignedCertificate        []byte
		OppositePublicCA            string
		SignedCertificateEncryption map[string]string
	}
	type args struct {
		s3Svc       s3iface.S3API
		s3Bucket    string
		s3ObjectKey string
	}
	tests := []struct {
		name       string
		wantFields fields
		args       args
		wantErr    bool
	}{
		{
			name: "Loads Signed Certificate from S3",
			wantFields: fields{
				CertificateType:  protocol.HostCertificate,
				Identity:         "test.example.com",
				Principals:       []string{"test.example.com"},
				ValidityInterval: 120 * time.Hour,
			},
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: protocol.S3CertStoragePrefix + "host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json",
			},
			wantErr: false,
		},
		{
			name: "Handles errors from S3 gracefully",
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: "invalid-key",
			},
			wantErr: true,
		},
		{
			name: "Handles corrupt stored certs gracefully",
			args: args{
				s3Svc:       &protocol.MockS3Client{T: t},
				s3Bucket:    protocol.TestValidBucket,
				s3ObjectKey: "empty-objects",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &protocol.SignedCertificateS3Object{}
			want := &protocol.SignedCertificateS3Object{
				CertificateType:             tt.wantFields.CertificateType,
				IssuedOn:                    tt.wantFields.IssuedOn,
				Identity:                    tt.wantFields.Identity,
				Principals:                  tt.wantFields.Principals,
				ValidityInterval:            tt.wantFields.ValidityInterval,
				RawSignedCertificate:        tt.wantFields.RawSignedCertificate,
				OppositePublicCA:            tt.wantFields.OppositePublicCA,
				SignedCertificateEncryption: tt.wantFields.SignedCertificateEncryption,
			}
			if err := c.LoadObject(tt.args.s3Svc, tt.args.s3Bucket, tt.args.s3ObjectKey); (err != nil) != tt.wantErr {
				t.Errorf("LoadObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(c, want) {
				t.Errorf("LoadObject() got = %+v, want %+v", c, want)
			}
		})
	}
}

func TestCertType_OppositeCA(t *testing.T) {
	tests := []struct {
		name string
		ct   protocol.CertType
		want protocol.CertType
	}{
		{
			name: "Host yields User",
			ct:   protocol.UserCertificate,
			want: protocol.HostCertificate,
		},
		{
			name: "User yields Host",
			ct:   protocol.HostCertificate,
			want: protocol.UserCertificate,
		},
		{
			name: "Invalid Type yields empty string",
			ct:   "admin",
			want: "",
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

func TestCertType_Expand(t *testing.T) {
	tests := []struct {
		name string
		ct   protocol.CertType
		want protocol.CertType
	}{
		{
			name: "h yields host",
			ct:   "h", want: "host",
		},
		{
			name: "u yields user", ct: "u", want: "user"},
		{
			name: "c yields cakp",
			ct:   "c", want: "cakp",
		},
		{
			name: "expanded type yields itself",
			ct:   "host", want: "host",
		},
		{
			name: "invalid type yields empty string",
			ct:   "f", want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.Expand(); got != tt.want {
				t.Errorf("Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}
