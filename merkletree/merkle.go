package merkletree

import (
	"bytes"
	"fmt"
	"io"
	"math"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

type MerkleBlock struct {
	Version    uint32
	PrevBlock  []byte
	MerkleRoot []byte
	Timestapm  uint32
	Bits       []byte
	Nonce      []byte
	TxCount    int32
	Hashes     [][]byte
	Flags      []byte
}

func Parse(r io.Reader) (*MerkleBlock, error) {
	var err error

	var version uint32
	err = u.DecodeInterfaceNumLittleEndian(r, &version)

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

	var timestamp uint32
	err = u.DecodeInterfaceNumLittleEndian(r, &timestamp)

	bits, err := u.Read(r, 4)
	if err != nil {
		return nil, err
	}

	nonce, err := u.Read(r, 4)
	if err != nil {
		return nil, err
	}

	mb := &MerkleBlock{
		Version:    version,
		PrevBlock:  prevBlock,
		MerkleRoot: merkleRoot,
		Timestapm:  timestamp,
		Bits:       bits,
		Nonce:      nonce,
	}

	err = u.DecodeInterfaceNumLittleEndian(r, &mb.TxCount)
	if err != nil {
		return nil, err
	}

	var hashesCount uint64
	hashesCount, err = u.ReadVariant(r)
	if err != nil {
		return nil, err
	}

	for i := uint64(0); i < hashesCount; i++ {
		b, err := u.ReadByetes(r, 32)
		if err != nil {
			return nil, err
		}
		u.ReverseBytes(b)
		mb.Hashes = append(mb.Hashes, b)
	}

	var flagCount uint64
	flagCount, err = u.ReadVariant(r)
	if err != nil {
		return nil, err
	}

	mb.Flags, err = u.ReadByetes(r, int64(flagCount))
	if err != nil {
		return nil, err
	}

	return mb, nil
}

func (mb *MerkleBlock) IsValid() bool {
	mt := NewMerkleTree(mb.TxCount)
	hashes := make([][]byte, 0, len(mb.Hashes))
	for _, v := range mb.Hashes {
		hashes = append(hashes, u.CopybAndReverse(v))
	}

	bits := u.ByteToBits(mb.Flags)
	mt.PopulateTree(bits, hashes)

	return bytes.Compare(u.CopybAndReverse(mt.Root()), mb.MerkleRoot) == 0
}

func MerkleParent(hash1, hash2 []byte) []byte {
	b := make([]byte, 0, len(hash1)+len(hash2))
	b = append(b, hash1...)
	b = append(b, hash2...)

	return u.Hash256(b)
}

func MerkleParentLevel(hashes [][]byte) [][]byte {
	hashesLen := len(hashes)
	newLevel := make([][]byte, 0, hashesLen/2+1)

	for i := 0; i < hashesLen-1; i += 2 {
		newLevel = append(newLevel, MerkleParent(hashes[i], hashes[i+1]))
	}
	if hashesLen&1 == 1 {
		newLevel = append(newLevel, MerkleParent(hashes[hashesLen-1], hashes[hashesLen-1]))
	}

	return newLevel
}

func MerkleRoot(hashes [][]byte) [][]byte {
	if len(hashes) == 1 {
		return hashes
	}

	return MerkleRoot(MerkleParentLevel(hashes))
}

func NewMerkleTree(numHashes int32) *MerkleTree {
	maxDepth := int(math.Ceil(math.Log2(float64(numHashes))))
	tree := make([][][]byte, maxDepth+1)

	for i := 0; i < (maxDepth + 1); i++ {
		tree[i] = make([][]byte, getNumOfLeafs(numHashes, maxDepth, i))
	}

	return &MerkleTree{tree, maxDepth, 0, 0}
}

func getNumOfLeafs(total int32, maxDepth, depth int) int {
	pow := 1 << uint(maxDepth-depth)
	n := float64(total) / float64(pow)
	numLeafs := int(math.Ceil(n))

	return numLeafs
}

type MerkleTree struct {
	Tree         [][][]byte
	MaxDepth     int
	CurrentDepth int
	CurrentIndex int
}

func (m *MerkleTree) PopulateTree(flagBits []byte, hashes [][]byte) {
	bitPosition := 0
	hashPosition := 0
	for m.Root() == nil {
		if m.IsLeaf() {
			m.SetCurrentNode(hashes[hashPosition])
			hashPosition++
			bitPosition++
			m.Up()
		} else {
			leftHash := m.GetLeftNode()
			if leftHash == nil {
				flagBit := flagBits[bitPosition]
				bitPosition++
				if flagBit == 0 {
					m.SetCurrentNode(hashes[hashPosition])
					hashPosition++
					m.Up()
				} else {
					m.Left()
				}
			} else if m.RightExists() {
				rightHash := m.GetRightNode()
				if rightHash == nil {
					m.Right()
				} else {
					m.SetCurrentNode(MerkleParent(leftHash, rightHash))
					m.Up()
				}
			} else {
				m.SetCurrentNode(MerkleParent(leftHash, leftHash))
				m.Up()
			}
		}
	}

	if hashPosition != len(hashes) {
		panic(fmt.Errorf("hashes not all consumed"))
	}
	for ; bitPosition < len(flagBits); bitPosition++ {
		if flagBits[bitPosition] != 0 {
			panic(fmt.Errorf("flag bits not all consumed"))
		}

	}
}

func (m *MerkleTree) Up() {
	m.CurrentDepth -= 1
	m.CurrentIndex /= 2
}

func (m *MerkleTree) Left() {
	m.CurrentDepth += 1
	m.CurrentIndex *= 2
}

func (m *MerkleTree) Right() {
	m.CurrentDepth += 1
	m.CurrentIndex = (m.CurrentIndex * 2) + 1
}

func (m *MerkleTree) SetCurrentNode(hash []byte) {
	m.Tree[m.CurrentDepth][m.CurrentIndex] = hash
}

func (m *MerkleTree) GetCurrentNode() []byte {
	return m.Tree[m.CurrentDepth][m.CurrentIndex]
}

func (m *MerkleTree) GetLeftNode() []byte {
	return m.Tree[m.CurrentDepth+1][m.CurrentIndex*2]
}

func (m *MerkleTree) GetRightNode() []byte {
	return m.Tree[m.CurrentDepth+1][(m.CurrentIndex*2)+1]
}

func (m *MerkleTree) RightExists() bool {
	if m.IsLeaf() {
		return false
	}

	return len(m.Tree[m.CurrentDepth+1]) > (m.CurrentIndex*2)+1
}

func (m *MerkleTree) IsLeaf() bool {
	return m.CurrentDepth == m.MaxDepth
}

func (m *MerkleTree) IsRoot() bool {
	return m.CurrentDepth == 0
}

func (m *MerkleTree) Root() []byte {
	return m.Tree[0][0]
}
