package protocol

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

const (
	TestValidBucket = "schism-test"
)

type MockS3Client struct {
	s3iface.S3API
	T *testing.T
}

func (m *MockS3Client) ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	output := &s3.ListObjectsV2Output{}
	var contents []*s3.Object
	switch *input.Prefix {
	case "Signed-Certs/host:55e8182ec4413d51":
		contents = []*s3.Object{{Key: aws.String("host:55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d.json")}}
	case "Signed-Certs/user:4":
		contents = []*s3.Object{
			{Key: aws.String("user:4e1586bed08190ccac4056078afed44daac058e8361b216dd078c7714b874cae.json")},
			{Key: aws.String("user:4d5b5d59343254c4fccafe48813ceeb99ae5ce44c1b97113b370a93f8411a01e.json")},
		}
	case "Signed-Certs/host:d0c671a71f190313":
		contents = []*s3.Object{{Key: aws.String("hosts/d0c671a71f190313333bb79ed1a98fe7414da1089b3740de4ad5056c215512e7.json")}}
	default:
		// No matches
	}

	if *input.Bucket != TestValidBucket {
		// return error
		return output, fmt.Errorf("NoSuchBucket")
	}
	output.SetContents(contents)
	output.SetKeyCount(int64(len(contents)))
	return output, nil
}

func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	output := &s3.GetObjectOutput{}
	var body string
	if strings.HasPrefix(*input.Key, S3CaPubkeyPrefix) {
		body = HelperLoadString(m.T, "valid_ca_pubkey_s3_body.json")
		// CA Pubkey Request
	} else if strings.HasPrefix(*input.Key, S3CertStoragePrefix) {
		body = HelperLoadString(m.T, "valid_signed_cert_s3_body.json")
	} else if strings.HasPrefix(*input.Key, "empty-objects") {
		body = ""
	} else {
		return nil, fmt.Errorf("error loading object %s", *input.Key)
	}
	output.SetBody(ioutil.NopCloser(strings.NewReader(body)))
	return output, nil
}

func HelperLoadString(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("testdata", name)
	bytes, _ := ioutil.ReadFile(path)
	return string(bytes)
}
