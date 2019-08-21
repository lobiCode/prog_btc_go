package network

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/lobiCode/prog_btc_go/block"
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

func (n *Node) WaitForCommand(command CommandMsg) (*Envelope, error) {
	for {
		envelope, err := n.Read()
		if err != nil {
			return nil, err
		}

		if envelope.Command.Eq(command) {
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

func GetHeaders(firstBlock *block.Block, node *Node) error {
	getHeadersMessage := &GetHeadersMessage{
		ProtocolVersion: 70015,
		StartBlock:      [][]byte{firstBlock.HashBytes()},
	}

	err := node.Send(getHeadersMessage)
	if err != nil {
		return err
	}

	envelope, err := node.WaitForCommand(HeadersCommand)
	if err != nil {
		return err
	}

	headersMessage := &HeadersMessage{}
	err = headersMessage.Parse(envelope.GetStream())
	if err != nil {
		return err
	}

	prevHash := firstBlock.Hash()

	for _, block := range headersMessage.Blocks {
		if !block.CheckPow() {
			return fmt.Errorf("bad proof of work at block %s", block.Hash())
		}

		hash := block.Hash()
		if strings.Compare(prevHash, block.GetPrevBlock()) != 0 {
			return fmt.Errorf("discontinuous at block %s", hash)
		}
		fmt.Println(prevHash, block.GetPrevBlock(), block)
		prevHash = hash

		// TODO
	}

	return nil
}
