package network

import (
	"bytes"
	"encoding/hex"
	"reflect"
	"testing"
)

func TestParseEnvelope(t *testing.T) {
	in := "f9beb4d976657261636b000000000000000000005df6e0e2"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseEnvelope(r)
	check(nil, err, t)

	check(VerackCommand, result.Command, t)

	check([]byte{}, result.Payload, t)
}

func TestParseEnvelope2(t *testing.T) {
	in := "f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001"
	inB, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(inB)
	result, err := ParseEnvelope(r)
	check(nil, err, t)

	check(VersionCommand, result.Command, t)

	check(inB[24:], result.Payload, t)

	check(in, result.Serialize(), t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
