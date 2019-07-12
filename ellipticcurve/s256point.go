package ellipticcurve

import (
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ff "github.com/lobiCode/prog_btc_go/finitefield"
)

type SECP256k1Curve struct {
	A, B         *ff.Element
	P, N, Gx, Gy *big.Int
	G            *Point
}

var BTCCurve *SECP256k1Curve

func init() {
	BTCCurve = GetSECP256k1Curve()
}

func GetSECP256k1Curve() *SECP256k1Curve {
	// TODO
	p := u.SubInt(u.SubInt((u.PowInt(big.NewInt(2), big.NewInt(256))), (u.PowInt(big.NewInt(2), big.NewInt(32)))), big.NewInt(977))
	n, _ := new(big.Int).SetString("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 0)
	gx, _ := new(big.Int).SetString("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 0)
	gy, _ := new(big.Int).SetString("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 0)
	a, err := ff.NewS256Field(big.NewInt(0), p)
	if err != nil {
		panic(err.Error())
	}

	b, err := ff.NewS256Field(big.NewInt(7), p)
	if err != nil {
		panic(err.Error())
	}

	x, err := ff.NewS256Field(gx, p)
	if err != nil {
		panic(err.Error())
	}
	y, err := ff.NewS256Field(gy, p)
	if err != nil {
		panic(err.Error())
	}

	g, err := NewPoint(x, y, a, b)
	if err != nil {
		panic(err.Error())
	}

	return &SECP256k1Curve{a, b, p, n, gx, gy, g}
}

func NewS256Point(x, y *ff.Element) (*Point, error) {
	return NewPoint(x, y, BTCCurve.A, BTCCurve.B)
}
