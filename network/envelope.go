package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

type NetMagic uint32

func (m NetMagic) Encode() []byte {
	return u.MustEncodeNumLittleEndian(m)
}

func ParseNetMagic(r io.Reader) (NetMagic, error) {
	var m NetMagic
	err := u.DecodeInterfaceNumLittleEndian(r, &m)
	return m, err
}

var (
	MainNet  NetMagic = 0xd9b4bef9
	TestNet  NetMagic = 0xdab5bffa
	TestNet3 NetMagic = 0x0709110b
)

type Envelope struct {
	Magic   NetMagic
	Command CommandMsg
	Payload []byte
}

func NewEnvelope(magic NetMagic, message Message) *Envelope {
	return &Envelope{
		Magic:   magic,
		Command: message.GetCommand(),
		Payload: message.Serialize(),
	}
}

func ParseEnvelope(r io.Reader) (*Envelope, error) {
	// TODO
	magic, err := ParseNetMagic(r)
	if err != nil {
		return nil, err
	}
	command, err := ParseCommandMsg(r)
	if err != nil {
		return nil, err
	}

	var payloadLength uint32
	err = u.DecodeInterfaceNumLittleEndian(r, &payloadLength)
	if err != nil {
		return nil, err
	}

	// TODO
	_, err = u.ReadByetes(r, 4)
	if err != nil {
		return nil, err
	}

	payload, err := u.ReadByetes(r, int64(payloadLength))
	if err != nil {

		return nil, err
	}

	return &Envelope{magic, command, payload}, nil
}

func (e *Envelope) serialize() []byte {
	result := make([]byte, 0, 24+len(e.Payload))

	result = append(result, e.Magic.Encode()...)

	result = append(result, e.Command.Encode()...)

	payloadLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLength, uint32(len(e.Payload)))
	result = append(result, payloadLength...)

	result = append(result, u.Hash256(e.Payload)[:4]...)
	result = append(result, e.Payload...)

	return result
}

func (e *Envelope) GetStream() io.Reader {
	return bytes.NewReader(e.Payload)
}

func (e *Envelope) Serialize() string {
	result := e.serialize()

	return hex.EncodeToString(result)
}
