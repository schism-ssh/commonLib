package protocol

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var validBucket = "schism-test"

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	output := &s3.ListObjectsV2Output{}
	var contents []*s3.Object
	switch *input.Prefix {
	case "hosts/55e8182ec4413d51":
		contents = []*s3.Object{{Key: aws.String("hosts/55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json")}}
	case "users/4":
		contents = []*s3.Object{
			{Key: aws.String("users/4e1586bed08190ccac4056078afed44daac058e8361b216dd078c7714b874cae.json")},
			{Key: aws.String("users/4d5b5d59343254c4fccafe48813ceeb99ae5ce44c1b97113b370a93f8411a01e.json")},
		}
	default:
		// No matches
	}

	if *input.Bucket != validBucket {
		// return error
		return output, fmt.Errorf("NoSuchBucket")
	}
	output.SetContents(contents)
	output.SetKeyCount(int64(len(contents)))
	return output, nil
}

func TestGenerateLookupKey(t *testing.T) {
	type args struct {
		ident      string
		principals []string
	}
	tests := []struct {
		name string
		args args
		want LookupKey
	}{
		{
			name: "returns a lookup key for a host",
			args: args{
				ident:      "test.example.com",
				principals: []string{"test.example.com"},
			},
			want: LookupKey("55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"),
		},
		{
			name: "returns a lookup key for a user",
			args: args{
				ident:      "someUser@dev1.example.com",
				principals: []string{"someUser", "admin"},
			},
			want: LookupKey("a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateLookupKey(tt.args.ident, tt.args.principals); got != tt.want {
				t.Errorf("GenerateLookupKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupKey_Expand(t *testing.T) {
	type args struct {
		s3Svc    s3iface.S3API
		s3Bucket string
		s3Prefix string
		certType CertType
	}

	hostCertArgs := args{
		s3Svc:    &mockS3Client{},
		s3Bucket: validBucket,
		certType: HostCertificate,
	}
	userCertArgs := args{
		s3Svc:    &mockS3Client{},
		s3Bucket: validBucket,
		certType: UserCertificate,
	}
	tests := []struct {
		name    string
		lk      LookupKey
		args    args
		want    LookupKey
		wantErr bool
	}{
		{
			name:    "valid lookup key that returns a single result",
			lk:      "55e8182ec4413d51",
			args:    hostCertArgs,
			want:    "55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d",
			wantErr: false,
		},
		{
			name:    "valid lookup key that returns multiple keys",
			lk:      "4",
			args:    userCertArgs,
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid lookup key that returns zero matches",
			lk:      "0f739d75b44acc5b",
			args:    userCertArgs,
			want:    "",
			wantErr: true,
		},
		{
			name: "non-existent bucket returns err from aws",
			lk:   "0f739d75b44acc5b",
			args: args{
				s3Svc:    &mockS3Client{},
				s3Bucket: "this-bucket-is-a-lie",
				certType: UserCertificate,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.lk.Expand(tt.args.s3Svc, tt.args.s3Bucket, tt.args.s3Prefix, tt.args.certType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Expand() got = %v, want %v", got, tt.want)
			}
		})
	}
}
