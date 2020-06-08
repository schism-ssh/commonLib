package protocol

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// LookupKey is used for storing certificates in S3, 64-character sha256sum
//
// Partial keys are allowed, use lk.Expand() to attempt to fetch the full key
type LookupKey string

// Expand expands a LookupKey into a full 64-character key,
// given a partial key that matches a singular certificate bundle of the given type
//
// Returns an error if there is not a singular match or S3 calls fail.
func (lk LookupKey) Expand(s3Svc s3iface.S3API, s3Bucket string, s3Prefix string, certType CertType) (LookupKey, error) {
	fullPrefix := fmt.Sprintf("%s%ss/%s", s3Prefix, certType, lk)
	objs, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s3Bucket),
		Prefix: aws.String(fullPrefix),
	})
	if err != nil {
		return "", err
	}
	switch count := *objs.KeyCount; {
	case count > 1:
		return "", fmt.Errorf("partial key '%s' matches multiple certificates", lk)
	case count == 1:
		expndPrts := strings.Split(*objs.Contents[0].Key, "/")
		expnd := strings.Split(expndPrts[len(expndPrts)-1], ".")[0]
		return LookupKey(expnd), nil
	default:
		return "", fmt.Errorf("partial key '%s' matches zero certificates", lk)
	}
}

// GenerateLookupKey creates a 64-character long sha256sum key
// based on the given Identity and Principals
func GenerateLookupKey(ident string, principals []string) LookupKey {
	sort.Strings(principals)
	lookupList := append([]string{ident}, principals...)
	lookupString := strings.Join(lookupList, ",")
	return LookupKey(fmt.Sprintf("%x", sha256.Sum256([]byte(lookupString))))
}
