package protocol

import (
	"fmt"
	"sort"
	"strings"

	"crypto/sha256"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// LookupKeySeparator is used to separate the cert type and the cert key
const LookupKeySeparator = ":"

// LookupKey is used for storing certificates in S3, 64-character sha256sum
//
// Partial keys are allowed, use lk.Expand() to attempt to fetch the full key
type LookupKey struct {
	Id   string
	Type CertType
}

// String returns the LookupKey as a string in the format:
//  "#{lk.Type}:#{lk.Id}"
func (lk *LookupKey) String() string {
	return fmt.Sprintf("%s%s%s", lk.Type, LookupKeySeparator, lk.Id)
}

// MarshalJSON returns the same thing as String but as a `[]byte`
//
// This will not return errors, including for the interface only.
func (lk *LookupKey) MarshalJSON() ([]byte, error) {
	return []byte(lk.String()), nil
}

// UnmarshalJSON parses the output of `LookupKey.MarshalJSON()`
// See `LookupKey.String()` for formatting
//
// Returns an error if the raw data cannot be parsed
func (lk *LookupKey) UnmarshalJSON(data []byte) error {
	id, cType, err := parseRawLookupKey(string(data))
	if err != nil {
		return err
	}
	lk.Id = id
	lk.Type = cType
	return nil
}

// Expand expands a LookupKey into a full 64-character key,
// given a partial key that matches a singular certificate bundle of the given type
// stored in the given S3 bucket (and prefix)
//
// Returns an error if there is not a singular match or S3 calls fail.
// Returns an error if the expanded key is in an invalid format
//
//  Example:
//   sampleKey := protocol.LookupKey{Id: "55e8182e", Type: protocol.HostCertificate}
//   err := sampleKey.Expand(s3Con, bucket, prfx)
//   if err != nil { panic(err) }
//   # sampleKey.Id => "55e8182ec4413d51676d1ba7480708a48c5b50f4a86b3afb9be6c43c648b373d"
func (lk *LookupKey) Expand(s3Svc s3iface.S3API, s3Bucket string, s3Prefix string) error {
	// Expand short key types first, even though the data we get back from AWS
	// \should\ be expanded already. It just makes searching easier
	lk.Type = lk.Type.Expand()
	fullPrefix := fmt.Sprintf("%s%s%s", s3Prefix, S3CertStoragePrefix, lk)
	objs, err := s3Svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s3Bucket),
		Prefix: aws.String(fullPrefix),
	})
	if err != nil {
		return err
	}
	switch count := *objs.KeyCount; {
	case count > 1:
		return fmt.Errorf("partial key '%s' matches multiple certificates", lk)
	case count == 1:
		expndPrts := strings.Split(*objs.Contents[0].Key, "/")
		expnd := strings.Split(expndPrts[len(expndPrts)-1], ".")[0]
		lk.Id, lk.Type, err = parseRawLookupKey(expnd)
		return err
	default:
		return fmt.Errorf("partial key '%s' matches zero certificates", lk)
	}
}

// parseRawLookupKey takes a string in the same format that `String()` provides
// and returns the two sub-components of the Key
//
// Returns an error if the key is improperly formatted
func parseRawLookupKey(rawKey string) (string, CertType, error) {
	var (
		typeInd = 0
		idInd   = 1
	)
	parts := strings.Split(rawKey, LookupKeySeparator)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unable to parse raw key '%s'", rawKey)
	}
	return parts[idInd], CertType(parts[typeInd]), nil
}

// ParseLookupKey takes a string in the same format that `String()` provides
// and returns a pointer to a new LookupKey object
//
// returns an error if the key is improperly formatted
//
// Expansion does NOT happen here. See `lk.Expand()`
// if you wish to resolve a partial key
func ParseLookupKey(rawKey string) (*LookupKey, error) {
	var err error
	lk := &LookupKey{}
	lk.Id, lk.Type, err = parseRawLookupKey(rawKey)
	return lk, err
}

// GenerateLookupKey creates a 64-character long sha256sum key
// based on the given Identity and Principals
//
// The Identity is joined with the sorted list of Principals
// then separated with commas and run through sha256.Sum256()
//
//  Example:
//   ident := "someUser@dev1.example.com"
//   princs := []string{"someUser", "admin"}
//   sampleKey := protocol.GenerateLookupKey(ident, princs, protocol.UserCertificate)
//   // sampleKey => {
//   //  Id: "a5ba427b532c152b3e9cded5ab36f040072f7582a455271fd26d1fc696c7ac64",
//   //  Type: "user",
//   // }
func GenerateLookupKey(ident string, principals []string, certType CertType) *LookupKey {
	sort.Strings(principals)
	lookupList := append([]string{ident}, principals...)
	lookupString := strings.Join(lookupList, ",")
	return &LookupKey{
		Id:   fmt.Sprintf("%x", sha256.Sum256([]byte(lookupString))),
		Type: certType,
	}
}
