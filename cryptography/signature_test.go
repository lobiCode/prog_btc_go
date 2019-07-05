package cryptography

import (
	"reflect"
	"testing"
)

func TestSign(t *testing.T) {
	text := "Bitcoin Bitcoin"
	pk := NewPrivateKey("ifkdafkfkfiasfiodidafpasfjadsf")

	signature := pk.Sign(text)

	ok := Verify(text, signature, pk.point)
	check(true, ok, t)
}

func TestSignFail(t *testing.T) {
	text := "Bitcoin Bitcoin"
	pk1 := NewPrivateKey("ifkdafkfkfiasfiodidafpasfjadsf")
	pk2 := NewPrivateKey("111899998900")

	signature := pk1.Sign(text)

	ok := Verify(text, signature, pk2.point)
	check(false, ok, t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
