package network

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/lobiCode/prog_btc_go/block"
	u "github.com/lobiCode/prog_btc_go/btcutils"
)

var (
	ErrHeadersTxsCount = errors.New("number of txs must be zero")
)

type CommandMsg []byte

func ParseCommandMsg(r io.Reader) (CommandMsg, error) {
	command, err := u.ReadByetes(r, 12)
	if err != nil {
		return nil, err
	}
	command = bytes.TrimFunc(command, u.IsZeroPrefix)

	return CommandMsg(command), nil
}

func (c CommandMsg) Eq(command CommandMsg) bool {
	if bytes.Compare(c, command) == 0 {
		return true
	}
	return false
}

func (c CommandMsg) Encode() []byte {
	command := make([]byte, 12)
	copy(command, c)

	return command
}

var (
	VersionCommand    CommandMsg = []byte("version")
	VerackCommand     CommandMsg = []byte("verack")
	GetHeadersCommand CommandMsg = []byte("getheaders")
	HeadersCommand    CommandMsg = []byte("headers")
	PongCommand       CommandMsg = []byte("pong")
	PingCommand       CommandMsg = []byte("ping")
)

type NetAddr struct {
	Services uint64
	Ip       []byte
	Port     uint16
}

func (addr *NetAddr) Encode() []byte {
	result := make([]byte, 0, 30)

	result = append(result, u.MustEncodeNumLittleEndian(addr.Services)...)

	ip := make([]byte, 10, 16)
	ip = append(ip, 0xff, 0xff)
	ip = append(ip, addr.Ip...)
	result = append(result, ip...)

	result = append(result, u.MustEncodeNumBigEndian(addr.Port)...)

	return result
}

type Message interface {
	Serialize() []byte
	Parse(io.Reader) error
	GetCommand() CommandMsg
}

type VesrionMessage struct {
	ProtocolVersion uint32
	Services        uint64
	Timestamp       int64
	RecvAddr        *NetAddr
	FromAddr        *NetAddr
	Nonce           uint64
	UserAgent       string
	Height          int32
	Relay           bool
}

func (m *VesrionMessage) String() string {
	b := m.Serialize()
	return hex.EncodeToString(b)
}

func (m *VesrionMessage) GetCommand() CommandMsg {
	return VersionCommand
}

func (m *VesrionMessage) Serialize() []byte {
	result := make([]byte, 0, 120)

	result = append(result, u.MustEncodeNumLittleEndian(m.ProtocolVersion)...)

	result = append(result, u.MustEncodeNumLittleEndian(m.Services)...)

	result = append(result, u.MustEncodeNumLittleEndian(m.Timestamp)...)

	result = append(result, m.RecvAddr.Encode()...)
	result = append(result, m.FromAddr.Encode()...)

	result = append(result, u.MustEncodeNumLittleEndian(m.Nonce)...)
	userAgent := []byte(m.UserAgent)
	if len(userAgent) == 0 {
		result = append(result, 0x00)
	} else {
		result = append(result, u.EncodeVariant(len(userAgent))...)
		result = append(result, userAgent...)
	}

	result = append(result, u.MustEncodeNumLittleEndian(m.Height)...)

	if m.Relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result
}

func (m *VesrionMessage) Parse(r io.Reader) error {
	return nil
}

func GetDefaultVersionMessage() *VesrionMessage {
	recvAddr := &NetAddr{
		Services: 128,
		Ip:       []byte{0x00, 0x00, 0x00, 0x00},
		Port:     18333,
	}
	fromAddr := &NetAddr{
		Services: 128,
		Ip:       []byte{0x00, 0x00, 0x00, 0x00},
		Port:     18333,
	}

	rand.Seed(time.Now().Unix())
	nonce := rand.Uint64()
	return &VesrionMessage{
		ProtocolVersion: 70013,
		Services:        128,
		Timestamp:       time.Now().Unix(),
		RecvAddr:        recvAddr,
		FromAddr:        fromAddr,
		Nonce:           nonce,
		UserAgent:       "/kr neki/",
		Height:          0,
	}
}

type VerackMessage struct{}

func (m *VerackMessage) GetCommand() CommandMsg {
	return VerackCommand
}
func (m *VerackMessage) Serialize() []byte {
	return []byte{}
}
func (m *VerackMessage) Parse(r io.Reader) error {
	return nil
}

type GetHeadersMessage struct {
	ProtocolVersion uint32
	StartBlock      [][]byte
	EndBlock        []byte
}

func (m *GetHeadersMessage) GetCommand() CommandMsg {
	return GetHeadersCommand
}

func (m *GetHeadersMessage) Parse(r io.Reader) error {
	return nil
}

func (m *GetHeadersMessage) Serialize() []byte {
	result := make([]byte, 0, 80)
	result = append(result, u.MustEncodeNumLittleEndian(m.ProtocolVersion)...)

	n := u.EncodeVariant(len(m.StartBlock))
	result = append(result, n...)

	for _, s := range m.StartBlock {
		start := u.Copyb(s)
		u.ReverseBytes(start)
		result = append(result, start...)
	}
	if m.EndBlock == nil {
		result = append(result, make([]byte, 32)...)
	} else {
		end := u.Copyb(m.EndBlock)
		u.ReverseBytes(end)
		result = append(result, end...)
	}

	return result
}

type HeadersMessage struct {
	Blocks []*block.Block
}

func (m *HeadersMessage) GetCommand() CommandMsg {
	return HeadersCommand
}

func (m *HeadersMessage) Parse(r io.Reader) error {
	blockCount, err := u.ReadVariant(r)
	if err != nil {
		return err
	}

	if m.Blocks == nil {
		m.Blocks = make([]*block.Block, 0, blockCount)
	}

	for i := uint64(0); i < blockCount; i++ {
		block, err := block.Parse(r)
		if err != nil {
			return err
		}
		m.Blocks = append(m.Blocks, block)

		txCount, err := u.ReadVariant(r)
		if err != nil {
			return err
		}

		if txCount > 0 {
			return ErrHeadersTxsCount
		}
	}

	return nil
}

type PongMessage struct {
	Nonce []byte
}

func (m *PongMessage) GetCommand() CommandMsg {
	return PongCommand
}

func (m *PongMessage) Serialize() []byte {
	return m.Nonce
}

func (m *PongMessage) Parse(r io.Reader) error {
	b, err := u.Read(r, 8)
	if err != nil {
		return err
	}

	m.Nonce = b

	return nil
}
