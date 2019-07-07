package cryptography

import (
	"crypto/sha256"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ec "github.com/lobiCode/prog_btc_go/ellipticcurve"
)

type PrivateKey struct {
	secret *big.Int
	point  *ec.Point
}

type Signature struct {
	r, s *big.Int
}

func NewPrivateKey(secret string) *PrivateKey {
	i := getBigInt(secret)

	return &PrivateKey{getBigInt(secret), ec.RMul(ec.BTCCurve.G, i)}

	return nil
}

func (pk *PrivateKey) Sign(message string) *Signature {
	z := getBigInt(message)
	k := getDeterministicK()
	r := ec.RMul(ec.BTCCurve.G, k).GetX().GetNum()
	kInv := u.InvInt(k, ec.BTCCurve.N)
	s := u.ModInt(u.MulInt(u.AddInt(z, u.MulInt(r, pk.secret)), kInv), ec.BTCCurve.N)
	if s.Cmp(u.DivInt(ec.BTCCurve.N, big.NewInt(2))) == 1 {
		s = u.SubInt(ec.BTCCurve.N, s)
	}

	return &Signature{r, s}
}

func Verify(message string, signature *Signature, publicKey *ec.Point) bool {
	sInv := u.InvInt(signature.s, ec.BTCCurve.N)
	z := getBigInt(message)
	uu := u.ModInt(u.MulInt(z, sInv), ec.BTCCurve.N)
	v := u.ModInt(u.MulInt(signature.r, sInv), ec.BTCCurve.N)
	sum := ec.Add(ec.RMul(ec.BTCCurve.G, uu), ec.RMul(publicKey, v)).GetX().GetNum()
	if sum.Cmp(signature.r) == 0 {
		return true
	}

	return false
}

func getBigInt(s string) *big.Int {
	sum := sha256.Sum256([]byte(s))
	i := new(big.Int).SetBytes(sum[:])

	return i
}

func getDeterministicK() *big.Int {
	// TODO
	return big.NewInt(100)
}
