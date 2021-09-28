package i18nc_test

import (
	"bytes"
	"fmt"
	"github.com/dcb9/i18nc"
	"testing"
)

func TestFromBytes(t *testing.T) {
	w := &bytes.Buffer{}

	if err := i18nc.FromBytes([]byte(`{
    "listen": "Stay a while and listen",
    "person": {
        "one": "{{.PluralCount}} person",
        "other": "{{.PluralCount}} people"
    },
    "hi_to": "hi {{.Name}}"
}`), "i18nc", w); err != nil {
		fmt.Println(err)
	}

	fmt.Println(w.String())
}
