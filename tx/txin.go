package tx

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	"github.com/lobiCode/prog_btc_go/script"
)

type TxIn struct {
	PreTxId      string
	PreTxIdx     uint32
	ScriptSig    *script.Script
	Sequence     uint32
	value        uint64
	scriptPubKey *script.Script
}

func (txIn *TxIn) String() string {
	return fmt.Sprintf(`%s:%d`, txIn.PreTxId, txIn.PreTxIdx)
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

func (txIn *TxIn) serializePreTxId() []byte {
	preTxId, err := hex.DecodeString(txIn.PreTxId)
	u.ReverseBytes(preTxId)
	// TODO
	if err != nil {
		panic(err)
	}
	return preTxId
}

func (txIn *TxIn) serializePreTxIdx() []byte {
	preTxIdx := make([]byte, 4)
	binary.LittleEndian.PutUint32(preTxIdx, txIn.PreTxIdx)
	return preTxIdx
}

func (txIn *TxIn) serializeSequence() []byte {
	/*if txIn.Sequence == 0 {*/
	//return []byte{}
	//}

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, txIn.Sequence)

	return sequence
}

func (txIn *TxIn) SerializeSigHash(replaceScriptSig, testnet bool) ([]byte, error) {
	result := []byte{}
	result = append(result, txIn.serializePreTxId()...)
	result = append(result, txIn.serializePreTxIdx()...)

	if replaceScriptSig {
		scriptPubKey, err := txIn.ScriptPubKey(testnet)
		if err != nil {
			return nil, err
		}
		result = append(result, scriptPubKey.Serialize()...)
	} else {
		result = append(result, 0x00)
	}

	result = append(result, txIn.serializeSequence()...)

	return result, nil
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

func (txIn *TxIn) ScriptPubKey(testnet bool) (*script.Script, error) {

	if txIn.scriptPubKey == nil {
		prevTx, err := FetchTx(txIn.PreTxId, testnet)
		if err != nil {
			return nil, err
		}

		txIn.scriptPubKey = prevTx.TxOuts[txIn.PreTxIdx].ScriptPubKey
	}

	return txIn.scriptPubKey, nil
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

	scriptSig, err := script.Parse(r)
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
