package ellipticcurve

import (
	"math/big"
	"reflect"
	"testing"

	ff "github.com/lobiCode/prog_btc_go/finitefield"
)

func TestNewPoint(t *testing.T) {
	tests := []struct {
		test              string
		x, y, a, b, prime int64
		expPoint          bool
		err               error
	}{
		{"new Point 1", 17, 64, 0, 7, 103, true, nil},
		{"new Point 2", 192, 105, 0, 7, 223, true, nil},
		{"new Point 3", 17, 56, 0, 7, 223, true, nil},
		{"new Point 4", 1, 193, 0, 7, 223, true, nil},
		{"new Point err 1", 200, 119, 0, 7, 223, false, ErrEllipticCurvePointNotOnCurve},
		{"new Point err 2", 42, 99, 0, 7, 223, false, ErrEllipticCurvePointNotOnCurve},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			x := newElement(test.x, test.prime)
			y := newElement(test.y, test.prime)
			a := newElement(test.a, test.prime)
			b := newElement(test.b, test.prime)
			p, err := NewPoint(x, y, a, b)
			check(test.err, err, t)
			if test.expPoint {
				checkPoint(&Point{x, y, a, b}, p, t)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		test       string
		prime      int64
		p1x, p1y   int64
		p2x, p2y   int64
		a, b       int64
		expx, expy int64
	}{
		{"add 1", 223, 192, 105, 17, 56, 0, 7, 170, 142},
		{"add 2", 223, 170, 142, 60, 139, 0, 7, 220, 181},
		{"add 3", 223, 47, 71, 17, 56, 0, 7, 215, 68},
		{"add 4", 223, 143, 98, 76, 66, 0, 7, 47, 71},
	}

	for _, test := range tests {
		p1 := &Point{
			newElement(test.p1x, test.prime),
			newElement(test.p1y, test.prime),
			newElement(test.a, test.prime),
			newElement(test.b, test.prime),
		}
		p2 := &Point{
			newElement(test.p2x, test.prime),
			newElement(test.p2y, test.prime),
			newElement(test.a, test.prime),
			newElement(test.b, test.prime),
		}
		exp := &Point{
			newElement(test.expx, test.prime),
			newElement(test.expy, test.prime),
			newElement(test.a, test.prime),
			newElement(test.b, test.prime),
		}
		t.Run(test.test, func(t *testing.T) {
			result := Add(p1, p2)
			checkPoint(exp, result, t)
		})
	}
}

func TestRmul(t *testing.T) {
	tests := []struct {
		test       string
		cof        int64
		prime      int64
		p1x, p1y   int64
		a, b       int64
		expx, expy int64
	}{
		{"rmul 1", 2, 223, 192, 105, 0, 7, 49, 71},
		{"rmul 2", 2, 223, 143, 98, 0, 7, 64, 168},
		{"rmul 3", 2, 223, 47, 71, 0, 7, 36, 111},
		{"rmul 4", 4, 223, 47, 71, 0, 7, 194, 51},
		{"rmul 5", 21, 223, 47, 71, 0, 7, -1, -1},
	}

	for _, test := range tests {
		p1 := &Point{
			newElement(test.p1x, test.prime),
			newElement(test.p1y, test.prime),
			newElement(test.a, test.prime),
			newElement(test.b, test.prime),
		}
		exp := &Point{
			newElement(test.expx, test.prime),
			newElement(test.expy, test.prime),
			newElement(test.a, test.prime),
			newElement(test.b, test.prime),
		}
		t.Run(test.test, func(t *testing.T) {
			result := RMul(p1, big.NewInt(test.cof))
			checkPoint(exp, result, t)
		})
	}
}
func newElement(num, prime int64) *ff.Element {
	if num == -1 {
		return nil
	}

	e, _ := ff.NewElement(num, prime)

	return e
}

func checkPoint(expected, recived *Point, t *testing.T) {
	t.Helper()
	if expected == nil && recived == nil {
		return
	}

	if Ne(expected, recived) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
