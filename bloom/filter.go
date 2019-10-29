package bloom

import (
	"github.com/twmb/murmur3"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

const BIP37_CONSTANT uint32 = 0xfba4c795

type Filter struct {
	Size          uint32
	FieldSize     uint32
	FunctionCount uint32
	Tweak         uint32
	BitField      []byte
}

func NewFilter(size, functionCount, tweak uint32) *Filter {
	return &Filter{
		Size:          size,
		FieldSize:     size * 8,
		FunctionCount: functionCount,
		Tweak:         tweak,
		BitField:      make([]byte, size*8),
	}
}

// TODO!!!!
func (f *Filter) Add(phrase []byte) {
	for i := uint32(0); i < f.FunctionCount; i++ {
		seed := uint32(i)*BIP37_CONSTANT + uint32(f.Tweak)
		sum := murmur3.SeedSum32(seed, phrase)
		f.BitField[sum%uint32(f.FieldSize)] = 1
	}
}

func (f *Filter) FilterBytes() []byte {
	return u.BitsToBytes(f.BitField)
}

func (f *Filter) FilterLoad(flag uint8) []byte {
	payload := make([]byte, 0, 9+f.Size+4)
	payload = append(payload, u.EncodeVariant(int(f.Size))...)
	payload = append(payload, f.FilterBytes()...)
	payload = append(payload, u.MustEncodeNumLittleEndian(f.FunctionCount)...)
	payload = append(payload, u.MustEncodeNumLittleEndian(f.Tweak)...)
	payload = append(payload, u.MustEncodeNumLittleEndian(flag)...)

	return payload
}
