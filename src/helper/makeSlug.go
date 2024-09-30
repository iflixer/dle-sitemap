package helper

import (
	"github.com/alexsergivan/transliterator"
	"regexp"
	"strings"
)

func MakeSlug(s string) (res string) {
	res = strings.TrimSpace(s)
	trans := transliterator.NewTransliterator(nil)
	res = trans.Transliterate(res, "")
	res = strings.ReplaceAll(res, " ", "-")
	reg, err := regexp.Compile("[^a-zA-Z0-9_]+")
	if err != nil {
		panic(err)
	}
	res = reg.ReplaceAllString(res, "")
	// Convert to lowercase
	res = strings.ToLower(res)

	return
}
