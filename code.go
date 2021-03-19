package chip8

import (
	"encoding/binary"
)

type Instructions []byte
type Instruction [2]byte

type Opcode byte

const (
	SYS Opcode = iota
	CLS
	RET
	JP
	CALL
	SE  // Skip if equal to byte
	SNE // Skip in not equal to byte
	SRE // Skip If registers equal
	LD  // Load to register
	OR  // Bitwise OR
	AND
	XOR
	SHR
	SUBN
	SHL
	SRNE // Skip uf registers not equal
	LDI  // Load value to Instruction Pointer register
	JP0  // Jump to value + V0
	RND
	DRW
	SKP
	SKNP
	LDDT  // Load delay timer to register
	LGKP  // Load keypress to register
	LDTDT // Load register to delay time
	LDST  // Load register to sound timer
	ADDI  // Add register to Instruction Pointer
	LDF   // Set Instruction Pointer to location of sprite
	LDB   // Load BCD to I, I+1, I+2
	LDR   // Load registers to memory starting at I
	LDV   // Read memory into registers, starting at I

)

func ReadUint16(ins Instruction) uint16 {
	return binary.BigEndian.Uint16(ins[0:2])
}

// ReadUint8 returns the lowest 8 bits of the instruction
func ReadUint8(ins Instruction) uint8 {
	return uint8(ins[1])
}

// ReadUint12 returns the lowest 12 bits of the instruction
func ReadUint12(ins Instruction) uint16 {
	return ReadUint16(ins) & 0xfff
}

// ReadNibble returns the lowest 4 bits of the instruction
func ReadNibble(ins Instruction) uint8 {
	return uint8(ins[1]) & 0xf
}

func ReadHighByteNibble(ins Instruction) uint8 {
	return uint8(ins[0]) & 0xf
}

func ReadHighByteHighNibble(ins Instruction) uint8 {
	return uint8(ins[0]) >> 4
}

func ReadLowByteHighNibble(ins Instruction) uint8 {
	return uint8(ins[1]) >> 4
}

func ReadHighestByte(ins Instruction) uint8 {
	return uint8(ins[0])
}
