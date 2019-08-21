package network

import (
	"encoding/hex"
	"testing"
)

func TestSerializeVesrionMessage(t *testing.T) {
	recvAddr := &NetAddr{
		Services: 0,
		Ip:       []byte{0x00, 0x00, 0x00, 0x00},
		Port:     8333,
	}
	fromAddr := &NetAddr{
		Services: 0,
		Ip:       []byte{0x00, 0x00, 0x00, 0x00},
		Port:     8333,
	}

	vm := &VesrionMessage{
		ProtocolVersion: 70015,
		Services:        0,
		Timestamp:       0,
		RecvAddr:        recvAddr,
		FromAddr:        fromAddr,
		Nonce:           0,
		UserAgent:       "/programmingbitcoin:0.1/",
		Height:          0,
	}
	expected := "7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000"

	check(expected, vm.String(), t)
}

func TestSerializeGetHeadersMsg(t *testing.T) {
	start, err := hex.DecodeString("0000000000000000001237f46acddf58578a37e213d2a6edc4884a2fcad05ba3")
	if err != nil {
		panic(err)
	}

	msg := &GetHeadersMessage{
		ProtocolVersion: 70015,
		StartBlock:      [][]byte{start},
	}

	check("7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		hex.EncodeToString(msg.Serialize()), t)
}
