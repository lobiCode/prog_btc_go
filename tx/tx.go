package tx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	c "github.com/lobiCode/prog_btc_go/cryptography"
	"github.com/lobiCode/prog_btc_go/script"
)

var ErrTxVersionLen = errors.New("wrong version len")
var ErrTxScripSig = errors.New("parsing script failed")
var ErrWrongSig = errors.New("wrong signature")

var (
	SIGHASH_ALL    uint32 = 1
	SIGHASH_NONE   uint32 = 2
	SIGHASH_SINGLE uint32 = 3
)

type Tx struct {
	Version  uint32
	TxIns    []*TxIn
	TxOuts   []*TxOut
	Locktime uint32
	Testnet  bool
}

func (tx *Tx) serializeVersion() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, tx.Version)
	return version
}

func (tx *Tx) serializeNTxIns() []byte {
	return u.EncodeVariant(len(tx.TxIns))
}

func (tx *Tx) serializeNTxOuts() []byte {
	return u.EncodeVariant(len(tx.TxOuts))
}

func (tx *Tx) serializeLocktime() []byte {
	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.Locktime)

	return locktime
}

func (tx *Tx) Serialize() string {
	result := make([]byte, 0, 16)

	result = append(result, tx.serializeVersion()...)

	result = append(result, tx.serializeNTxIns()...)
	for _, txIn := range tx.TxIns {
		result = append(result, txIn.Serialize()...)
	}
	result = append(result, tx.serializeNTxOuts()...)
	for _, txOut := range tx.TxOuts {
		result = append(result, txOut.Serialize()...)
	}
	result = append(result, tx.serializeLocktime()...)

	return hex.EncodeToString(result)
}

func (tx *Tx) String() string {
	return fmt.Sprintf(
		`version: %d
txins:
%s
txouts:
%s`,
		tx.Version, tx.TxIns, tx.TxOuts)
}

func (tx *Tx) SigHash(replaceScriptSig int, redeemScript *script.Script) ([]byte, error) {
	result := []byte{}

	result = append(result, tx.serializeVersion()...)
	result = append(result, tx.serializeNTxIns()...)
	for i, txIn := range tx.TxIns {
		var ok bool
		if i == replaceScriptSig {
			ok = true
		}
		txInSer, err := txIn.SerializeSigHash(ok, tx.Testnet, redeemScript)
		if err != nil {
			return nil, err
		}
		result = append(result, txInSer...)
	}

	result = append(result, tx.serializeNTxOuts()...)
	for _, txOut := range tx.TxOuts {
		result = append(result, txOut.Serialize()...)
	}
	result = append(result, tx.serializeLocktime()...)

	sigHashAll := make([]byte, 4)
	binary.LittleEndian.PutUint32(sigHashAll, SIGHASH_ALL)

	result = append(result, sigHashAll...)

	r := u.Hash256(result)[:]
	return r, nil
}

func (tx *Tx) Verify() bool {
	fee, err := tx.Fee()
	if fee < 0 && err != nil {
		return false
	}

	for i, _ := range tx.TxIns {
		if !tx.verifyInput(i) {
			return false
		}
	}

	return true
}

func (tx *Tx) getReedemScript(replaceScriptSig int) (*script.Script, error) {
	return nil, nil
}

func (tx *Tx) verifyInput(replaceScriptSig int) bool {

	scriptPubKey, err := tx.TxIns[replaceScriptSig].ScriptPubKey(tx.Testnet)
	if err != nil {
		return false
	}

	var redeemScript *script.Script
	if scriptPubKey.IsP2shScriptPubkeys() {
		r, err := tx.TxIns[replaceScriptSig].ScriptSig.GetRedeemScript()
		if err != nil {
			return false
		}
		redeemScript = r
	}

	z, err := tx.SigHash(replaceScriptSig, redeemScript)
	if err != nil {
		return false
	}

	ok := script.Evaluate(z, tx.TxIns[replaceScriptSig].ScriptSig, scriptPubKey)
	return ok
}

func (tx *Tx) SingInput(i int, key *c.PrivateKey) error {
	redeemScript, err := tx.getReedemScript(i)

	sec := key.Sec(true)
	v := tx.TxIns[i]
	z, err := tx.SigHash(i, redeemScript)
	if err != nil {
		return err
	}
	zi := u.ParseBytes(z)
	sigd := key.Sign(zi)
	sig := sigd.Der()
	sig = append(sig, byte(SIGHASH_ALL))
	scriptSig := &script.Script{Cmds: [][]byte{sig, sec}}
	v.ScriptSig = scriptSig

	ok := tx.verifyInput(i)
	if !ok {
		return ErrWrongSig

	}
	return nil
}

func (tx *Tx) SingInputs(key *c.PrivateKey) error {
	var err error
	for i, _ := range tx.TxIns {
		if err = tx.SingInput(i, key); err != nil {
			return err
		}
	}
	return nil
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

func (tx *Tx) IsCoinbase() bool {
	if len(tx.TxIns) != 1 {
		return false
	}

	if tx.TxIns[0].PreTxId != "0000000000000000000000000000000000000000000000000000000000000000" {
		return false
	}

	if tx.TxIns[0].PreTxIdx != 0xffffffff {
		return false
	}

	return true
}

func (tx *Tx) CoinbaseHeight() (int64, error) {
	b := tx.TxIns[0].ScriptSig.Cmds[0]
	return u.DecodeNumLittleEndian(b)
}

func FetchTx(txId string, testnet bool) (*Tx, error) {
	var url string
	// TODO
	if testnet {
		url = fmt.Sprintf("http://testnet.programmingbitcoin.com/tx/%s.hex", txId)
	} else {
		url = fmt.Sprintf("https://blockchain.info/rawtx/%s?format=hex", txId)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	b = bytes.Trim(b, "\n")
	i, err := hex.Decode(b, b)
	if err != nil {
		return nil, err
	}

	b = b[:i]

	var z bool
	if b[4] == 0 {
		b = append(b[:4], b[6:]...)
		z = true
	}

	r := bytes.NewReader(b)
	tx, err := ParseTx(r, testnet)
	if err != nil {
		return nil, err
	}

	if z {
		locktime := binary.LittleEndian.Uint32(b[len(b)-4:])
		tx.Locktime = locktime
	}

	// TODO check tr id

	return tx, err
}
