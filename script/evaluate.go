package script

import (
	"math/big"

	u "github.com/lobiCode/prog_btc_go/btcutils"
)

type stack struct {
	s [][]byte
}

func (s *stack) popFirst() []byte {
	if s.length() == 0 {
		return nil
	}

	sf := s.s[0]
	s.s = s.s[1:]

	return sf
}

func (s *stack) pop() []byte {
	l := s.length()
	if l == 0 {
		return nil
	}

	sl := s.s[l-1]
	s.s = s.s[:l-1]

	return sl
}

func (s *stack) get() []byte {
	l := s.length()
	if l == 0 {
		return nil
	}

	e := make([]byte, len(s.s[l-1]))
	copy(e, s.s[l-1])

	return e
}

func (s *stack) getN(i int) []byte {
	l := s.length()
	if l == 0 {
		return nil
	}

	p := 0
	if i < 0 {
		p = l
	}

	p = p + (i)
	if p > l-1 {
		return nil
	}

	e := make([]byte, len(s.s[p]))
	copy(e, s.s[p])

	return e
}

func (s *stack) length() int {
	return len(s.s)
}

func (s *stack) push(b ...[]byte) {
	s.s = append(s.s, b...)
}

func newStack(capacity int) *stack {
	s := make([][]byte, 0, capacity)

	return &stack{s}
}

func Evaluate(z []byte, scriptSig, scriptPubKey *Script) bool {
	cmds := newStack(len(scriptSig.Cmds) + len(scriptPubKey.Cmds))
	cmds.push(scriptSig.Cmds...)
	cmds.push(scriptPubKey.Cmds...)

	realStack := newStack(0)
	altStack := newStack(0)

	return evaluate(u.ParseBytes(z), cmds, realStack, altStack)
}

func evaluate(z *big.Int, cmds, realStack, altStack *stack) bool {
	for cmds.length() > 0 {
		cmd := cmds.popFirst()

		var operationFunc OperationFunc
		if len(cmd) == 1 {
			if s := GetOpCodeName(cmd[0]); s != "" {
				if operationFunc = GetOperationFunction(s); operationFunc == nil {
					return false
				}
			}
		}

		if operationFunc != nil {
			if !operationFunc(z, cmds, realStack, altStack) {
				return false
			}
		} else {
			realStack.push(cmd)
		}
	}

	if realStack.length() == 0 {
		return false
	}

	i, err := decodeNum(realStack.pop())
	if err != nil || i != 1 {
		return false
	}

	return true
}
