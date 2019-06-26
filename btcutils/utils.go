package btcutils

import "math/big"

func AddInt(x, y *big.Int) *big.Int {
	z := NewInt(0)
	z.Add(x, y)

	return z
}

func MulInt(x, y *big.Int) *big.Int {
	z := NewInt(0)
	z.Mul(x, y)

	return z
}

func SubInt(x, y *big.Int) *big.Int {
	z := NewInt(0)
	z.Sub(x, y)

	return z
}

func ExpInt(x, exp, m *big.Int) *big.Int {
	z := NewInt(0)
	z.Exp(x, exp, m)

	return z
}

func ModInt(x, m *big.Int) *big.Int {
	z := NewInt(0)
	z.Mod(x, m)

	return z
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
