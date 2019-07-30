package main

import (
	"fmt"

	c "github.com/lobiCode/prog_btc_go/cryptography"
	"github.com/lobiCode/prog_btc_go/script"
	"github.com/lobiCode/prog_btc_go/tx"
)

func main() {
	secret := c.GetHash256Int("")
	key := c.NewPrivateKey(secret)
	fmt.Println(key.Address(true, true))

	targetH160, err := c.GetH160Address("miKegze5FQNCnGw6PKyqUbYUeBa4x2hFeM")
	targetH1602, err := c.GetH160Address("2MtwTo5PCjTiGdKHfVVWFp4HGEdRk1TmZ9K")
	targetH1603, err := c.GetH160Address(key.Address(true, true))
	if err != nil {
		panic(err)
	}

	txIn1 := &tx.TxIn{
		PreTxId:  "5b27a7941a230bd74848694950b3e356360ac1240f0e786674c2756c3390194c",
		PreTxIdx: 1,
		Sequence: 0,
	}
	txIn2 := &tx.TxIn{
		PreTxId:  "92004eac6d9e29d91b63db4151c51b507651efc45cb5da22484e4c32f1dde4fd",
		PreTxIdx: 0,
		Sequence: 0,
	}
	txOut1 := &tx.TxOut{
		Amount:       uint64(0.00000023 * 100000000),
		ScriptPubKey: script.P2pkh(targetH160),
	}
	txOut2 := &tx.TxOut{
		Amount:       uint64(0.00000023 * 100000000),
		ScriptPubKey: script.P2pkh(targetH1602),
	}
	txOut3 := &tx.TxOut{
		Amount:       uint64(0.04115 * 100000000),
		ScriptPubKey: script.P2pkh(targetH1603),
	}

	transaction := &tx.Tx{
		Version:  1,
		Locktime: 0,
		Testnet:  true,
		TxIns:    []*tx.TxIn{txIn1, txIn2},
		TxOuts:   []*tx.TxOut{txOut1, txOut2, txOut3},
	}

	err = transaction.SingInput(0, key)
	err = transaction.SingInput(1, key)

	if err != nil {
		panic(err)
	}
	fmt.Println(transaction.Serialize())
}
