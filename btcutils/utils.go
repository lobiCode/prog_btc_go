package btcutils

import (
	"math/big"
)

func AddInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Add(x, y)

	return z
}

func MulInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Mul(x, y)

	return z
}

func InvInt(x, m *big.Int) *big.Int {
	exp := SubInt(m, NewInt(2))
	return ExpInt(x, exp, m)
}

func SubInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Sub(x, y)

	return z
}

func DivInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Div(x, y)

	return z
}

func DivModInt(x, y *big.Int) (*big.Int, *big.Int) {
	z := new(big.Int)
	m := new(big.Int)
	return z.DivMod(x, y, m)
}

func ExpInt(x, exp, m *big.Int) *big.Int {
	z := new(big.Int)
	z.Exp(x, exp, m)

	return z
}

func PowInt(x, exp *big.Int) *big.Int {
	z := new(big.Int)
	z.Exp(x, exp, nil)

	return z
}

func ModInt(x, m *big.Int) *big.Int {
	z := new(big.Int)
	z.Mod(x, m)

	return z
}

func IsEvenInt(x *big.Int) bool {
	z := ModInt(x, NewInt(2))
	return z.Cmp(NewInt(0)) == 0
}

func ParseInt(s string, base int) (*big.Int, bool) {
	return new(big.Int).SetString(s, base)
}

func ParseBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func NewInt(i int64) *big.Int {
	return big.NewInt(i)
}

func AddFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Add(x, y)

	return z
}

func SubFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Sub(x, y)

	return z
}
func MulFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Mul(x, y)

	return z
}

func DivFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	// TODO check 0
	z.Quo(x, y)

	return z
}

func PowFloat(x *big.Float, exp uint64) *big.Float {
	z := NewZeroFloat().Copy(x)
	for i := uint64(0); i < (exp - 1); i++ {
		z.Mul(z, x)
	}

	return z
}

func NewZeroFloat() *big.Float {
	return big.NewFloat(0.0)
}

func LittleEndianToBigInt(b []byte) *big.Int {
	tmp := make([]byte, len(b))

	copy(tmp, b)
	ReverseBytes(tmp)
	i := new(big.Int).SetBytes(tmp)

	return i
}

func BigIntToLittleEndian(i *big.Int, l uint) []byte {
	tmp := make([]byte, l)

	b := i.Bytes()
	ReverseBytes(b)
	copy(tmp, b)

	return tmp
}

func ReverseBytes(b []byte) {
	l := len(b)
	for i := 0; i < l/2; i++ {
		b[i], b[l-1-i] = b[l-1-i], b[i]
	}
}
