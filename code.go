package chip8

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

type Definition struct {
	Name          string
	OperandWidths []int // bits
}

const (
	UNKNOWN Opcode = iota
	SYS
	CLS
	RET
	JP
	CALL
	// SE  // Skip if equal to byte
	// SNE // Skip in not equal to byte
	// SRE // Skip If registers equal
	// LD  // Load to register
	// ADD
	// LDR
	// OR // Bitwise OR
	// AND
	// XOR
	// ADDR
	// SUB
	// SHR
	// SUBN
	// SHL
	// SRNE // Skip uf registers not equal
	// LDI  // Load value to Instruction Pointer register
	// JP0  // Jump to value + V0
	// RND
	// DRW
	// SKP
	// SKNP
	// LDDT  // Load delay timer to register
	// LGKP  // Load keypress to register
	// LDTDT // Load register to delay time
	// LDST  // Load register to sound timer
	// ADDI  // Add register to Instruction Pointer
	// LDF   // Set Instruction Pointer to location of sprite
	// LDB   // Load BCD to I, I+1, I+2
	// LDR   // Load registers to memory starting at I
	// LDV   // Read memory into registers, starting at I
)

var definitions = map[Opcode]*Definition{
	SYS:  {"SYS", []int{12}},
	CLS:  {"CLS", []int{}},
	RET:  {"RET", []int{}},
	JP:   {"JP", []int{12}},
	CALL: {"CALL", []int{12}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

func ParseOpcode(ins Instructions) Opcode {
	switch ReadHighByteHighNibble(ins) {
	case 0x0:
		tribble := ReadUint12(ins)
		switch tribble {
		case 0xE0: // 00E0 - CLS
			return CLS
		case 0xEE: // 00EE - RET
			return RET
		default: // 0nnn - SYS addr. Presumably will be unuused
			return SYS
		}
	case 0x1:
		return JP
	case 0x2:
		return CALL
		// case 0x3:
		// 	return SE
		// case 0x4:
		// 	return SNE
		// case 0x5:
		// 	return SRE
		// case 0x6:
		// 	return LD
		// case 0x7:
		// 	return ADD
		// case 0x8:
		// 	nibble := ReadNibble(ins)
		// 	switch nibble {
		// 	case 0x0:
		// 		return LDR
		// 	case 0x1:
		// 		return OR
		// 	case 0x2:
		// 		return AND
		// 	case 0x3:
		// 		return XOR
		// 	case 0x4:

		// 	}
	}
	return UNKNOWN
}

func ReadOperands(def *Definition, ins Instructions) []int {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 12:
			operands[i] = int(ReadUint12((ins)))
		case 8:
			operands[i] = int(ReadUint8(ins))
		case 4:
			switch offset {
			case 0:
				operands[i] = int(ReadHighByteNibble(ins))
			case 1:
				operands[i] = int(ReadLowByteHighNibble(ins))
			case 2:
				operands[i] = int(ReadNibble(ins))
			}
		}
		offset++
	}

	return operands
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins[0:2])
}

// ReadUint8 returns the lowest 8 bits of the instruction
func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[1])
}

// ReadUint12 returns the lowest 12 bits of the instruction
func ReadUint12(ins Instructions) uint16 {
	return ReadUint16(ins) & 0xfff
}

// ReadNibble returns the lowest 4 bits of the instruction
func ReadNibble(ins Instructions) uint8 {
	return uint8(ins[1]) & 0xf
}

func ReadHighByteNibble(ins Instructions) uint8 {
	return uint8(ins[0]) & 0xf
}

func ReadHighByteHighNibble(ins Instructions) uint8 {
	return uint8(ins[0]) >> 4
}

func ReadLowByteHighNibble(ins Instructions) uint8 {
	return uint8(ins[1]) >> 4
}

func ReadHighestByte(ins Instructions) uint8 {
	return uint8(ins[0])
}

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		op := ParseOpcode(ins)
		def, err := Lookup(byte(op))

		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		} else {

			operands := ReadOperands(def, ins)

			fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		}
		i += 2
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if operandCount != len(operands) {
		return fmt.Sprintf("ERROR: operand length %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	case 3:
		return fmt.Sprintf("%s %d %d %d", def.Name, operands[0], operands[1], operands[2])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}
