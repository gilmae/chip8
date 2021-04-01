package chip8

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	font_start_addr    uint16 = 0x50
	program_start_addr uint16 = 0x200
)

var (
	DefaultLogger = log.New(os.Stdout, "", 0)
)

type cpu struct {
	memory    [4096]byte
	registers [16]byte

	I uint16 // Instruction Register

	delay byte // delay timer register
	sound byte //sound timer register

	stack [16]uint16

	pc uint16 // Program counter
	sp uint8  // stack pointer

	d        display
	keyboard *keyboard
	logger   *log.Logger
	clock    <-chan time.Time
	stop     chan struct{}
	r        Renderer
}

func NewCpu(k *keyboard, r Renderer) *cpu {
	c := &cpu{
		pc:       program_start_addr, // First 512 bytes are "reserved" for the Chip-8 "interpreter"
		d:        NewDisplay(),
		keyboard: k,
		logger:   DefaultLogger,
		clock:    time.Tick(time.Duration(60)),
		stop:     make(chan struct{}),
		r:        r,
	}
	c.loadFont()

	return c
}

func (c *cpu) LoadBytes(program []byte) (int, error) {
	reader := bytes.NewReader(program)
	return c.Load(reader)
}

func (c *cpu) Load(reader io.Reader) (int, error) {
	return c.load(reader, program_start_addr)
}

func (c *cpu) Run() error {
	for {
		select {
		case <-c.stop:
			return nil
		case <-c.clock:
			err := c.Tick()
			if err != nil {
				return err
			}
		}
	}
}

func (c *cpu) Stop() {
	close(c.stop)
}

func (c *cpu) Tick() error {
	if c.delay > 0 {
		c.delay--
	}
	if c.sound > 0 {
		c.sound--
	}

	ins := c.memory[c.pc : c.pc+2]
	c.logger.Println(Instructions(ins).String())
	c.pc += 2
	op := ParseOpcode(ins)

	switch op {
	case SYS:
		// nothing, ignore
	case CLS:
		c.d.Clear()
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
	case LDVxVy:
		xregister := ReadHighByteNibble(ins)
		yregister := ReadLowByteHighNibble(ins)

		c.registers[xregister] = c.registers[yregister]
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
	case ADD:
		register := ReadHighByteNibble(ins)
		val := ReadUint8(ins)
		c.registers[register] += val
	case ADDVxVy:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		value := c.registers[registerx] + c.registers[registery]
		overflow := 0
		if value > 255 {
			overflow = 1
			value = value & 0xff
		}

		c.registers[0xf] = byte(overflow)
		c.registers[registerx] = value
	case ADDIVx:
		register := ReadHighByteNibble(ins)
		c.I += uint16(c.registers[register])
	case OR:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		c.registers[registerx] |= c.registers[registery]
	case AND:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		c.registers[registerx] &= c.registers[registery]
	case XOR:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		c.registers[registerx] ^= c.registers[registery]
	case SUB:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		vx := c.registers[registerx]
		vy := c.registers[registery]
		if vx > vy {
			c.registers[0xf] = 1
		}
		c.registers[registerx] = byte(vx - vy)
	case SUBN:
		registerx := ReadHighByteNibble(ins)
		registery := ReadLowByteHighNibble(ins)
		vx := c.registers[registerx]
		vy := c.registers[registery]
		if vy > vx {
			c.registers[0xf] = 1
		}
		c.registers[registerx] = byte(vy - vx)
	case SHR:
		register := ReadHighByteNibble(ins)
		c.registers[0xf] = c.registers[register] & 0x1
		c.registers[register] = c.registers[register] >> 1
	case SHL:
		register := ReadHighByteNibble(ins)
		if c.registers[register] >= 0x8 {
			c.registers[0xf] = 1
		} else {
			c.registers[0xf] = 0
		}

		c.registers[register] = byte(c.registers[register] << 1)
	case JPV0:
		addr := ReadUint12(ins)
		c.pc = uint16(c.registers[0]) + addr
	case RND:
		register := ReadHighByteNibble(ins)
		val := ReadUint8(ins)
		random := uint8(rand.Intn(255))
		c.registers[register] = random & val
	case DRW:
		x := int(ReadHighByteNibble(ins))
		y := int(ReadLowByteHighNibble(ins))
		sprite_size := ReadNibble(ins)
		sprite := make([]byte, sprite_size)

		for idx := 0; idx < int(sprite_size); idx++ {
			sprite[idx] = c.memory[int(c.I)+idx]
		}

		collision := c.d.DrawSprite(sprite, x, y)
		if collision {
			c.registers[0xf] = 1
		}
	case LDF:
		register := ReadHighByteNibble(ins)
		c.I = uint16(font_start_addr + uint16(fontwidth)*uint16(c.registers[register]))
	case LDK:
		register := ReadHighByteNibble(ins)
		key, ok := c.keyboard.readKey()
		if !ok {
			c.pc -= 2
		} else {
			c.registers[register] = key
		}
	case SKP:
		register := ReadHighByteNibble(ins)
		key, ok := c.keyboard.readKey()
		if ok && key == c.registers[register] {
			c.pc += 2
		}
	case SKNP:
		register := ReadHighByteNibble(ins)
		key, ok := c.keyboard.readKey()
		if !ok || key != c.registers[register] {
			c.pc += 2
		}
	}

	if c.d.isDirty {
		c.drawScreen()
	}

	if c.sound > 0 {
		c.buzz()
	}

	return nil
}

func (c *cpu) buzz() {

}

func (c *cpu) drawScreen() {
}

func (c *cpu) loadFont() {
	reader := bytes.NewReader(fontset)
	n, err := c.load(reader, font_start_addr)
	if err != nil {
		panic(err)
	}
	if n != len(fontset) {
		panic("fontset loaded incorrectly")
	}
}

func (c *cpu) load(reader io.Reader, offset uint16) (int, error) {
	return reader.Read(c.memory[offset:])
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
