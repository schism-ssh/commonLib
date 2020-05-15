package protocol

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

type LookupKey string

func GenerateLookupKey(ident string, principals []string) LookupKey {
	sort.Strings(principals)
	lookupList := append([]string{ident}, principals...)
	lookupString := strings.Join(lookupList, ",")
	return LookupKey(fmt.Sprintf("%x", sha256.Sum256([]byte(lookupString))))
}
