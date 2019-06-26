package ellipticcurve

import (
	. "math/big"
	"reflect"
	"testing"
)

func TestNewPoint(t *testing.T) {
	tests := []struct {
		test       string
		x, y, a, b *Float
		expPoint   bool
		err        error
	}{
		{"error expected 1", NewFloat(2.0), NewFloat(2.0), NewFloat(2.0), NewFloat(2.0), false, ErrEllipticCurvePointNotOnCurve},
		{"error expected 2", nil, NewFloat(2.0), NewFloat(2.0), NewFloat(2.0), false, ErrEllipticCurvePointNotOnCurve},
		{"error expected 3", NewFloat(2.0), nil, NewFloat(2.0), NewFloat(2.0), false, ErrEllipticCurvePointNotOnCurve},
		{"new point 1", nil, nil, NewFloat(2.0), NewFloat(2.0), true, nil},
		{"new point 2", NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0), false, nil},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			p, err := NewPoint(test.x, test.y, test.a, test.b)
			check(test.err, err, t)
			if test.expPoint {
				checkFloat(test.x, p.x, t)
				checkFloat(test.y, p.y, t)
				checkFloat(test.a, p.a, t)
				checkFloat(test.b, p.b, t)
			}
		})
	}
}

func TestEq(t *testing.T) {
	tests := []struct {
		test     string
		p1, p2   *point
		expected bool
	}{
		{"eq 1", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(18.0), NewFloat(77.0), NewFloat(5.0), NewFloat(7.0)}, false},
		{"eq 2", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)}, true},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			result := Eq(test.p1, test.p2)
			check(test.expected, result, t)
		})
	}
}

func TestNe(t *testing.T) {
	tests := []struct {
		test     string
		p1, p2   *point
		expected bool
	}{
		{"ne 1", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(18.0), NewFloat(77.0), NewFloat(5.0), NewFloat(7.0)}, true},
		{"ne 2", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)}, false},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			result := Ne(test.p1, test.p2)
			check(test.expected, result, t)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		test             string
		p1, p2, expPoint *point
		err              error
	}{
		{"add 1", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(6.0), NewFloat(7.0)},
			&point{NewFloat(18.0), NewFloat(77.0), NewFloat(5.0), NewFloat(7.0)}, nil, ErrEllipticCurvePointsNotOnSameCurve},
		{"add 2", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(18.0), NewFloat(77.0), NewFloat(5.0), NewFloat(8.0)}, nil, ErrEllipticCurvePointsNotOnSameCurve},

		{"add 3", &point{NewFloat(3.0), NewFloat(-7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(3.0), NewFloat(77.0), NewFloat(5.0), NewFloat(7.0)},
			&point{nil, nil, NewFloat(5.0), NewFloat(7.0)}, nil},
		{"add 4",
			&point{NewFloat(3.0), NewFloat(7.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(-1.0), NewFloat(-1.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(2.0), NewFloat(-5), NewFloat(5.0), NewFloat(7.0)}, nil,
		},
		{"add 5",
			&point{NewFloat(-1.0), NewFloat(-1.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(-1.0), NewFloat(-1.0), NewFloat(5.0), NewFloat(7.0)},
			&point{NewFloat(18.0), NewFloat(77.0), NewFloat(5.0), NewFloat(7.0)}, nil,
		},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			result, err := Add(test.p1, test.p2)
			check(test.err, err, t)
			if test.expPoint != nil {
				checkFloat(test.expPoint.x, result.x, t)
				checkFloat(test.expPoint.y, result.y, t)
				checkFloat(test.expPoint.a, result.a, t)
				checkFloat(test.expPoint.b, result.b, t)
			}
		})
	}
}

func checkFloat(expected, recived *Float, t *testing.T) {
	t.Helper()
	if expected == nil && recived == nil {
		return
	}

	if expected.Cmp(recived) != 0 {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
