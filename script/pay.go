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

func P2sh(h160 []byte) *Script {
	cmds := [][]byte{
		{0xa9},
		h160,
		{0x87},
	}

	return &Script{cmds}
}

func P2wpkh(h160 []byte) *Script {
	cmds := [][]byte{
		{0x00},
		h160,
	}

	return &Script{cmds}
}

func P2wsh(h256 []byte) *Script {
	cmds := [][]byte{
		{0x00},
		h256,
	}

	return &Script{cmds}
}
