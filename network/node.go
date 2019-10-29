package network

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

type Node struct {
	conn     net.Conn
	address  string
	port     string
	testnet  bool
	netMagic NetMagic
}

func (n *Node) Close() error {
	return n.conn.Close()
}

func (n *Node) Read() (*Envelope, error) {
	envelope, err := ParseEnvelope(n.conn)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}

func (n *Node) Send(message Message) error {
	envelope := NewEnvelope(n.netMagic, message)

	_, err := n.conn.Write(envelope.serialize())
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) WaitForCommand(commands ...CommandMsg) (*Envelope, error) {
	fmt.Println(string(commands[0]))
	for {
		fmt.Println("bbbbbbbbbb")
		envelope, err := n.Read()
		if err != nil {
			fmt.Println("rrrrrrrrrr")
			return nil, err
		}
		fmt.Println(string(envelope.Command), string(commands[0]))

		if envelope.Command.Eqs(commands...) {
			return envelope, nil
		}

		if envelope.Command.Eq(PingCommand) {
			pong := &PongMessage{Nonce: envelope.Payload}
			n.Send(pong)
		}
	}
}

func NewNode(address, port string, testnet bool) (*Node, error) {
	//conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	netMagic := MainNet
	if testnet {
		netMagic = TestNet3
	}

	return &Node{
		conn:     conn,
		address:  address,
		port:     port,
		testnet:  testnet,
		netMagic: netMagic,
	}, nil
}

func Handshake(node *Node) error {
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
		Timestamp:       time.Now().Unix(),
		RecvAddr:        recvAddr,
		FromAddr:        fromAddr,
		Nonce:           0,
		UserAgent:       "/kr neki/",
		Height:          0,
		Relay:           false,
	}

	err := node.Send(vm)
	if err != nil {
		return err
	}

	var version, verack bool
	for !version || !verack {
		envelope, err := ParseEnvelope(node.conn)
		if err != nil {
			return err
		}
		if bytes.Compare(envelope.Command, VersionCommand) == 0 {
			version = true
			// TODO
			fmt.Println("Received version")
		} else if bytes.Compare(envelope.Command, VerackCommand) == 0 {
			fmt.Println("Received verack")
			verack = true
		}
	}

	verackMsg := &VerackMessage{}
	err = node.Send(verackMsg)

	if err != nil {
		return err
	}

	return nil
}

func GetHeaders(firstBlock string, node *Node) (*HeadersMessage, error) {
	b, _ := hex.DecodeString(firstBlock)
	getHeadersMessage := &GetHeadersMessage{
		ProtocolVersion: 70015,
		StartBlock:      [][]byte{b},
	}

	err := node.Send(getHeadersMessage)
	if err != nil {
		return nil, err
	}

	envelope, err := node.WaitForCommand(HeadersCommand)
	if err != nil {
		return nil, err
	}

	headersMessage := &HeadersMessage{}
	err = headersMessage.Parse(envelope.GetStream(), node.testnet)
	if err != nil {
		return nil, err
	}

	prevHash := firstBlock

	for _, block := range headersMessage.Blocks {
		if !block.CheckPow() {
			return nil, fmt.Errorf("bad proof of work at block %s", block.Hash())
		}

		hash := block.Hash()
		if strings.Compare(prevHash, block.GetPrevBlock()) != 0 {
			return nil, fmt.Errorf("discontinuous at block %s", hash)
		}
		prevHash = hash

		// TODO
	}

	return headersMessage, nil
}
