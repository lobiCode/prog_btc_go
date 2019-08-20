package block

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

func TestParseSerialize(t *testing.T) {
	in := "020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := Parse(r)
	check(nil, err, t)

	check(uint32(536870914), result.Version, t)
	check("000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e",
		hex.EncodeToString(result.PrevBlock), t)
	check("be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b",
		hex.EncodeToString(result.MerkleRoot), t)
	check(uint32(1504147230), result.Timestapm, t)
	expBits, _ := hex.DecodeString("e93c0118")
	check(expBits, result.Bits, t)
	expNonce, _ := hex.DecodeString("a4ffd71d")
	check(expNonce, result.Nonce, t)

	check(in, result.Serialize(), t)

	check("0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523",
		result.Hash(), t)

	check(true, result.Bip9(), t)
}
func TestBip9(t *testing.T) {
	in := "0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := Parse(r)
	check(nil, err, t)

	check(false, result.Bip9(), t)
}

func TestDificulty(t *testing.T) {
	bits, _ := hex.DecodeString("e93c0118")
	block := &Block{Bits: bits}
	difficulty := block.Difficulty()
	check("888171856257", difficulty.String(), t)
}

func TestCheckPow(t *testing.T) {
	in := "04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := Parse(r)
	check(nil, err, t)

	check(true, result.CheckPow(), t)
}

func TestCheckPow2(t *testing.T) {
	in := "04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := Parse(r)
	check(nil, err, t)

	check(false, result.CheckPow(), t)
}

func TestCalNewBits(t *testing.T) {
	inFirst := "000000203471101bbda3fe307664b3283a9ef0e97d9a38a7eacd8800000000000000000010c8aba8479bbaa5e0848152fd3c2289ca50e1c3e58c9a4faaafbdf5803c5448ddb845597e8b0118e43a81d3"
	inBFirst, err := hex.DecodeString(inFirst)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inBFirst)
	blockFirst, err := Parse(r)
	check(nil, err, t)

	inLast := "02000020f1472d9db4b563c35f97c428ac903f23b7fc055d1cfc26000000000000000000b3f449fcbe1bc4cfbcb8283a0d2c037f961a3fdf2b8bedc144973735eea707e1264258597e8b0118e5f00474"
	inBLast, err := hex.DecodeString(inLast)
	if err != nil {
		panic(err)
	}

	r = bytes.NewReader(inBLast)
	blockLast, err := Parse(r)
	check(nil, err, t)

	timeDiff := int64(blockLast.Timestapm - blockFirst.Timestapm)

	expected, _ := hex.DecodeString("308d0118")

	result := u.CalculateNewBits(timeDiff, blockLast.Bits)

	check(expected, result, t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
