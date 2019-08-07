package script

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"strings"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

var ErrTxVersionLen = errors.New("wrong version len")
var ErrScripParse = errors.New("parsing script failed")

type Script struct {
	Cmds [][]byte
}

func (s *Script) Serialize() []byte {
	result := make([]byte, 8)

	var i byte
	var dataLen []byte

	for _, v := range s.Cmds {
		l := len(v)
		if l == 1 && (v[0] == 0 || v[0] > 77) {
			result = append(result, v...)
		} else {
			if l >= 1 && l < 76 {
				dataLen = []byte{}
				i = byte(l)
			} else if l > 75 && l < 0x100 {
				i = 76
				dataLen = []byte{byte(l)}
			} else if l >= 0x100 && l <= 520 {
				i = 77
				dataLen = []byte{0, 0}
				binary.LittleEndian.PutUint16(dataLen, uint16(l))
			} else {
				panic("too long cmd")
			}

			result = append(result, i)
			if len(dataLen) > 0 {
				result = append(result, dataLen...)
			}
			result = append(result, v...)
		}
	}

	variant := u.EncodeVariant(len(result) - 8)
	vl := 8 - len(variant)
	result = result[vl:]
	for i, v := range variant {
		result[i] = v
	}
	return result
}

func (s *Script) IsP2shScriptPubkeys() bool {
	return isP2sh(s.Cmds)
}

func (s *Script) GetRedeemScript() (*Script, error) {
	return Parse(bytes.NewReader(s.Cmds[2]))
}

func (s *Script) String() string {
	outs := []string{}
	for _, v := range s.Cmds {
		if len(v) == 1 {
			if out := GetOpCodeName(v[0]); out != "" {
				outs = append(outs, out)
				continue
			}
		}
		outs = append(outs, hex.EncodeToString(v))
	}
	return strings.Join(outs, " ")
}

func Parse(r io.Reader) (*Script, error) {
	n, err := u.ReadVariant(r)
	if err != nil {
		return nil, err
	}

	cmds := make([][]byte, 0)
	// TODO
	b := make([]byte, 520)

	count := uint64(0)
	for count < n {
		_, err = io.ReadFull(r, b[:1])
		if err != nil {
			return nil, err
		}
		count += 1

		// TODO
		i := uint64(b[0])

		switch {
		case (i >= 1 && i <= 75):
			_, err = io.ReadFull(r, b[:i])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, u.Copyb(b[:i]))
			count += i
		case i == 76:
			_, err = io.ReadFull(r, b[:1])
			if err != nil {
				return nil, err
			}
			dataLen := uint64(b[0])
			_, err = io.ReadFull(r, b[:dataLen])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, u.Copyb(b[:dataLen]))
			count += (dataLen + 1)
		case i == 77:
			_, err = io.ReadFull(r, b[:2])
			if err != nil {
				return nil, err
			}
			dataLen := uint64(binary.LittleEndian.Uint16(b[:2]))
			_, err = io.ReadFull(r, b[:dataLen])
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, u.Copyb(b[:dataLen]))
			count += (dataLen + 2)
		default:
			cmds = append(cmds, u.Copyb(b[:1]))
		}
	}

	if count != n {
		return nil, ErrScripParse
	}

	return &Script{cmds}, nil
}
