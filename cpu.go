package chip8

type cpu struct {
	memory    [4096]byte
	registers [16]byte

	I int16 // Instruction Register

	delay byte // delay timer register
	sound byte //sound timer register

	stack [16]int16

	pc int16 // Program counter
	sp int8  // stack pointer
}

func NewCpu() *cpu {
	return &cpu{
		pc: 0x200, // First 512 bytes are "reserved" for the Chip-8 "interpreter"
	}
}
