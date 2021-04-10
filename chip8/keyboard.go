package chip8

type keyboard struct {
	buffer  []byte
	mapping map[rune]byte
}

const buffer_size int = 8

var default_mapping = map[rune]byte{
	'1': 0x1,
	'2': 0x2,
	'3': 0x3,
	'4': 0xc,
	'q': 0x4,
	'w': 0x5,
	'e': 0x6,
	'r': 0xd,
	'a': 0x7,
	's': 0x8,
	'd': 0x9,
	'f': 0xe,
	'z': 0xa,
	'x': 0x0,
	'c': 0xb,
	'v': 0xf,
}

func NewKeyboard() *keyboard {
	return &keyboard{buffer: make([]byte, 0), mapping: default_mapping}
}

func (k *keyboard) addToBuffer(keys []byte) {
	tmp := append(k.buffer, keys...)
	if len(tmp) > buffer_size {
		tmp = tmp[len(tmp)-buffer_size:]
	}

	k.buffer = tmp
}

func (k *keyboard) readKey() (byte, bool) {
	if len(k.buffer) < 1 {
		return 0, false
	}

	ch := rune(k.buffer[0])
	if len(k.buffer) == 1 {
		k.buffer = []byte{}
	} else {
		k.buffer = k.buffer[1:]
	}

	b, ok := k.mapping[ch]

	if !ok {
		return k.readKey()
	}

	return b, true
}
