package bloom

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	bf := NewFilter(10, 5, 99)
	item := []byte("Hello World")
	bf.Add(item)
	expected := "0000000a080000000140"
	result := bf.FilterBytes()
	resultS := hex.EncodeToString(result)
	check(expected, resultS, t)
	item = []byte("Goodbye!")
	bf.Add(item)
	expected = "4000600a080000010940"
	result = bf.FilterBytes()
	resultS = hex.EncodeToString(result)
	check(expected, resultS, t)
}

func TestFilterLoad(t *testing.T) {
	bf := NewFilter(10, 5, 99)
	item := []byte("Hello World")
	bf.Add(item)
	expected := "0000000a080000000140"
	result := bf.FilterBytes()
	resultS := hex.EncodeToString(result)
	check(expected, resultS, t)
	item = []byte("Goodbye!")
	bf.Add(item)
	expected = "0a4000600a080000010940050000006300000001"
	result = bf.FilterLoad(1)
	resultS = hex.EncodeToString(result)
	check(expected, resultS, t)
}

func check(expected, recived interface{}, t *testing.T) {
	t.Helper()
	if !reflect.DeepEqual(recived, expected) {
		t.Errorf("Received\n%+v\ndoesn't match expected\n%+v\n", recived, expected)
	}
}
