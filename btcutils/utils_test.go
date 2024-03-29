package btcutils

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
)

func TestLittleEndianToBigInt(t *testing.T) {
	b, _ := hex.DecodeString("99c3980000000000")
	expected := big.NewInt(10011545)
	i := LittleEndianToBigInt(b)
	check(0, i.Cmp(expected), t)
}

func TestBigToIntLittleEndian(t *testing.T) {
	tests := []struct {
		test     string
		i        *big.Int
		l        uint
		expected []byte
	}{
		{"int to le 1", big.NewInt(1), 4, []byte{1, 0, 0, 0}},
		{"int to le 2", big.NewInt(10011545), 8, []byte{153, 195, 152, 0, 0, 0, 0, 0}},
	}
	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			result := BigIntToLittleEndian(test.i, test.l)
			check(test.expected, result, t)
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

func TestBase58Decode(t *testing.T) {
	testCase := []struct {
		test     string
		expected []byte
		s        string
	}{
		{"t 1", []byte("\x00\x00"), "11"},
		{"t 2", []byte("yeijskloilk49"), "B7TY7kdDFMU2gmpVqr"},
	}

	for _, test := range testCase {
		t.Run(test.test, func(t *testing.T) {
			result := Base58Decode(test.s)
			check(test.expected, result, t)
		})
	}
}

func TestED(t *testing.T) {
	b, _ := hex.DecodeString("507b27411ccf7f16f10297de6cef3f291623eddf")
	s := EncodeBase58Checksum(b)
	d, err := DecodeBase58Checksum(s)

	check(nil, err, t)
	check(b, d, t)

}
func TestBitsToTarget(t *testing.T) {
	bits, _ := hex.DecodeString("e93c0118")
	target := BitsToTarget(bits)
	check("30353962581764818649842367179120467226026534727449575424", target.String(), t)
}

func TestCalculateNewBits(t *testing.T) {
	prevBits, _ := hex.DecodeString("54d80118")
	timeDiff := int64(302400)
	expected, _ := hex.DecodeString("00157617")

	result := CalculateNewBits(timeDiff, prevBits)

	check(expected, result, t)
}

func TestByteToBits(t *testing.T) {
	expected := []byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	input := "4000600a080000010940"
	inB, _ := hex.DecodeString(input)

	result := ByteToBits(inB)
	check(expected, result, t)
}

func TestBitsToBytes(t *testing.T) {
	input := []byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	expected := "4000600a080000010940"
	expectedB, _ := hex.DecodeString(expected)

	result := BitsToBytes(input)
	check(expectedB, result, t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
