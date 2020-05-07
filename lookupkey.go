package commonLib

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

func LookupKey(ident string, principals []string) string {
	sort.Strings(principals)
	lookupList := append([]string{ident}, principals...)
	lookupString := strings.Join(lookupList, ",")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(lookupString)))
}
