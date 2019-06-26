package finitefield

import (
	"errors"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
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

	num := u.AddInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Sub(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	num := u.SubInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Mul(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	num := u.MulInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Div(x, y *element) (*element, error) {
	if err := Cmp(x, y); err != nil {
		return nil, err
	}

	powNum := u.SubInt(x.prime, big.NewInt(2))

	num := u.ExpInt(y.num, powNum, x.prime)
	num = u.MulInt(x.num, num)
	num.Mod(num, x.prime)

	return &element{num, x.prime}, nil
}

func Pow(x *element, exp *big.Int) *element {
	modNum := u.SubInt(x.prime, big.NewInt(1))

	num := u.ModInt(exp, modNum)
	num = u.ExpInt(x.num, num, x.prime)

	return &element{num, x.prime}
}
