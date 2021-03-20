package chip8

import "fmt"

type cpu struct {
	memory    [4096]byte
	registers [16]byte

	I uint16 // Instruction Register

	delay byte // delay timer register
	sound byte //sound timer register

	stack [16]uint16

	pc uint16 // Program counter
	sp uint8  // stack pointer
}

func NewCpu() *cpu {
	return &cpu{
		pc: 0x200, // First 512 bytes are "reserved" for the Chip-8 "interpreter"
	}
}

func (c *cpu) Tick() error {
	ins := c.memory[c.pc : c.pc+2]
	c.pc += 2
	op := ParseOpcode(ins)
	// def, err := Lookup(byte(op))
	// if err != nil {
	// 	return
	// }
	// operands := ReadOperands(def, ins)
	switch op {
	case SYS:
		// nothing, ignore
	case CLS:
		// clearDisplay
	case RET:
		val, err := c.pop()
		if err != nil {
			return err
		}
		c.pc = val
	case JP:
		addr := ReadUint12(ins)
		c.pc = addr
	case CALL:
		addr := ReadUint12(ins)
		err := c.push(c.pc)
		if err != nil {
			return err
		}
		c.pc = addr
	case SE:
		val := ReadUint8(ins)
		register := ReadHighByteNibble(ins)
		if c.registers[register] == val {
			c.pc += 2
		}
	case SNE:
		val := ReadUint8(ins)
		register := ReadHighByteNibble(ins)
		if c.registers[register] != val {
			c.pc += 2
		}
	case SRE:
		xregister := ReadHighByteNibble(ins)
		yregister := ReadLowByteHighNibble(ins)

		if c.registers[xregister] == c.registers[yregister] {
			c.pc += 2
		}
	case SRNE:
		xregister := ReadHighByteNibble(ins)
		yregister := ReadLowByteHighNibble(ins)

		if c.registers[xregister] != c.registers[yregister] {
			c.pc += 2
		}
	case LD:
		register := ReadHighByteNibble(ins)
		val := ReadUint8(ins)
		c.registers[register] = val
	case LDI:
		addr := ReadUint12(ins)
		c.I = addr
	case LDVxDT:
		register := ReadHighByteNibble(ins)
		c.registers[register] = c.delay
	case LDDTVx:
		register := ReadHighByteNibble(ins)
		c.delay = c.registers[register]
	case LDSTVx:
		register := ReadHighByteNibble(ins)
		c.sound = c.registers[register]
	case LDB:
		register := ReadHighByteNibble(ins)
		value := int(c.registers[register])
		c.memory[c.I] = byte(value / 100)
		c.memory[c.I+1] = byte((value % 100) / 10)
		c.memory[c.I+2] = byte(value % 10)
	case LDIVx:
		register := ReadHighByteNibble(ins)
		for idx := 0; idx <= int(register); idx++ {
			c.memory[int(c.I)+idx] = c.registers[idx]
		}
	case LDVxI:
		register := ReadHighByteNibble(ins)
		for idx := 0; idx <= int(register); idx++ {
			c.registers[idx] = c.memory[int(c.I)+idx]
		}
	}
	return nil
}

func (c *cpu) pop() (uint16, error) {
	if c.sp == 0 {
		return 0, fmt.Errorf("cannot pop from stack empty")
	}
	v := c.stack[c.sp-1]
	c.sp--
	return v, nil
}

func (c *cpu) push(value uint16) error {
	if int(c.sp) >= len(c.stack) {
		return fmt.Errorf("stack overflow")
	}
	c.stack[c.sp] = value
	c.sp++
	return nil
}
