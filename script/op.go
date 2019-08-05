package script

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
	c "github.com/lobiCode/prog_btc_go/cryptography"
)

type OperationFunc = func(z *big.Int, cmds, realStack, altStack *stack) bool

func _add_number(i int64, realStack *stack) bool {
	b, err := encodeNum(i)
	if err != nil {
		return false
	}

	realStack.push(b)

	return true
}

func op0(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(0, realStack)
}

func op1(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(1, realStack)
}

func op2(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(2, realStack)
}

func op3(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(3, realStack)
}

func op4(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(4, realStack)
}

func op5(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(5, realStack)
}

func op6(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(6, realStack)
}

func op7(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(7, realStack)
}

func op8(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(8, realStack)
}

func op9(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(9, realStack)
}

func op10(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(10, realStack)
}

func op11(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(11, realStack)
}

func op12(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(12, realStack)
}

func op13(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(13, realStack)
}

func op14(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(14, realStack)
}

func op15(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(15, realStack)
}

func op16(z *big.Int, cmds, realStack, altStack *stack) bool {
	return _add_number(16, realStack)
}

func opDup(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	realStack.push(realStack.get())
	return true
}

func opHash256(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	e := realStack.pop()
	realStack.push(u.Hash256(e))

	return true
}

func opHash160(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	e := realStack.pop()
	realStack.push(u.Hash160(e))

	return true
}

func opChecksig(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 2 {
		return false
	}

	publicKeyB := realStack.pop()
	signatureB := realStack.pop()

	publicKey, err := c.ParsePublicKey(publicKeyB)
	if err != nil {
		return false
	}

	signature, err := c.ParseSignature(signatureB[:len(signatureB)-1])
	if err != nil {
		return false
	}

	var i int64 = 0

	if c.Verify(z, signature, publicKey) {
		i = 1
	}

	return _add_number(i, realStack)
}

func opEqual(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 2 {
		return false
	}

	e1 := realStack.pop()
	e2 := realStack.pop()

	ok := bytes.Equal(e1, e2)
	var i int64 = 0
	if ok {
		i = 1
	}

	return _add_number(i, realStack)
}

func opVerify(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	e := realStack.pop()
	i, err := decodeNum(e)
	if err != nil {
		return false
	}

	if i == 0 {
		return false
	}

	return true
}

func opEqualverify(z *big.Int, cmds, realStack, altStack *stack) bool {
	return opEqual(z, cmds, realStack, altStack) && opVerify(z, cmds, realStack, altStack)
}

func op2dup(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 2 {
		return false
	}

	e1 := realStack.getN(-1)
	e2 := realStack.getN(-2)

	realStack.push(e1, e2)

	return true
}

func opSwap(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 2 {
		return false
	}
	e1 := realStack.pop()
	e2 := realStack.pop()

	realStack.push(e1, e2)

	return true
}

func opNot(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}
	e := realStack.pop()

	i, err := decodeNum(e)
	if err != nil {
		return false
	}

	if i == 0 {
		i = 1
	} else {
		i = 0
	}

	return _add_number(i, realStack)
}

func opSha1(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	e := realStack.pop()
	sum := sha1.Sum(e)
	realStack.push(sum[:])

	return true
}

func opCheckmultisig(z *big.Int, cmds, realStack, altStack *stack) bool {
	if realStack.length() < 1 {
		return false
	}

	e := realStack.pop()
	i, err := decodeNum(e)
	if err != nil || int64(realStack.length()) < i+1 || i > 15 {
		return false
	}

	publicKeys := make([][]byte, 0, i)
	for ; i > 0; i-- {
		publicKeys = append(publicKeys, realStack.pop())
	}

	e = realStack.pop()
	i, err = decodeNum(e)
	if err != nil || int64(realStack.length()) < i+1 || i > 15 {
		return false
	}

	sigs := make([][]byte, 0, i)
	for ; i > 0; i-- {
		sigs = append(sigs, realStack.pop())
	}

	realStack.pop()

	if len(sigs) > len(publicKeys) {
		return false
	}

SIG:
	for i, sig := range sigs {
		for j := i; j < len(publicKeys); j++ {
			publicKey, err := c.ParsePublicKey(publicKeys[j])
			if err != nil {
				return false
			}
			signature, err := c.ParseSignature(sig[:len(sig)-1])
			if err != nil {
				return false
			}

			if c.Verify(z, signature, publicKey) {
				break SIG
			}

			return false
		}

	}

	_add_number(1, realStack)

	return true
}

func encodeNum(i int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, i)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func decodeNum(b []byte) (int64, error) {
	var i int64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &i)
	if err != nil {
		return 0, err
	}

	return i, nil
}

var operation_functions = map[string]OperationFunc{
	"OP_DUP":           opDup,
	"OP_HASH256":       opHash256,
	"OP_HASH160":       opHash160,
	"OP_CHECKSIG":      opChecksig,
	"OP_EQUAL":         opEqual,
	"OP_EQUALVERIFY":   opEqualverify,
	"OP_VERIFY":        opVerify,
	"OP_2DUP":          op2dup,
	"OP_SWAP":          opSwap,
	"OP_NOT":           opNot,
	"OP_SHA1":          opSha1,
	"OP_0":             op0,
	"OP_1":             op1,
	"OP_2":             op2,
	"OP_3":             op3,
	"OP_4":             op4,
	"OP_5":             op5,
	"OP_6":             op6,
	"OP_7":             op7,
	"OP_8":             op8,
	"OP_9":             op9,
	"OP_10":            op10,
	"OP_11":            op11,
	"OP_12":            op12,
	"OP_13":            op13,
	"OP_14":            op14,
	"OP_15":            op15,
	"OP_16":            op16,
	"OP_CHECKMULTISIG": opCheckmultisig,
}

var op_codes_names = map[byte]string{
	0:   "OP_0",
	76:  "OP_PUSHDATA1",
	77:  "OP_PUSHDATA2",
	78:  "OP_PUSHDATA4",
	79:  "OP_1NEGATE",
	81:  "OP_1",
	82:  "OP_2",
	83:  "OP_3",
	84:  "OP_4",
	85:  "OP_5",
	86:  "OP_6",
	87:  "OP_7",
	88:  "OP_8",
	89:  "OP_9",
	90:  "OP_10",
	91:  "OP_11",
	92:  "OP_12",
	93:  "OP_13",
	94:  "OP_14",
	95:  "OP_15",
	96:  "OP_16",
	97:  "OP_NOP",
	99:  "OP_IF",
	100: "OP_NOTIF",
	103: "OP_ELSE",
	104: "OP_ENDIF",
	105: "OP_VERIFY",
	106: "OP_RETURN",
	107: "OP_TOALTSTACK",
	108: "OP_FROMALTSTACK",
	109: "OP_2DROP",
	110: "OP_2DUP",
	111: "OP_3DUP",
	112: "OP_2OVER",
	113: "OP_2ROT",
	114: "OP_2SWAP",
	115: "OP_IFDUP",
	116: "OP_DEPTH",
	117: "OP_DROP",
	118: "OP_DUP",
	119: "OP_NIP",
	120: "OP_OVER",
	121: "OP_PICK",
	122: "OP_ROLL",
	123: "OP_ROT",
	124: "OP_SWAP",
	125: "OP_TUCK",
	130: "OP_SIZE",
	135: "OP_EQUAL",
	136: "OP_EQUALVERIFY",
	139: "OP_1ADD",
	140: "OP_1SUB",
	143: "OP_NEGATE",
	144: "OP_ABS",
	145: "OP_NOT",
	146: "OP_0NOTEQUAL",
	147: "OP_ADD",
	148: "OP_SUB",
	149: "OP_MUL",
	154: "OP_BOOLAND",
	155: "OP_BOOLOR",
	156: "OP_NUMEQUAL",
	157: "OP_NUMEQUALVERIFY",
	158: "OP_NUMNOTEQUAL",
	159: "OP_LESSTHAN",
	160: "OP_GREATERTHAN",
	161: "OP_LESSTHANOREQUAL",
	162: "OP_GREATERTHANOREQUAL",
	163: "OP_MIN",
	164: "OP_MAX",
	165: "OP_WITHIN",
	166: "OP_RIPEMD160",
	167: "OP_SHA1",
	168: "OP_SHA256",
	169: "OP_HASH160",
	170: "OP_HASH256",
	171: "OP_CODESEPARATOR",
	172: "OP_CHECKSIG",
	173: "OP_CHECKSIGVERIFY",
	174: "OP_CHECKMULTISIG",
	175: "OP_CHECKMULTISIGVERIFY",
	176: "OP_NOP1",
	177: "OP_CHECKLOCKTIMEVERIFY",
	178: "OP_CHECKSEQUENCEVERIFY",
	179: "OP_NOP4",
	180: "OP_NOP5",
	181: "OP_NOP6",
	182: "OP_NOP7",
	183: "OP_NOP8",
	184: "OP_NOP9",
	185: "OP_NOP10",
}

func GetOpCodeName(b byte) string {
	if c, ok := op_codes_names[b]; ok {
		return c
	}

	return ""
}

func GetOperationFunction(code string) OperationFunc {
	if f, ok := operation_functions[code]; ok {
		return f
	}

	return nil
}
