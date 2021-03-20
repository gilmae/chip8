package chip8

import (
	"testing"
)

func TestParseOpcode(t *testing.T) {
	tests := []struct {
		input          Instructions
		expectedOpcode Opcode
	}{
		{[]byte{0x01, 0x23}, SYS},
		{[]byte{0x00, 0xe0}, CLS},
		{[]byte{0x00, 0xee}, RET},
	}

	for _, tt := range tests {
		actualOpcode := ParseOpcode(tt.input)

		if actualOpcode != tt.expectedOpcode {
			t.Errorf("wrong opcode, want=%d, got=%d", tt.expectedOpcode, actualOpcode)
		}
	}
}

func TestInstructionString(t *testing.T) {
	tests := []struct {
		input          Instructions
		expectedString string
	}{
		{[]byte{0x01, 0x23}, "0000 SYS 291\n"},
		{[]byte{0x00, 0xe0}, "0000 CLS\n"},
		{[]byte{0x00, 0xee}, "0000 RET\n"},
	}

	for _, tt := range tests {
		insString := tt.input.String()

		if insString != tt.expectedString {
			t.Errorf("instruction wrongly formatted, want=%s, got=%s", tt.expectedString, insString)
		}
	}
}

func TestReadNibble(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 15},
		{[]byte{0x00, 0xff}, 15},
		{[]byte{0xff, 0x01}, 1},
		{[]byte{0xff, 0x0f}, 15},
		{[]byte{0xff, 0x0a}, 10},
	}

	for _, tt := range tests {
		actualValue := ReadNibble(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestReadUint16(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint16
	}{
		{[]byte{0xff, 0xff}, 65535},
		{[]byte{0xff, 0x00}, 65280},
		{[]byte{0x00, 0xff}, 255},
		{[]byte{0xff, 0x01}, 65281},
		{[]byte{0xff, 0x0f}, 65295},
		{[]byte{0xff, 0x0a}, 65290},
	}

	for _, tt := range tests {
		actualValue := ReadUint16(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestReadUint12(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint16
	}{
		{[]byte{0xff, 0xff}, 4095},
		{[]byte{0xff, 0x00}, 3840},
		{[]byte{0x00, 0xff}, 255},
		{[]byte{0xff, 0x01}, 3841},
		{[]byte{0xff, 0x0f}, 3855},
		{[]byte{0xff, 0x0a}, 3850},
	}

	for _, tt := range tests {
		actualValue := ReadUint12(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestReadUint8(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 255},
		{[]byte{0xff, 0x00}, 0},
		{[]byte{0x00, 0xff}, 255},
		{[]byte{0xff, 0x01}, 1},
		{[]byte{0xff, 0x0f}, 15},
		{[]byte{0xff, 0x0a}, 10},
	}

	for _, tt := range tests {
		actualValue := ReadUint8(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestHighByteNibble(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 15},
		{[]byte{0x00, 0xff}, 0},
		{[]byte{0xff, 0x01}, 15},
		{[]byte{0xfa, 0x0f}, 10},
		{[]byte{0x01, 0x0a}, 1},
	}

	for _, tt := range tests {
		actualValue := ReadHighByteNibble(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestLowByteHighNibble(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 15},
		{[]byte{0x00, 0xff}, 15},
		{[]byte{0xff, 0x01}, 0},
		{[]byte{0xfa, 0x0f}, 0},
		{[]byte{0x01, 0xa0}, 10},
	}

	for _, tt := range tests {
		actualValue := ReadLowByteHighNibble(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestHighByteHighNibble(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 15},
		{[]byte{0x00, 0xff}, 0},
		{[]byte{0xff, 0x01}, 15},
		{[]byte{0xfa, 0x0f}, 15},
		{[]byte{0x01, 0xa0}, 0},
		{[]byte{0xaf, 0xa0}, 10},
	}

	for _, tt := range tests {
		actualValue := ReadHighByteHighNibble(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}

func TestHighestByte(t *testing.T) {
	tests := []struct {
		input          []byte
		expectedResult uint8
	}{
		{[]byte{0xff, 0xff}, 255},
		{[]byte{0x00, 0xff}, 0},
		{[]byte{0xff, 0x01}, 255},
		{[]byte{0xfa, 0x0f}, 250},
		{[]byte{0x01, 0xa0}, 1},
		{[]byte{0xaf, 0xa0}, 175},
	}

	for _, tt := range tests {
		actualValue := ReadHighestByte(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}
