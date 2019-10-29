package cryptography

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ec "github.com/lobiCode/prog_btc_go/ellipticcurve"
	ff "github.com/lobiCode/prog_btc_go/finitefield"
)

var (
	ErrPubKeyInvalidFormat = errors.New("publick key invalid format")
	ErrBadSig              = errors.New("bad signature")
	ErrBadSigLength        = errors.New("bad signature length")
)

type Signature struct {
	r, s *big.Int
}

func (sig *Signature) Der() []byte {
	rbin := sig.r.Bytes()
	rbin = bytes.TrimLeftFunc(rbin, u.IsZeroPrefix)
	if rbin[0]&0x80 > 0 {
		rbin = append([]byte{0x00}, rbin...)
	}
	// TODO
	result := append([]byte{0x02, byte(len(rbin))}, rbin...)

	sbin := sig.s.Bytes()
	sbin = bytes.TrimLeftFunc(sbin, u.IsZeroPrefix)
	if sbin[0] >= 0x80 {
		sbin = append([]byte{0x00}, sbin...)
	}
	sbin = append([]byte{0x02, byte(len(sbin))}, sbin...)

	result = append(result, sbin...)

	result = append([]byte{0x30, byte(len(result))}, result...)

	return result
}

func ParseSignature(sig []byte) (*Signature, error) {
	check := []struct {
		expected byte
		err      error
		b        []byte
	}{
		{0x30, ErrBadSig, nil},
		{byte(len(sig) - 2), ErrBadSigLength, nil},
		{0x02, ErrBadSig, nil},
		{0, nil, nil},
		{0x02, ErrBadSig, nil},
		{0, nil, nil},
	}

	br := bytes.NewReader(sig)
	for i, v := range check {
		b, err := u.ReadByetes(br, 1)
		if err != nil {
			return nil, err
		}
		if i == 3 || i == 5 {
			b, err = u.ReadByetes(br, int64(b[0]))
			if err != nil {
				return nil, err
			}
			check[i].b = b
		} else {
			if b[0] != v.expected {
				return nil, v.err
			}
		}
	}

	total := len(check[3].b) + len(check[5].b) + 6
	if total != len(sig) {
		return nil, ErrBadSig
	}

	r := u.ParseBytes(check[3].b)
	s := u.ParseBytes(check[5].b)

	return &Signature{r, s}, nil
}

func (sig *Signature) String() string {
	return hex.EncodeToString(sig.Der())
}

type PrivateKey struct {
	secret *big.Int
	point  *ec.Point
}

func (pk *PrivateKey) Point() *ec.Point {
	return pk.point
}

func (pk *PrivateKey) Sec(compressed bool) []byte {

	xb := pk.point.GetXbytes()
	result := make([]byte, 1, len(xb)+1)
	result = append(result, xb...)

	if compressed {
		if pk.point.IsYeven() {
			result[0] = 0x02
		} else {
			result[0] = 0x03
		}
	} else {
		result[0] = 0x04
		result = append(result, pk.point.GetYbytes()...)
	}

	return result
}

func (pk *PrivateKey) AddressP2pkh(compressed, testnet bool) string {
	b160 := u.Hash160(pk.Sec(compressed))

	return u.AddressP2pkh(b160, testnet)
}

func (pk *PrivateKey) Wif(compressed, testnet bool) string {
	// TODO
	sb := pk.secret.Bytes()
	result := make([]byte, 1, len(sb)+2)
	result = append(result, sb...)

	if testnet {
		result[0] = 0xef
	} else {
		result[0] = 0x80
	}

	if compressed {
		result = append(result, 0x01)
	}

	return u.EncodeBase58Checksum(result)
}

func ParsePublicKey(key []byte) (*ec.Point, error) {
	// TODO len

	x := u.ParseBytes(key[1:33])
	xf, err := ff.NewS256Field(x, ec.BTCCurve.P)
	if err != nil {
		return nil, err
	}

	format := key[0]

	if format == 4 {
		y := u.ParseBytes(key[33:65])
		yf, err := ff.NewS256Field(y, ec.BTCCurve.P)
		if err != nil {
			return nil, err
		}
		return ec.NewS256Point(xf, yf)
	}

	isEven := false
	if format == 2 {
		isEven = true
	}

	yf := calculateY(xf)
	if yf.IsEven() {
		if !isEven {
			yf, err = ff.NewS256Field(u.SubInt(ec.BTCCurve.P, yf.GetNum()), ec.BTCCurve.P)
		}
	} else {
		if isEven {
			yf, err = ff.NewS256Field(u.SubInt(ec.BTCCurve.P, yf.GetNum()), ec.BTCCurve.P)
		}
	}

	if err != nil {
		return nil, err
	}

	return ec.NewS256Point(xf, yf)
}

func NewPrivateKey(secret *big.Int) *PrivateKey {
	// TODO
	return &PrivateKey{secret, ec.RMul(ec.BTCCurve.G, secret)}
}

func (pk *PrivateKey) Sign(z *big.Int) *Signature {
	k := getDeterministicK()
	r := ec.RMul(ec.BTCCurve.G, k).GetX().GetNum()
	kInv := u.InvInt(k, ec.BTCCurve.N)
	s := u.ModInt(u.MulInt(u.AddInt(z, u.MulInt(r, pk.secret)), kInv), ec.BTCCurve.N)
	if s.Cmp(u.DivInt(ec.BTCCurve.N, big.NewInt(2))) == 1 {
		s = u.SubInt(ec.BTCCurve.N, s)
	}

	return &Signature{r, s}
}

func Verify(z *big.Int, signature *Signature, publicKey *ec.Point) bool {
	sInv := u.InvInt(signature.s, ec.BTCCurve.N)
	uu := u.ModInt(u.MulInt(z, sInv), ec.BTCCurve.N)
	v := u.ModInt(u.MulInt(signature.r, sInv), ec.BTCCurve.N)
	sum := ec.Add(ec.RMul(ec.BTCCurve.G, uu), ec.RMul(publicKey, v)).GetX().GetNum()
	if sum.Cmp(signature.r) == 0 {
		return true
	}

	return false
}

func GetH160Address(address string) ([]byte, error) {
	b, err := u.DecodeBase58Checksum(address)
	if err != nil {
		return nil, err
	}

	return b[1:], nil
}

func GetHash256Int(s string) *big.Int {
	sum := u.Hash256([]byte(s))
	i := new(big.Int).SetBytes(sum[:])

	return i
}

func getDeterministicK() *big.Int {
	// TODO
	return big.NewInt(100)
}

func calculateY(x *ff.Element) *ff.Element {
	y := ec.CalcEcRightSide(x, ec.BTCCurve.A, ec.BTCCurve.B)
	y = ff.Pow(y, u.DivInt(u.AddInt(ec.BTCCurve.P, u.NewInt(1)), u.NewInt(4)))
	return y
}
