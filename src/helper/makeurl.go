package helper

import (
	"github.com/alexsergivan/transliterator"
	"strings"
)

func MakeURL(s string) (u string) {
	trans := transliterator.NewTransliterator(nil)
	u = strings.ReplaceAll(strings.ToLower(trans.Transliterate(s, "")), " ", "_")
	return
}
