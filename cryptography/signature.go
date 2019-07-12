package cryptography

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	ec "github.com/lobiCode/prog_btc_go/ellipticcurve"
	ff "github.com/lobiCode/prog_btc_go/finitefield"
	"golang.org/x/crypto/ripemd160"
)

var ErrPubKeyInvalidFormat = errors.New("Publick key invalid format")

const (
	alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

type Signature struct {
	r, s *big.Int
}

func (sig *Signature) Der() []byte {
	rbin := sig.r.Bytes()
	if rbin[0] >= 0x80 {
		rbin = append([]byte{0x00}, rbin...)
	}
	// TODO
	result := append([]byte{0x02, byte(len(rbin))}, rbin...)

	sbin := sig.s.Bytes()
	sbin = bytes.TrimLeftFunc(sbin, isZeroPrefix)
	if sbin[0] >= 0x80 {
		sbin = append([]byte{0x00}, sbin...)
	}
	sbin = append([]byte{0x02, byte(len(sbin))}, sbin...)

	result = append(result, sbin...)

	result = append([]byte{0x03, byte(len(result))}, result...)

	return result
}

func (sig *Signature) String() string {
	return hex.EncodeToString(sig.Der())
}

type PrivateKey struct {
	secret *big.Int
	point  *ec.Point
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

func (pk *PrivateKey) Address(compressed, testnet bool) string {
	b160 := Hash160(pk.Sec(compressed))

	result := make([]byte, 1, len(b160)+1)
	if testnet {
		result[0] = 0x6f
	} else {
		result[0] = 0x00
	}

	result = append(result, b160...)

	return EncodeBase58Checksum(result)
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

	return EncodeBase58Checksum(result)
}

func ParsePublicKey(key string) (*ec.Point, error) {
	// TODO len
	bKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	x := u.ParseBytes(bKey[1:33])
	xf, err := ff.NewS256Field(x, ec.BTCCurve.P)
	if err != nil {
		return nil, err
	}

	format := bKey[0]

	if format == 4 {
		y := u.ParseBytes(bKey[33:65])
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

func (pk *PrivateKey) Sign(message string) *Signature {
	z := getHash256Int(message)
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
	z := getHash256Int(message)
	uu := u.ModInt(u.MulInt(z, sInv), ec.BTCCurve.N)
	v := u.ModInt(u.MulInt(signature.r, sInv), ec.BTCCurve.N)
	sum := ec.Add(ec.RMul(ec.BTCCurve.G, uu), ec.RMul(publicKey, v)).GetX().GetNum()
	if sum.Cmp(signature.r) == 0 {
		return true
	}

	return false
}

func Base58Encode(b []byte) string {
	x := u.ParseBytes(b)
	zero := u.NewInt(0)
	x58 := u.NewInt(58)
	mod := new(big.Int)

	base58 := make([]byte, 0, len(b))

	for x.Cmp(zero) > 0 {
		x, mod = u.DivModInt(x, x58)
		base58 = append(base58, alphabet[mod.Int64()])
	}

	for _, i := range b {
		if i != 0 {
			break
		}
		base58 = append(base58, '1')
	}

	l := len(base58)
	for i := 0; i < l/2; i++ {
		base58[i], base58[l-1-i] = base58[l-1-i], base58[i]
	}

	return string(base58)
}

func getHash256Int(s string) *big.Int {
	sum := Hash256([]byte(s))
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

func isZeroPrefix(r rune) bool {
	if uint32(r) == uint32(0) {
		return true
	}
	return false
}

func Hash256(b []byte) []byte {
	sum := sha256.Sum256(b)
	sum = sha256.Sum256(sum[:])
	return sum[:]
}

func Hash160(b []byte) []byte {
	sum := sha256.Sum256(b)
	h := ripemd160.New()
	h.Write(sum[:])
	return h.Sum(nil)
}

func EncodeBase58Checksum(b []byte) string {
	b256 := Hash256(b)[:4]
	asum := append(b, b256...)

	return Base58Encode(asum)

}
