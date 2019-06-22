package finitefield

import (
	"errors"
	"math/big"
)

var (
	ErrFiniteFieldNumNotInRange = errors.New("Number not in field range 0 to prime - 1")
	ErrFiniteFieldDiffFields    = errors.New("Number of different fields")
)

type element struct {
	num   *big.Int
	prime *big.Int
}

func NewElement(num, prime int64) (*element, error) {
	if num >= prime {
		return nil, ErrFiniteFieldNumNotInRange
	}

	return &element{big.NewInt(num), big.NewInt(prime)}, nil
}

func (f *element) GetPrime() *big.Int {
	return f.prime
}

func (f *element) GetNum() *big.Int {
	return f.num
}

func Eq(x, y *element) bool {
	if x == nil || y == nil {
		return false
	}

	return (x.num.Cmp(y.num) == 0 && x.prime.Cmp(y.prime) == 0)
}

func Ne(x, y *element) bool {
	return !Eq(x, y)
}

func Cmp(x, y *element) error {
	if x != nil && y != nil && x.prime.Cmp(y.prime) == 0 {
		return nil
	}

	return ErrFiniteFieldDiffFields
}

func Add(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	num := &big.Int{}
	num.Add(x.num, y.num).Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Sub(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	num := &big.Int{}
	num.Sub(x.num, y.num).Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Mul(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	num := &big.Int{}
	num.Mul(x.num, y.num).Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Div(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	powNum := &big.Int{}
	powNum.Sub(x.prime, big.NewInt(2))

	num := &big.Int{}
	num.Exp(y.num, powNum, x.prime)
	num.Mul(x.num, num)

	num.Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Pow(x *element, exp *big.Int) *element {
	modNum := &big.Int{}
	modNum.Sub(x.prime, big.NewInt(1))

	num := &big.Int{}
	num.Mod(exp, modNum)

	num.Exp(x.num, num, x.prime)

	return &element{num, x.prime}
}
