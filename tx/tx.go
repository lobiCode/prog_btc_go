package tx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

var ErrTxVersionLen = errors.New("wrong version len")
var ErrTxScripSig = errors.New("parsing script failed")

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
}

func (tx *Tx) Serialize() string {
	result := make([]byte, 0, 16)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, tx.Version)
	result = append(result, version...)

	nTxIns := encodeVariant(len(tx.TxIns))
	result = append(result, nTxIns...)
	for _, txIn := range tx.TxIns {
		result = append(result, txIn.Serialize()...)
	}
	nTxOuts := encodeVariant(len(tx.TxOuts))
	result = append(result, nTxOuts...)
	for _, txOut := range tx.TxOuts {
		result = append(result, txOut.Serialize()...)
	}
	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.Locktime)
	result = append(result, locktime...)

	return hex.EncodeToString(result)
}

func (tx *Tx) String() string {
	return fmt.Sprintf("version: %d, txins: %s", tx.Version, tx.TxIns)
}

func (tx *Tx) Fee() (uint64, error) {
	var fee uint64 = 0

	for _, v := range tx.TxIns {
		f, err := v.Value(tx.Testnet)
		if err != nil {
			return 0, err
		}
		fee += f
	}
	for _, v := range tx.TxOuts {
		fee -= v.Amount
	}

	return fee, nil
}

type TxIn struct {
	PreTxId      string
	PreTxIdx     uint32
	ScriptSig    *ScriptSig
	Sequence     uint32
	value        uint64
	scriptPubKey []byte
}

func (txIn *TxIn) String() string {
	return fmt.Sprintf("id: %s, preTxIdx: %d, scriptSig: %s, sequence: %d", txIn.PreTxId, txIn.PreTxIdx, txIn.ScriptSig, txIn.Sequence)
}

func (txIn *TxIn) Serialize() []byte {
	preTxId, err := hex.DecodeString(txIn.PreTxId)
	u.ReverseBytes(preTxId)
	// TODO
	if err != nil {
		panic(err)
	}

	preTxIdx := make([]byte, 4)
	binary.LittleEndian.PutUint32(preTxIdx, txIn.PreTxIdx)

	scriptSig := txIn.ScriptSig.Serialize()

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, txIn.Sequence)

	total := 40 + len(scriptSig)
	result := make([]byte, 0, total)
	result = append(result, preTxId...)
	result = append(result, preTxIdx...)
	result = append(result, scriptSig...)
	result = append(result, sequence...)

	return result
}

func (txIn *TxIn) Value(testnet bool) (uint64, error) {

	if txIn.value == 0 {
		prevTx, err := FetchTx(txIn.PreTxId, testnet)
		if err != nil {
			return 0, err
		}

		txIn.value = prevTx.TxOuts[txIn.PreTxIdx].Amount
	}

	return txIn.value, nil
}

func (txIn *TxIn) ScriptPubKey(testnet bool) ([]byte, error) {

	if txIn.scriptPubKey == nil {
		prevTx, err := FetchTx(txIn.PreTxId, testnet)
		if err != nil {
			return nil, err
		}

		txIn.scriptPubKey = prevTx.TxOuts[txIn.PreTxIdx].ScriptPubKey
	}

	return txIn.scriptPubKey, nil
}

type TxOut struct {
	Amount       uint64
	ScriptPubKey []byte
}

func (txOut *TxOut) Serialize() []byte {
	amount := make([]byte, 8, 34)
	binary.LittleEndian.PutUint64(amount, txOut.Amount)

	scriptPubKey := copyb(txOut.ScriptPubKey)
	u.ReverseBytes(scriptPubKey)

	return append(amount, scriptPubKey...)
}

type ScriptSig struct {
	Cmds [][]byte
}

func (s *ScriptSig) Serialize() []byte {
	result := make([]byte, 8)

	var i byte
	var dataLen []byte
	for _, v := range s.Cmds {
		l := len(v)
		if l >= 1 && l < 76 {
			dataLen = []byte{}
			i = byte(l)
		} else if l > 75 && l < 0x100 {
			i = 76
			dataLen = []byte{byte(l)}
		} else if l >= 0x100 && l <= 520 {
			i = 77
			dataLen = []byte{0, 0}
			binary.LittleEndian.PutUint16(dataLen, uint16(l))
		} else {
			panic("too long cmd")
		}

		result = append(result, i)
		if len(dataLen) > 0 {
			result = append(result, dataLen...)
		}
		result = append(result, v...)
	}

	variant := encodeVariant(len(result) - 8)
	vl := 8 - len(variant)
	result = result[vl:]
	for i, v := range variant {
		result[i] = v
	}

	return result
}

func (s *ScriptSig) String() string {
	return fmt.Sprintf("cmds: %b", s.Cmds)
}

func ParseTx(r io.Reader, testnet bool) (*Tx, error) {
	b := make([]byte, 4)

	// read version
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	version := binary.LittleEndian.Uint32(b)

	// inTx number
	n, err := readVariant(r)
	if err != nil {
		return nil, err
	}

	txIns := []*TxIn{}
	for i := uint64(0); i < n; i++ {
		txIn, err := ParseTxIn(r)
		if err != nil {
			return nil, err
		}
		txIns = append(txIns, txIn)
	}

	n, err = readVariant(r)
	if err != nil {
		return nil, err
	}

	txOuts := make([]*TxOut, 0, n)

	for i := uint64(0); i < n; i++ {
		txOut, err := ParseTxOut(r)
		if err != nil {
			return nil, err
		}
		txOuts = append(txOuts, txOut)
	}

	// read locktime
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	locktime := binary.LittleEndian.Uint32(b)

	tx := &Tx{version, txIns, txOuts, locktime, testnet}

	return tx, nil
}

func ParseTxIn(r io.Reader) (*TxIn, error) {
	b := make([]byte, 32)

	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	u.ReverseBytes(b)
	preTxId := hex.EncodeToString(b)

	_, err = io.ReadFull(r, b[:4])
	if err != nil {
		return nil, err
	}
	preTxIdx := binary.LittleEndian.Uint32(b[:4])

	scriptSig, err := ParseScriptSig(r)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(r, b[:4])
	if err != nil {
		return nil, err
	}
	sequence := binary.LittleEndian.Uint32(b[:4])

	txIn := &TxIn{
		PreTxId:   preTxId,
		PreTxIdx:  preTxIdx,
		ScriptSig: scriptSig,
		Sequence:  sequence,
	}

	return txIn, nil
}

func ParseTxOut(r io.Reader) (*TxOut, error) {
	b := make([]byte, 26)

	_, err := io.ReadFull(r, b[:8])
	if err != nil {
		return nil, err
	}

	amount := binary.LittleEndian.Uint64(b[:8])

	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	u.ReverseBytes(b)
	scriptPubKey := copyb(b)

	return &TxOut{amount, scriptPubKey}, nil
}

func ParseScriptSig(r io.Reader) (*ScriptSig, error) {
	n, err := readVariant(r)
	if err != nil {
		return nil, err
	}

	cmds := make([][]byte, 0)
	// TODO
	b := make([]byte, 520)

	count := uint64(0)
	for count < n {
		_, err = io.ReadFull(r, b[:1])
		if err != nil {
			return nil, err
		}
		count += 1

		// TODO
		i := uint64(b[0])
		switch {
		case (i >= 1 && i <= 75):
			_, err = io.ReadFull(r, b[:i])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, copyb(b[:i]))
			count += i
		case i == 76:
			_, err = io.ReadFull(r, b[:1])
			if err != nil {
				return nil, err
			}
			dataLen := uint64(b[0])
			_, err = io.ReadFull(r, b[:dataLen])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, copyb(b[:dataLen]))
			count += (dataLen + 1)
		case i == 77:
			_, err = io.ReadFull(r, b[:2])
			if err != nil {
				return nil, err
			}
			dataLen := uint64(binary.LittleEndian.Uint16(b[:2]))
			_, err = io.ReadFull(r, b[:dataLen])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, copyb(b[:dataLen]))
			count += (dataLen + 2)
		default:
			cmds = append(cmds, copyb(b[:0]))
		}
	}

	if count != n {
		return nil, ErrTxScripSig
	}

	return &ScriptSig{cmds}, nil
}

func copyb(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

func readVariant(r io.Reader) (uint64, error) {
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

func encodeVariant(i int) []byte {
	switch {
	case i < 0xfd:
		return []byte{byte(i)}
	case i < 0x10000:
		return _encodeVariant(0xfd, uint16(i))
	case i < 0x100000000:
		return _encodeVariant(0xfe, uint32(i))
	case i < 0x7FFFFFFFFFFFFFFF:
		return _encodeVariant(0xff, uint64(i))
	default:
		panic("integer to large")
	}
}

func _encodeVariant(b byte, i interface{}) []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte(b)

	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func FetchTx(txId string, testnet bool) (*Tx, error) {
	var url string
	// TODO
	if testnet {
		url = fmt.Sprintf("https://blockchain.info/rawtx/%s?format=hex", txId)
	} else {
		url = fmt.Sprintf("https://blockchain.info/rawtx/%s?format=hex", txId)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r := hex.NewDecoder(resp.Body)
	tx, err := ParseTx(r, testnet)
	return tx, err
}
