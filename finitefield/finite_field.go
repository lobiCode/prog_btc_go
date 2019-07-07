package finitefield

import (
	"errors"
	"fmt"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

var (
	ErrFiniteFieldNumNotInRange = errors.New("Number not in field range 0 to prime - 1")
	ErrFiniteFieldDiffFields    = errors.New("Number of different fields")
)

type Element struct {
	num   *big.Int
	prime *big.Int
}

func NewElement(num, prime int64) (*Element, error) {
	if num >= prime {
		return nil, ErrFiniteFieldNumNotInRange
	}

	return &Element{big.NewInt(num), big.NewInt(prime)}, nil
}

func (f *Element) GetPrime() *big.Int {
	return f.prime
}

func (f *Element) GetNum() *big.Int {
	return f.num
}

func (f *Element) String() string {
	return fmt.Sprintf("%s, %s", f.num, f.prime)
}

func Eq(x, y *Element) bool {
	if x == nil && y == nil {
		return true
	}

	if x == nil || y == nil {
		return false
	}

	return (x.num.Cmp(y.num) == 0 && x.prime.Cmp(y.prime) == 0)
}

func Ne(x, y *Element) bool {
	return !Eq(x, y)
}

func Cmp(x, y *Element) error {
	if x != nil && y != nil && x.prime.Cmp(y.prime) == 0 {
		return nil
	}

	return ErrFiniteFieldDiffFields
}

func Add(x, y *Element) *Element {
	if err := Cmp(x, y); err != nil {
		panic(err.Error())
	}

	num := u.AddInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &Element{num, x.prime}
}

func Sub(x, y *Element) *Element {
	if err := Cmp(x, y); err != nil {
		panic(err.Error())
	}

	num := u.SubInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &Element{num, x.prime}
}

func Mul(x, y *Element) *Element {
	if err := Cmp(x, y); err != nil {
		panic(err.Error())
	}

	num := u.MulInt(x.num, y.num)
	num = num.Mod(num, x.prime)

	return &Element{num, x.prime}
}

func Div(x, y *Element) *Element {
	if err := Cmp(x, y); err != nil {
		panic(err.Error())
	}

	powNum := u.SubInt(x.prime, big.NewInt(2))

	num := u.ExpInt(y.num, powNum, x.prime)
	num = u.MulInt(x.num, num)
	num.Mod(num, x.prime)

	return &Element{num, x.prime}
}

func Pow(x *Element, exp *big.Int) *Element {
	modNum := u.SubInt(x.prime, big.NewInt(1))

	num := u.ModInt(exp, modNum)
	num = u.ExpInt(x.num, num, x.prime)

	return &Element{num, x.prime}
}

func RMul(x *Element, cof *big.Int) *Element {
	num := u.ModInt(u.MulInt(x.num, cof), x.prime)

	return &Element{num, x.prime}
}
