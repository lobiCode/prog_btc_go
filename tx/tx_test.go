package tx

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseTx(t *testing.T) {
	in := "010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600"

	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseTx(r, true)

	// check error
	check(nil, err, t)

	// cehck version
	check(uint32(1), result.Version, t)

	// check txIns
	txIns := []*TxIn{
		{PreTxId: "9e067aedc661fca148e13953df75f8ca6eada9ce3b3d8d68631769ac60999156", PreTxIdx: 1},
		{PreTxId: "d37f9e7282f81b7fd3af0fde8b462a1c28024f1d83cf13637ec18d03f4518feb", PreTxIdx: 0},
		{PreTxId: "75d7454b7010fa28b00f16cccb640b1756fd6e357c03a3b81b9d119505f47b56", PreTxIdx: 0},
		{PreTxId: "45f3f79066d251addc04fd889f776c73afab1cb22559376ff820e6166c5e3ad6", PreTxIdx: 1},
	}
	check(len(txIns), len(result.TxIns), t)
	for i, v := range txIns {
		check(v.PreTxId, result.TxIns[i].PreTxId, t)
		check(v.PreTxIdx, result.TxIns[i].PreTxIdx, t)
	}

	txOuts := []*TxOut{
		{Amount: 1000273},
		{Amount: 40000000},
	}

	check(len(txOuts), len(result.TxOuts), t)
	for i, v := range txOuts {
		check(v.Amount, result.TxOuts[i].Amount, t)
	}

	s := result.Serialize()
	check(in, s, t)
}

func TestParseTxIn(t *testing.T) {
	in := "56919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffff"

	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseTxIn(r)
	if err != nil {
		panic(err)
	}

	s := hex.EncodeToString(result.Serialize())
	check(in, s, t)
}
func TestParseTxOut(t *testing.T) {
	in := "51430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac"

	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseTxOut(r)
	if err != nil {
		panic(err)
	}

	s := hex.EncodeToString(result.Serialize())
	check(in, s, t)
}

func TestFee(t *testing.T) {
	in := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseTx(r, true)
	if err != nil {
		panic(err)
	}

	fee, err := result.Fee()
	check(nil, err, t)
	check(uint64(40000), fee, t)
}

func TestLocktime(t *testing.T) {
	in := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseTx(r, true)
	if err != nil {
		panic(err)
	}

	check(uint32(410393), result.Locktime, t)
}

/*func TestFetchTx(t *testing.T) {*/
//id := "08e16d81608810565dcaaa8c4ee61b8c840fe59762c3cb01fd21a90cc71d96b3"
//tx, err := FetchTx(id, true)
//fmt.Println(tx, err)
/*}*/

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
