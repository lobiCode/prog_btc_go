package finitefield

import (
	"math/big"
	"reflect"
	"testing"
)

type testOperatorFunc func(x, y *element) (*element, error)
type testOperatorCase struct {
	test          string
	xNum          int64
	xPrime        int64
	yNum          int64
	yPrime        int64
	expectedNum   int64
	expectedPrime int64
	err           error
}

func TestNewElement(t *testing.T) {
	tests := []struct {
		test  string
		num   int64
		prime int64
		err   error
	}{
		{"error expected", 19, 19, ErrFiniteFieldNumNotInRange},
		{"error not expected", 5, 19, nil},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			f, err := NewElement(test.num, test.prime)
			check(err, test.err, t)
			if test.err == nil {
				check(big.NewInt(test.num), f.GetNum(), t)
				check(big.NewInt(test.prime), f.GetPrime(), t)
			}
		})
	}
}

func TestEq(t *testing.T) {
	tests := []struct {
		test     string
		xNum     int64
		xPrime   int64
		yNum     int64
		yPrime   int64
		expected bool
	}{
		{"eq 1", 5, 7, 5, 7, true},
		{"eq 2", 5, 7, 5, 11, false},
		{"eq 3", 5, 11, 7, 11, false},
		{"eq 4", 5, 7, 8, 11, false},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			x, _ := NewElement(test.xNum, test.xPrime)
			y, _ := NewElement(test.yNum, test.yPrime)
			result := Eq(x, y)
			check(test.expected, result, t)
		})
	}
}

func TestNe(t *testing.T) {
	tests := []struct {
		test     string
		xNum     int64
		xPrime   int64
		yNum     int64
		yPrime   int64
		expected bool
	}{
		{"ne 1", 5, 7, 5, 7, false},
		{"ne 2", 5, 7, 5, 11, true},
		{"ne 3", 5, 11, 7, 11, true},
		{"ne 4", 5, 7, 8, 11, true},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			x, _ := NewElement(test.xNum, test.xPrime)
			y, _ := NewElement(test.yNum, test.yPrime)
			result := Ne(x, y)
			check(test.expected, result, t)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []testOperatorCase{
		{"add 1", 2, 31, 15, 31, 17, 31, nil},
		{"add error expected", 2, 19, 15, 31, 0, 0, ErrFiniteFieldDiffFields},
	}
	testOperators(tests, Add, t)
}

func TestSub(t *testing.T) {
	tests := []testOperatorCase{
		{"sub 1", 29, 31, 4, 31, 25, 31, nil},
		{"sub 2", 4, 31, 29, 31, 6, 31, nil},
		{"sub error expected", 2, 19, 15, 31, 0, 0, ErrFiniteFieldDiffFields},
	}
	testOperators(tests, Sub, t)
}

func TestMul(t *testing.T) {
	tests := []testOperatorCase{
		{"mul 1", 24, 31, 19, 31, 22, 31, nil},
		{"mul error expected", 2, 19, 15, 31, 0, 0, ErrFiniteFieldDiffFields},
	}
	testOperators(tests, Mul, t)
}

func TestDiv(t *testing.T) {
	tests := []testOperatorCase{
		{"div 1", 3, 31, 24, 31, 4, 31, nil},
		{"div 2", 7, 19, 5, 19, 9, 19, nil},
		{"div 3", 2, 19, 7, 19, 3, 19, nil},
		{"div 4", 0, 19, 1, 19, 0, 19, nil},
		{"div error expected", 2, 19, 15, 31, 0, 0, ErrFiniteFieldDiffFields},
	}
	testOperators(tests, Div, t)
}

func TestPow(t *testing.T) {
	tests := []struct {
		test          string
		xNum          int64
		xPrime        int64
		exp           int64
		expectedNum   int64
		expectedPrime int64
	}{
		{"pow 1", 17, 31, 3, 15, 31},
		{"pow 2", 17, 31, -3, 29, 31},
		{"pow 3", 17, 31, 0, 1, 31},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			x, _ := NewElement(test.xNum, test.xPrime)
			result := Pow(x, big.NewInt(test.exp))
			checkInt(big.NewInt(test.expectedNum), result.GetNum(), t)
			checkInt(big.NewInt(test.expectedPrime), result.GetPrime(), t)
		})
	}
}

func testOperators(tests []testOperatorCase, f testOperatorFunc, t *testing.T) {
	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			x, _ := NewElement(test.xNum, test.xPrime)
			y, _ := NewElement(test.yNum, test.yPrime)
			result, err := f(x, y)
			check(err, test.err, t)
			if test.err == nil {
				checkInt(big.NewInt(test.expectedNum), result.GetNum(), t)
				checkInt(big.NewInt(test.expectedPrime), result.GetPrime(), t)
			}
		})
	}
}

func checkInt(expected, recived *big.Int, t *testing.T) {
	t.Helper()
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
