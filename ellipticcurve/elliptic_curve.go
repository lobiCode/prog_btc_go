package ellipticcurve

import (
	"errors"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

var (
	ErrEllipticCurvePointNotOnCurve      = errors.New("Point not on the curve")
	ErrEllipticCurvePointsNotOnSameCurve = errors.New("Points not on the same curve")
)

type point struct {
	x, y, a, b *big.Float
}

func NewPoint(x, y, a, b *big.Float) (*point, error) {
	if (x != nil && y == nil) || (x == nil && y != nil) {
		return nil, ErrEllipticCurvePointNotOnCurve
	}

	if x == nil && y == nil {
		return &point{x, y, a, b}, nil
	}

	l := u.PowFloat(y, 2)
	axb := u.AddFloat(u.MulFloat(x, a), b)
	x3 := u.PowFloat(x, 3)
	r := u.AddFloat(x3, axb)

	if l.Cmp(r) == 0 {
		return &point{x, y, a, b}, nil
	}

	return nil, ErrEllipticCurvePointNotOnCurve
}

func Eq(p1, p2 *point) bool {
	return p1.x.Cmp(p2.x) == 0 && p1.y.Cmp(p2.y) == 0 && p1.a.Cmp(p2.a) == 0 && p1.b.Cmp(p2.b) == 0
}

func Ne(p1, p2 *point) bool {
	return !Eq(p1, p2)
}

func Add(p1, p2 *point) (*point, error) {
	if p1.a.Cmp(p2.a) != 0 || p1.b.Cmp(p2.b) != 0 {
		return nil, ErrEllipticCurvePointsNotOnSameCurve
	}

	if p1.x == nil {
		return p2, nil
	}

	if p2.x == nil {
		return p1, nil
	}

	cmpX := p1.x.Cmp(p2.x)
	cmpY := p1.y.Cmp(p2.y)
	zeroY := p1.y.Cmp(u.NewZeroFloat())
	if (cmpX == 0 && cmpY != 0) || (cmpX == 0 && cmpY == 0 && zeroY == 0) {
		return &point{nil, nil, p1.a, p1.b}, nil
	}

	var slope *big.Float

	if cmpX != 0 {
		// s = (y2 - y1)/(x2 - x1)
		slope = u.DivFloat(u.SubFloat(p2.y, p1.y), u.SubFloat(p2.x, p1.x))
	} else {
		slope = u.MulFloat(big.NewFloat(3.0), u.PowFloat(p1.x, 2))
		slope = u.AddFloat(slope, p1.a)
		slope = u.DivFloat(slope, u.MulFloat(big.NewFloat(2.0), p1.y))
	}

	x3 := calculateX3(slope, p1.x, p2.x)
	y3 := calculateY3(slope, p1.x, x3, p1.y)

	return &point{x3, y3, p1.a, p1.b}, nil
}

func calculateX3(slope, x1, x2 *big.Float) *big.Float {
	// x3 = s^2 - x1 - x2
	x3 := u.SubFloat(u.PowFloat(slope, 2), x1)
	x3 = u.SubFloat(x3, x2)

	return x3
}

func calculateY3(slope, x1, x3, y1 *big.Float) *big.Float {
	// y3 = s(x1 - x3) - y1
	y3 := u.MulFloat(slope, u.SubFloat(x1, x3))
	y3 = u.SubFloat(y3, y1)

	return y3
}
