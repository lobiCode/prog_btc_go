package cryptography

import (
	"encoding/hex"
	"reflect"
	"testing"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ec "github.com/lobiCode/prog_btc_go/ellipticcurve"
)

func TestSign(t *testing.T) {
	text := "Bitcoin Bitcoin"
	pk := NewPrivateKey(getHash256Int("ifkdafkfkfiasfiodidafpasfjadsf"))

	signature := pk.Sign(text)

	ok := Verify(text, signature, pk.point)
	check(true, ok, t)
}

func TestSignFail(t *testing.T) {
	text := "Bitcoin Bitcoin"
	pk1 := NewPrivateKey(getHash256Int("ifkdafkfkfiasfiodidafpasfjadsf"))
	pk2 := NewPrivateKey(getHash256Int("111899998900"))

	signature := pk1.Sign(text)

	ok := Verify(text, signature, pk2.point)
	check(false, ok, t)
}

func TestSignParse(t *testing.T) {
	testCase := []struct {
		test       string
		secret     string
		compressed bool
	}{
		{"t 1", "rriioeoreifkjdkafkasfjeo", false},
		{"t 2", "rriioeoreifkjdkafkasfjeo", false},
		{"t 3", "rdkfakdfkdariioeoreieo", false},
		{"t 4", "rriioeokfdkasfkakreieo", true},
		{"t 5", "kfkasdkfrriioeoreieo", true},
		{"t 6", "kjfdasfkasjfasfkasdkfrriioeoreieo", true},
		{"t 7", "kfdjsafdafldasdfdasfkasjfasfkasdkfrriioeoreieo", true},
		{"t 8", "kjdddfdasfkasjfasfkasdkfrriiiiiioeoreieo", true},
	}

	for _, test := range testCase {
		t.Run(test.test, func(t *testing.T) {
			pk := NewPrivateKey(getHash256Int(test.secret))
			pub := hex.EncodeToString(pk.Sec(test.compressed))
			parsePub, _ := ParsePublicKey(pub)

			check(true, ec.Eq(pk.point, parsePub), t)
		})
	}
}

func TestBase58Encode(t *testing.T) {
	testCase := []struct {
		test     string
		s        string
		expected string
	}{
		{"t 1", "\x00\x00", "11"},
		{"t 2", "yeijskloilk49", "B7TY7kdDFMU2gmpVqr"},
	}

	for _, test := range testCase {
		t.Run(test.test, func(t *testing.T) {
			result := Base58Encode([]byte(test.s))
			check(test.expected, result, t)
		})
	}
}

func TestAddress(t *testing.T) {
	secret, _ := u.ParseInt("0x12345deadbeef", 0)
	pk := NewPrivateKey(secret)
	address := pk.Address(true, false)
	check("1F1Pn2y6pDb68E5nYJJeba4TLg2U7B6KF1", address, t)
}

func TestDer(t *testing.T) {
	expected := "0345022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"
	r, _ := u.ParseInt("0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 0)
	s, _ := u.ParseInt("0x8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 0)
	sig := Signature{r, s}
	check(expected, sig.String(), t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
