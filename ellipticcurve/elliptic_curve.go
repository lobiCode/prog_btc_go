package ellipticcurve

import (
	"errors"
	"fmt"
	"math/big"

	ff "github.com/lobiCode/prog_btc_go/finitefield"
)

var (
	ErrEllipticCurvePointNotOnCurve      = errors.New("Point not on the curve")
	ErrEllipticCurvePointsNotOnSameCurve = errors.New("Points not on the same curve")
)

type Point struct {
	x, y, a, b *ff.Element
}

func (p *Point) String() string {
	return fmt.Sprintf("%s, %s, %s, %s", p.x, p.y, p.a, p.b)
}

func (p *Point) GetX() *ff.Element {
	return p.x
}

func (p *Point) GetXbytes() []byte {
	// TODO nil
	return p.x.GetNumBytes()
}

func (p *Point) GetYbytes() []byte {
	// TODO nil
	return p.y.GetNumBytes()
}

func (p *Point) GetXhex() string {
	// TODO nil
	return p.x.GetNumHex()
}

func (p *Point) GetYhex() string {
	// TODO nil
	return p.y.GetNumHex()
}

func (p *Point) IsYeven() bool {
	// TODO nil
	return p.y.IsEven()
}

func NewPoint(x, y, a, b *ff.Element) (*Point, error) {
	if (x != nil && y == nil) || (x == nil && y != nil) {
		return nil, ErrEllipticCurvePointNotOnCurve
	}

	if x == nil && y == nil {
		return &Point{x, y, a, b}, nil
	}

	l := ff.Pow(y, big.NewInt(2))
	r := CalcEcRightSide(x, a, b)

	if ff.Eq(l, r) {
		return &Point{x, y, a, b}, nil
	}

	return nil, ErrEllipticCurvePointNotOnCurve
}

func CalcEcRightSide(x, a, b *ff.Element) *ff.Element {
	axb := ff.Add(ff.Mul(x, a), b)
	x3 := ff.Pow(x, big.NewInt(3))
	r := ff.Add(x3, axb)

	return r
}

func Eq(p1, p2 *Point) bool {
	return ff.Eq(p1.x, p2.x) && ff.Eq(p1.y, p2.y) && ff.Eq(p1.a, p2.a) && ff.Eq(p1.b, p2.b)
}

func Ne(p1, p2 *Point) bool {
	return !Eq(p1, p2)
}

func Add(p1, p2 *Point) *Point {
	if ff.Ne(p1.a, p2.a) || ff.Ne(p1.b, p2.b) {
		panic(ErrEllipticCurvePointsNotOnSameCurve.Error())
	}

	if p1.x == nil {
		return p2
	}

	if p2.x == nil {
		return p1
	}

	cmpX := ff.Eq(p1.x, p2.x)
	cmpY := ff.Eq(p1.y, p2.y)
	// TODO
	zeroY := ff.Eq(p1.y, ff.RMul(p1.x, big.NewInt(0)))
	if (cmpX && !cmpY) || (cmpX && cmpY && zeroY) {
		return &Point{nil, nil, p1.a, p1.b}
	}

	var slope *ff.Element

	if !cmpX {
		// s = (y2 - y1)/(x2 - x1)
		slope = ff.Div(ff.Sub(p2.y, p1.y), ff.Sub(p2.x, p1.x))
	} else {
		slope = ff.RMul(ff.Pow(p1.x, big.NewInt(2)), big.NewInt(3))
		slope = ff.Add(slope, p1.a)
		slope = ff.Div(slope, ff.RMul(p1.y, big.NewInt(2)))
	}

	x3 := calculateX3(slope, p1.x, p2.x)
	y3 := calculateY3(slope, p1.x, x3, p1.y)

	return &Point{x3, y3, p1.a, p1.b}
}

func RMul(p *Point, cof *big.Int) *Point {
	result := &Point{nil, nil, p.a, p.b}
	// TODO
	curent := p

	z := new(big.Int)
	c := new(big.Int).Set(cof)
	zero := big.NewInt(0)
	one := big.NewInt(1)

	for c.Cmp(zero) > 0 {
		z = z.And(c, one)
		if z.Int64() > int64(0) {
			result = Add(result, curent)
		}
		curent = Add(curent, curent)
		c.Rsh(c, 1)
	}

	return result
}

func calculateX3(slope, x1, x2 *ff.Element) *ff.Element {
	// x3 = s^2 - x1 - x2
	x3 := ff.Sub(ff.Pow(slope, big.NewInt(2)), x1)
	x3 = ff.Sub(x3, x2)

	return x3
}

func calculateY3(slope, x1, x3, y1 *ff.Element) *ff.Element {
	// y3 = s(x1 - x3) - y1
	y3 := ff.Mul(slope, ff.Sub(x1, x3))
	y3 = ff.Sub(y3, y1)

	return y3
}
