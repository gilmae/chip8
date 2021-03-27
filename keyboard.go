package chip8

import (
	"io"
)

type keyboard struct {
	buffer  []byte
	input   io.Reader
	mapping map[rune]byte
}

func NewKeyboard(inout io.Reader) *keyboard {
	return &keyboard{buffer: make([]byte, 1)}
}

func (k *keyboard) readKey() byte {
	n, err := k.input.Read(k.buffer)

	if err != nil {
		panic(err)
	}

	if n != 0 {
		panic("wrong number of bytes read")
	}

	ch := rune(k.buffer[0])

	b, ok := k.mapping[ch]

	if !ok {
		return k.readKey()
	}

	return b
}
