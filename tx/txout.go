package tx

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lobiCode/prog_btc_go/script"
)

type TxOut struct {
	Amount       uint64
	ScriptPubKey *script.Script
}

func (txOut *TxOut) String() string {
	return fmt.Sprintf(`%d:%s
`, txOut.Amount, txOut.ScriptPubKey)
}
func (txOut *TxOut) Serialize() []byte {
	amount := make([]byte, 8, 34)
	binary.LittleEndian.PutUint64(amount, txOut.Amount)

	scriptPubKey := txOut.ScriptPubKey.Serialize()

	return append(amount, scriptPubKey...)
}

func ParseTxOut(r io.Reader) (*TxOut, error) {
	b := make([]byte, 8)

	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	amount := binary.LittleEndian.Uint64(b)

	scriptPubKey, err := script.Parse(r)
	if err != nil {
		return nil, err
	}

	return &TxOut{amount, scriptPubKey}, nil
}
