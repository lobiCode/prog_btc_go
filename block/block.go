package block

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

type Block struct {
	Version    uint32
	PrevBlock  []byte
	MerkleRoot []byte
	Timestapm  uint32
	Bits       []byte
	Nonce      []byte
}

func Parse(r io.Reader) (*Block, error) {

	b, err := u.Read(r, 4)
	if err != nil {
		return nil, err
	}
	version := binary.LittleEndian.Uint32(b)

	prevBlock, err := u.Read(r, 32)
	if err != nil {
		return nil, err
	}
	u.ReverseBytes(prevBlock)

	merkleRoot, err := u.Read(r, 32)
	if err != nil {
		return nil, err
	}
	u.ReverseBytes(merkleRoot)

	b, err = u.Read(r, 4)
	if err != nil {
		return nil, err
	}
	timestamp := binary.LittleEndian.Uint32(b)

	bits, err := u.Read(r, 4)
	if err != nil {
		return nil, err
	}

	nonce, err := u.Read(r, 4)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Version:    version,
		PrevBlock:  prevBlock,
		MerkleRoot: merkleRoot,
		Timestapm:  timestamp,
		Bits:       bits,
		Nonce:      nonce,
	}
	return block, nil

}

func (b *Block) serializeVersion() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, b.Version)
	return version
}
func (b *Block) serializeTimestamp() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, b.Timestapm)
	return version
}

func (b *Block) Serialize() string {
	return hex.EncodeToString(b.serialize())
}

func (b *Block) serialize() []byte {
	result := make([]byte, 0, 80)
	result = append(result, b.serializeVersion()...)

	result = append(result, b.PrevBlock...)
	u.ReverseBytes(result[4:36])

	result = append(result, b.MerkleRoot...)
	u.ReverseBytes(result[36:68])

	result = append(result, b.serializeTimestamp()...)

	result = append(result, b.Bits...)
	result = append(result, b.Nonce...)

	return result
}

func (b *Block) Hash() string {
	result := u.Hash256(b.serialize())
	u.ReverseBytes(result)
	return hex.EncodeToString(result)
}

func (b *Block) Bip9() bool {
	return b.Version>>29 == 1
}

func (b *Block) Target() *big.Int {
	return u.BitsToTarget(b.Bits)
}

func (b *Block) Difficulty() *big.Int {
	target := b.Target()
	div := u.MulInt(u.NewInt(65535), u.PowInt(u.NewInt(256), u.NewInt(26)))
	difficulty := u.DivInt(div, target)

	return difficulty
}

func (b *Block) CheckPow() bool {
	target := b.Target()
	proof := u.LittleEndianToBigInt(u.Hash256(b.serialize()))

	return proof.Cmp(target) == -1
}
