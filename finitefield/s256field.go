package finitefield

import (
	"math/big"
)

func NewS256Field(num, p *big.Int) (*Element, error) {
	// TODO

	return &Element{num, p}, nil
}
