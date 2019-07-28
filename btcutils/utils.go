package btcutils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

func AddInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Add(x, y)

	return z
}

func MulInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Mul(x, y)

	return z
}

func InvInt(x, m *big.Int) *big.Int {
	exp := SubInt(m, NewInt(2))
	return ExpInt(x, exp, m)
}

func SubInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Sub(x, y)

	return z
}

func DivInt(x, y *big.Int) *big.Int {
	z := new(big.Int)
	z.Div(x, y)

	return z
}

func DivModInt(x, y *big.Int) (*big.Int, *big.Int) {
	z := new(big.Int)
	m := new(big.Int)
	return z.DivMod(x, y, m)
}

func ExpInt(x, exp, m *big.Int) *big.Int {
	z := new(big.Int)
	z.Exp(x, exp, m)

	return z
}

func PowInt(x, exp *big.Int) *big.Int {
	z := new(big.Int)
	z.Exp(x, exp, nil)

	return z
}

func ModInt(x, m *big.Int) *big.Int {
	z := new(big.Int)
	z.Mod(x, m)

	return z
}

func IsEvenInt(x *big.Int) bool {
	z := ModInt(x, NewInt(2))
	return z.Cmp(NewInt(0)) == 0
}

func ParseInt(s string, base int) (*big.Int, bool) {
	return new(big.Int).SetString(s, base)
}

func ParseBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func NewInt(i int64) *big.Int {
	return big.NewInt(i)
}

func AddFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Add(x, y)

	return z
}

func SubFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Sub(x, y)

	return z
}
func MulFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	z.Mul(x, y)

	return z
}

func DivFloat(x, y *big.Float) *big.Float {
	z := NewZeroFloat()
	// TODO check 0
	z.Quo(x, y)

	return z
}

func PowFloat(x *big.Float, exp uint64) *big.Float {
	z := NewZeroFloat().Copy(x)
	for i := uint64(0); i < (exp - 1); i++ {
		z.Mul(z, x)
	}

	return z
}

func NewZeroFloat() *big.Float {
	return big.NewFloat(0.0)
}

func LittleEndianToBigInt(b []byte) *big.Int {
	tmp := make([]byte, len(b))

	copy(tmp, b)
	ReverseBytes(tmp)
	i := new(big.Int).SetBytes(tmp)

	return i
}

func BigIntToLittleEndian(i *big.Int, l uint) []byte {
	tmp := make([]byte, l)

	b := i.Bytes()
	ReverseBytes(b)
	copy(tmp, b)

	return tmp
}

func ReverseBytes(b []byte) {
	l := len(b)
	for i := 0; i < l/2; i++ {
		b[i], b[l-1-i] = b[l-1-i], b[i]
	}
}

func EncodeVariant(i int) []byte {
	switch {
	case i < 0xfd:
		return []byte{byte(i)}
	case i < 0x10000:
		return encodeVariant(0xfd, uint16(i))
	case i < 0x100000000:
		return encodeVariant(0xfe, uint32(i))
	case i < 0x7FFFFFFFFFFFFFFF:
		return encodeVariant(0xff, uint64(i))
	default:
		panic("integer to large")
	}
}

func encodeVariant(b byte, i interface{}) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(b)

	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func ReadVariant(r io.Reader) (uint64, error) {
	var i int64
	b := make([]byte, 8)
	_, err := io.ReadFull(r, b[:1])
	if err != nil {
		return 0, err
	}

	switch b[0] {
	case 0xfd:
		i = 2
	case 0xfe:
		i = 4
	case 0xff:
		i = 8
	default:
		return uint64(b[0]), nil
	}

	_, err = io.ReadFull(r, b[:i])
	if err != nil {
		return 0, err
	}

	n := binary.LittleEndian.Uint64(b)

	return n, nil
}

func Copyb(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
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

func Base58Encode(b []byte) string {
	x := ParseBytes(b)
	zero := NewInt(0)
	x58 := NewInt(58)
	mod := new(big.Int)

	base58 := make([]byte, 0, len(b))

	for x.Cmp(zero) > 0 {
		x, mod = DivModInt(x, x58)
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

func ReadByetes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
