package chip8

import (
	"io"
)

type keyboard struct {
	buffer  []byte
	input   io.Reader
	mapping map[rune]byte
}

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

func NewKeyboard(input io.Reader) *keyboard {
	return &keyboard{buffer: make([]byte, 1), input: input, mapping: default_mapping}
}

func (k *keyboard) readKey() (byte, bool) {
	n, err := k.input.Read(k.buffer)

	if err == io.EOF {
		return 0, false
	}
	if err != nil {
		panic(err)
	}

	if n != 1 {
		panic("wrong number of bytes read")
	}

	ch := rune(k.buffer[0])

	b, ok := k.mapping[ch]

	if !ok {
		return k.readKey()
	}

	return b, true
}
