package tx

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	"github.com/lobiCode/prog_btc_go/script"
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

	nTxIns := u.EncodeVariant(len(tx.TxIns))
	result = append(result, nTxIns...)
	for _, txIn := range tx.TxIns {
		result = append(result, txIn.Serialize()...)
	}
	nTxOuts := u.EncodeVariant(len(tx.TxOuts))
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
	ScriptSig    *script.Script
	Sequence     uint32
	value        uint64
	scriptPubKey *script.Script
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

func ParseTx(r io.Reader, testnet bool) (*Tx, error) {
	b := make([]byte, 4)

	// read version
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	version := binary.LittleEndian.Uint32(b)

	// inTx number
	n, err := u.ReadVariant(r)
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

	n, err = u.ReadVariant(r)
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
