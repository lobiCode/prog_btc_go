package btcutils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

const (
	alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

var ErrBadAddress = errors.New("bad address")

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

func DecodeBase58Checksum(s string) ([]byte, error) {
	b := Base58Decode(s)
	checksum := b[len(b)-4:]
	b256 := Hash256(b[:len(b)-4])[:4]
	if !bytes.Equal(b256, checksum) {
		return nil, ErrBadAddress
	}

	return b[0 : len(b)-4], nil
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

func Base58Decode(s string) []byte {
	// remove ones
	countZero := 0
	b := []byte{}
	for i := 0; i < len(s); i++ {
		if s[i] != '1' {
			break
		}
		b = append(b, 0)
		countZero++
	}

	s = s[countZero:]

	num := NewInt(0)
	x58 := NewInt(58)

	for i := 0; i < len(s); i++ {
		num = MulInt(num, x58)
		m := strings.IndexByte(alphabet, s[i])
		num = AddInt(num, NewInt(int64(m)))
	}

	b = append(b, num.Bytes()...)
	return b
}

func H160ToP2pkhAddress(h160 []byte, testnet bool) string {
	address := make([]byte, 1, len(h160)+1)
	if testnet {
		address[0] = 0x6f
	} else {
		address[0] = 0x00
	}

	return EncodeBase58Checksum(address)
}

func H160ToP2shAddress(h160 []byte, testnet bool) string {
	address := make([]byte, 1, len(h160)+1)
	if testnet {
		address[0] = 0xc4
	} else {
		address[0] = 0x05
	}

	return EncodeBase58Checksum(address)
}

func ReadByetes(r io.Reader, n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeNum(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeNum(b []byte) (int64, error) {
	var i int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func DecodeNumLittleEndian(b []byte) (int64, error) {
	var i int64
	bb := make([]byte, 8)
	copy(bb, b)
	buf := bytes.NewReader(bb)
	err := binary.Read(buf, binary.LittleEndian, &i)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func Read(r io.Reader, u int64) ([]byte, error) {
	b := make([]byte, u)
	_, err := io.ReadFull(r, b)
	return b, err
}

func BitsToTarget(bits []byte) *big.Int {
	coefficient, err := DecodeNumLittleEndian(bits[:3])
	// TODO
	if err != nil {
		panic(err)
	}
	exp := NewInt(int64(bits[3]))

	target := MulInt(NewInt(coefficient), PowInt(NewInt(256), SubInt(exp, NewInt(3))))
	return target
}

func TargetToBits(target *big.Int) []byte {
	targetB := target.Bytes()
	targetB = bytes.TrimLeftFunc(targetB, IsZeroPrefix)

	var exp int
	coefficient := make([]byte, 3, 4)

	if targetB[0] == 0x7f {
		exp = len(targetB) + 1
		coefficient[0] = 0x00
		copy(coefficient[1:], targetB[:2])
	} else {
		exp = len(targetB)
		copy(coefficient, targetB)
	}
	ReverseBytes(coefficient)
	coefficient = append(coefficient, byte(exp))

	return coefficient
}

func IsZeroPrefix(r rune) bool {
	if uint32(r) == uint32(0) {
		return true
	}
	return false
}

var (
	TWO_WEEKS   int64 = 60 * 60 * 24 * 14
	EIGHT_WEEKS int64 = TWO_WEEKS * 4
	HALF_WEEK   int64 = TWO_WEEKS / 4
)

func CalculateNewBits(timeDiff int64, prevBits []byte) []byte {

	if timeDiff > EIGHT_WEEKS {
		timeDiff = EIGHT_WEEKS
	} else if timeDiff < HALF_WEEK {
		timeDiff = HALF_WEEK
	}

	prevTarger := BitsToTarget(prevBits)
	newTarget := DivInt(MulInt(prevTarger, NewInt(timeDiff)), NewInt(TWO_WEEKS))

	return TargetToBits(newTarget)
}
