package model

import (
	"fmt"
	"testing"
)

func TestStringArray(t *testing.T) {
	var emo = []string{"ava", "qwq", "ovo"}
	enc, _ := EncodeStringArray(emo)
	fmt.Println(enc)
	dec := DecodeStringArray(enc)
	fmt.Println(dec)
}
