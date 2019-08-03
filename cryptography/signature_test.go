package cryptography

import (
	"encoding/hex"
	"reflect"
	"testing"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ec "github.com/lobiCode/prog_btc_go/ellipticcurve"
)

func TestSign(t *testing.T) {
	z := GetHash256Int("Bitcoin Bitcoin")
	pk := NewPrivateKey(GetHash256Int("ifkdafkfkfiasfiodidafpasfjadsf"))

	signature := pk.Sign(z)

	ok := Verify(z, signature, pk.point)
	check(true, ok, t)
}

func TestSignFail(t *testing.T) {
	z := GetHash256Int("Bitcoin Bitcoin")
	pk1 := NewPrivateKey(GetHash256Int("ifkdafkfkfiasfiodidafpasfjadsf"))
	pk2 := NewPrivateKey(GetHash256Int("111899998900"))

	signature := pk1.Sign(z)

	ok := Verify(z, signature, pk2.point)
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
			pk := NewPrivateKey(GetHash256Int(test.secret))
			pub := pk.Sec(test.compressed)
			parsePub, _ := ParsePublicKey(pub)

			check(true, ec.Eq(pk.point, parsePub), t)
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
	expected := "3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec"
	r, _ := u.ParseInt("0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 0)
	s, _ := u.ParseInt("0x8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 0)
	sig := &Signature{r, s}
	check(expected, sig.String(), t)

	b, _ := hex.DecodeString(expected)
	sigResult, err := ParseSignature(b)

	check(nil, err, t)
	check(sig, sigResult, t)
}

func TestDer2(t *testing.T) {
	text := "3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6"
	sig, _ := hex.DecodeString(text)

	sigResult, err := ParseSignature(sig)

	check(nil, err, t)
	check(text, sigResult.String(), t)
}

func TestGetH160Address(t *testing.T) {
	addr := "mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf"
	b, err := GetH160Address(addr)
	result := hex.EncodeToString(b)
	expected := "507b27411ccf7f16f10297de6cef3f291623eddf"
	check(nil, err, t)
	check(expected, result, t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
