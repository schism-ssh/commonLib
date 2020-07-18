package protocol_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3/s3iface"

	"src.doom.fm/schism/commonLib/protocol"
)

var (
	hostTestExampleComKey = "55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"
	userSomeUserDevKey    = "a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64"
)

func TestGenerateLookupKey(t *testing.T) {
	type args struct {
		ident      string
		principals []string
		certType   protocol.CertType
	}
	tests := []struct {
		name string
		args args
		want *protocol.LookupKey
	}{
		{
			name: "valid LookupKey for a host",
			args: args{
				ident:      "test.example.com",
				principals: []string{"test.example.com"},
				certType:   protocol.HostCertificate,
			},
			want: &protocol.LookupKey{
				Id:   hostTestExampleComKey,
				Type: protocol.HostCertificate,
			},
		},
		{
			name: "valid LookupKey for a user",
			args: args{
				ident:      "someUser@dev1.example.com",
				principals: []string{"someUser", "admin"},
				certType:   protocol.UserCertificate,
			},
			want: &protocol.LookupKey{
				Id:   userSomeUserDevKey,
				Type: protocol.UserCertificate,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := protocol.GenerateLookupKey(tt.args.ident, tt.args.principals, tt.args.certType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateLookupKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupKey_Expand(t *testing.T) {
	type fields struct {
		Id   string
		Type protocol.CertType
	}
	type args struct {
		s3Svc    s3iface.S3API
		s3Bucket string
		s3Prefix string
	}
	validBucketArgs := args{
		s3Svc:    &protocol.MockS3Client{T: t},
		s3Bucket: protocol.TestValidBucket,
	}
	validHostLookupKey := &protocol.LookupKey{
		Id:   hostTestExampleComKey,
		Type: protocol.HostCertificate,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *protocol.LookupKey
		wantErr bool
	}{
		{
			name: "valid lookup key that returns a single result",
			fields: fields{
				Id:   "55e8182ec4413d51",
				Type: protocol.HostCertificate,
			},
			args:    validBucketArgs,
			want:    validHostLookupKey,
			wantErr: false,
		},
		{
			name: "valid lookup key with short key and type that returns a single result",
			fields: fields{
				Id:   "55e8182ec4413d51",
				Type: "h",
			},
			args:    validBucketArgs,
			want:    validHostLookupKey,
			wantErr: false,
		},
		{
			name: "valid lookup key that returns multiple keys",
			fields: fields{
				Id:   "4",
				Type: protocol.UserCertificate,
			},
			args:    validBucketArgs,
			wantErr: true,
		},
		{
			name: "invalid lookup key that returns zero matches",
			fields: fields{
				Id:   "0f739d75b44acc5b",
				Type: protocol.UserCertificate,
			},
			args:    validBucketArgs,
			wantErr: true,
		},
		{
			name: "non-existent bucket returns err from aws",
			fields: fields{
				Id:   "0f739d75b44acc5b",
				Type: protocol.UserCertificate,
			},
			args: args{
				s3Svc:    &protocol.MockS3Client{T: t},
				s3Bucket: "this-bucket-is-a-lie",
			},
			wantErr: true,
		},
		{
			name: "expanded key is in invalid format",
			fields: fields{
				Id:   "d0c671a71f190313",
				Type: protocol.HostCertificate,
			},
			args: args{
				s3Svc:    &protocol.MockS3Client{T: t},
				s3Bucket: protocol.TestValidBucket,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lk := &protocol.LookupKey{
				Id:   tt.fields.Id,
				Type: tt.fields.Type,
			}
			err := lk.Expand(tt.args.s3Svc, tt.args.s3Bucket, tt.args.s3Prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(lk, tt.want) {
					t.Errorf("Expand() got = %v, want = %v", lk, tt.want)
				}
			}
		})
	}
}

func TestLookupKey_MarshalJSON(t *testing.T) {
	type fields struct {
		Id   string
		Type protocol.CertType
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Host LookupKey",
			fields: fields{
				Id:   hostTestExampleComKey,
				Type: protocol.HostCertificate,
			},
			want: []byte("host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"),
		},
		{
			name: "User LookupKey",
			fields: fields{
				Id:   userSomeUserDevKey,
				Type: protocol.UserCertificate,
			},
			want: []byte("user:a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lk := &protocol.LookupKey{
				Id:   tt.fields.Id,
				Type: tt.fields.Type,
			}
			got, err := lk.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupKey_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *protocol.LookupKey
		wantErr bool
	}{
		{
			name: "Raw User LookupKey",
			args: args{
				data: []byte("user:a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64"),
			},
			want: &protocol.LookupKey{
				Id:   userSomeUserDevKey,
				Type: protocol.UserCertificate,
			},
			wantErr: false,
		},
		{
			name: "Invalid Host LookupKey",
			args: args{
				data: []byte("hosts/55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"),
			},
			want:    &protocol.LookupKey{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lk := &protocol.LookupKey{}
			if err := lk.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(lk, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", lk, tt.want)
			}
		})
	}
}

func TestLookupKey_String(t *testing.T) {
	type fields struct {
		Id   string
		Type protocol.CertType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Host LookupKey",
			fields: fields{
				Id:   hostTestExampleComKey,
				Type: protocol.HostCertificate,
			},
			want: "host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d",
		},
		{
			name: "User LookupKey",
			fields: fields{
				Id:   userSomeUserDevKey,
				Type: protocol.UserCertificate,
			},
			want: "user:a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lk := &protocol.LookupKey{
				Id:   tt.fields.Id,
				Type: tt.fields.Type,
			}
			if got := lk.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLookupKey(t *testing.T) {
	type args struct {
		rawKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *protocol.LookupKey
		wantErr bool
	}{
		{
			name: "correctly parses host LookupKey",
			args: args{rawKey: "host:55e8182ec4413d51"},
			want: &protocol.LookupKey{
				Id:   "55e8182ec4413d51",
				Type: protocol.HostCertificate,
			},
			wantErr: false,
		},
		{
			name: "correctly parses user LookupKey",
			args: args{rawKey: "user:a5ba427b532c152b"},
			want: &protocol.LookupKey{
				Id:   "a5ba427b532c152b",
				Type: protocol.UserCertificate,
			},
			wantErr: false,
		},
		{
			name: "correctly parses short LookupKey",
			args: args{rawKey: "h:55e818"},
			want: &protocol.LookupKey{
				Id:   "55e818",
				Type: "h",
			},
			wantErr: false,
		},
		{
			name:    "returns an error if the Key is invalid",
			args:    args{rawKey: "hosts/55e8182ec4413d51"},
			want:    &protocol.LookupKey{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := protocol.ParseLookupKey(tt.args.rawKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLookupKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseLookupKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
