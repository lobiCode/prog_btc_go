package script

func P2pkh(h160 []byte) *Script {
	cmds := [][]byte{
		{0x76},
		{0xa9},
		h160,
		{0x88},
		{0xac},
	}

	return &Script{cmds}
}
