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

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
