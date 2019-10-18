package merkletree

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

func TestMerkleParent(t *testing.T) {
	hash1, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	hash2, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	expected, _ := hex.DecodeString("8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd'")

	result := MerkleParent(hash1, hash2)
	check(expected, result, t)

}

func TestMerkleParentLevelOdd(t *testing.T) {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
	}
	hashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		hashes = append(hashes, b)
	}

	expectedHexHashes := []string{
		"8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd",
		"7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800",
		"ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7",
		"68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069",
		"43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27",
		"1796cd3ca4fef00236e07b723d3ed88e1ac433acaaa21da64c4b33c946cf3d10",
	}
	expected := make([][]byte, 0, len(expectedHexHashes))
	for _, v := range expectedHexHashes {
		b, _ := hex.DecodeString(v)
		expected = append(expected, b)
	}

	result := MerkleParentLevel(hashes)
	check(expected, result, t)
}

func TestMerkleParentLevelEven(t *testing.T) {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
	}
	hashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		hashes = append(hashes, b)
	}

	expectedHexHashes := []string{
		"8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd",
		"7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800",
		"ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7",
		"68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069",
		"43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27",
	}
	expected := make([][]byte, 0, len(expectedHexHashes))
	for _, v := range expectedHexHashes {
		b, _ := hex.DecodeString(v)
		expected = append(expected, b)
	}

	result := MerkleParentLevel(hashes)
	check(expected, result, t)
}

func TestMerkleRoot(t *testing.T) {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
		"b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0",
	}
	hashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		hashes = append(hashes, b)
	}

	expectedHash, _ := hex.DecodeString("acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6")
	expected := [][]byte{expectedHash}

	result := MerkleRoot(hashes)

	check(expected, result, t)
}

func TestMerkleParse(t *testing.T) {
	merkleBlock := "00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635"
	inB, _ := hex.DecodeString(merkleBlock)
	r := bytes.NewReader(inB)
	mb, err := Parse(r)
	check(err, nil, t)
	expectedVersion := uint32(0x20000000)
	check(expectedVersion, mb.Version, t)

	merkleRootHex := "ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4"
	expectedB, _ := hex.DecodeString(merkleRootHex)
	u.ReverseBytes(expectedB)
	check(expectedB, mb.MerkleRoot, t)
	prevBlockHex := "df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000"
	expectedB, _ = hex.DecodeString(prevBlockHex)
	u.ReverseBytes(expectedB)
	check(expectedB, mb.PrevBlock, t)

	expectedTimestamp := uint32(1535343836)
	check(expectedTimestamp, mb.Timestapm, t)
	expectedB, _ = hex.DecodeString("67d8001a")
	check(expectedB, mb.Bits, t)
	expectedB, _ = hex.DecodeString("c157e670")
	check(expectedB, mb.Nonce, t)
	expectedTotal := int32(3519)
	check(expectedTotal, mb.TxCount, t)
	hexHashes := []string{
		"ba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a",
		"7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d",
		"34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2",
		"158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cba",
		"ee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763ce",
		"f8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097",
		"c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d",
		"6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543",
		"d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274c",
		"dfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb62261",
	}

	expectedHashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		u.ReverseBytes(b)
		expectedHashes = append(expectedHashes, b)
	}
	check(expectedHashes, mb.Hashes, t)

	expectedB, _ = hex.DecodeString("b55635")
	check(expectedB, mb.Flags, t)
}

func TestMerkleBlockIsValid(t *testing.T) {
	merkleBlock := "00000020df3b053dc46f162a9b00c7f0d5124e2676d47bbe7c5d0793a500000000000000ef445fef2ed495c275892206ca533e7411907971013ab83e3b47bd0d692d14d4dc7c835b67d8001ac157e670bf0d00000aba412a0d1480e370173072c9562becffe87aa661c1e4a6dbc305d38ec5dc088a7cf92e6458aca7b32edae818f9c2c98c37e06bf72ae0ce80649a38655ee1e27d34d9421d940b16732f24b94023e9d572a7f9ab8023434a4feb532d2adfc8c2c2158785d1bd04eb99df2e86c54bc13e139862897217400def5d72c280222c4cbaee7261831e1550dbb8fa82853e9fe506fc5fda3f7b919d8fe74b6282f92763cef8e625f977af7c8619c32a369b832bc2d051ecd9c73c51e76370ceabd4f25097c256597fa898d404ed53425de608ac6bfe426f6e2bb457f1c554866eb69dcb8d6bf6f880e9a59b3cd053e6c7060eeacaacf4dac6697dac20e4bd3f38a2ea2543d1ab7953e3430790a9f81e1c67f5b58c825acf46bd02848384eebe9af917274cdfbb1a28a5d58a23a17977def0de10d644258d9c54f886d47d293a411cb6226103b55635"
	inB, _ := hex.DecodeString(merkleBlock)
	r := bytes.NewReader(inB)
	mb, err := Parse(r)
	check(err, nil, t)

	check(true, mb.IsValid(), t)
}

func TestMerkle1(t *testing.T) {
	hexHashes := []string{
		"9745f7173ef14ee4155722d1cbf13304339fd00d900b759c6f9d58579b5765fb",
		"5573c8ede34936c29cdfdfe743f7f5fdfbd4f54ba0705259e62f39917065cb9b",
		"82a02ecbb6623b4274dfcab82b336dc017a27136e08521091e443e62582e8f05",
		"507ccae5ed9b340363a0e6d765af148be9cb1c8766ccc922f83e4ae681658308",
		"a7a4aec28e7162e1e9ef33dfa30f0bc0526e6cf4b11a576f6c5de58593898330",
		"bb6267664bd833fd9fc82582853ab144fece26b7a8a5bf328f8a059445b59add",
		"ea6d7ac1ee77fbacee58fc717b990c4fcccf1b19af43103c090f601677fd8836",
		"457743861de496c429912558a106b810b0507975a49773228aa788df40730d41",
		"7688029288efc9e9a0011c960a6ed9e5466581abf3e3a6c26ee317461add619a",
		"b1ae7f15836cb2286cdd4e2c37bf9bb7da0a2846d06867a429f654b2e7f383c9",
		"9b74f89fa3f93e71ff2c241f32945d877281a6a50a6bf94adac002980aafe5ab",
		"b3a92b5b255019bdaf754875633c2de9fec2ab03e6b8ce669d07cb5b18804638",
		"b5c0b915312b9bdaedd2b86aa2d0f8feffc73a2d37668fd9010179261e25e263",
		"c9d52c5cb1e557b92c84c52e7c4bfbce859408bedffc8a5560fd6e35e10b8800",
		"c555bc5fc3bc096df0a0c9532f07640bfb76bfe4fc1ace214b8b228a1297a4c2",
		"f9dbfafc3af3400954975da24eb325e326960a25b87fffe23eef3e7ed2fb610e",
	}
	mt := NewMerkleTree(int32(len(hexHashes)))
	hashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		hashes = append(hashes, b)
	}
	flagBits := make([]byte, 31)
	for i, _ := range flagBits {
		flagBits[i] = 1
	}

	mt.PopulateTree(flagBits, hashes)
	expectedRoot, _ := hex.DecodeString("597c4bafe3832b17cbbabe56f878f4fc2ad0f6a402cee7fa851a9cb205f87ed1")
	check(expectedRoot, mt.Root(), t)
}

func TestMerkle2(t *testing.T) {
	hexHashes := []string{
		"42f6f52f17620653dcc909e58bb352e0bd4bd1381e2955d19c00959a22122b2e",
		"94c3af34b9667bf787e1c6a0a009201589755d01d02fe2877cc69b929d2418d4",
		"959428d7c48113cb9149d0566bde3d46e98cf028053c522b8fa8f735241aa953",
		"a9f27b99d5d108dede755710d4a1ffa2c74af70b4ca71726fa57d68454e609a2",
		"62af110031e29de1efcad103b3ad4bec7bdcf6cb9c9f4afdd586981795516577",
	}
	mt := NewMerkleTree(int32(len(hexHashes)))
	hashes := make([][]byte, 0, len(hexHashes))
	for _, v := range hexHashes {
		b, _ := hex.DecodeString(v)
		hashes = append(hashes, b)
	}
	flagBits := make([]byte, 11)
	for i, _ := range flagBits {
		flagBits[i] = 1
	}

	mt.PopulateTree(flagBits, hashes)
	expectedRoot, _ := hex.DecodeString("a8e8bd023169b81bc56854137a135b97ef47a6a7237f4c6e037baed16285a5ab")
	check(expectedRoot, mt.Root(), t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
