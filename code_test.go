package chip8

import (
	"testing"
)

func TestReadNibble(t *testing.T) {
	tests := []struct {
		input          [2]byte
		expectedResult uint8
	}{
		{[2]byte{0xff, 0xff}, 15},
		{[2]byte{0x00, 0xff}, 15},
		{[2]byte{0xff, 0x01}, 1},
		{[2]byte{0xff, 0x0f}, 15},
		{[2]byte{0xff, 0x0a}, 10},
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
		input          [2]byte
		expectedResult uint16
	}{
		{[2]byte{0xff, 0xff}, 65535},
		{[2]byte{0xff, 0x00}, 65280},
		{[2]byte{0x00, 0xff}, 255},
		{[2]byte{0xff, 0x01}, 65281},
		{[2]byte{0xff, 0x0f}, 65295},
		{[2]byte{0xff, 0x0a}, 65290},
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
		input          [2]byte
		expectedResult uint16
	}{
		{[2]byte{0xff, 0xff}, 4095},
		{[2]byte{0xff, 0x00}, 3840},
		{[2]byte{0x00, 0xff}, 255},
		{[2]byte{0xff, 0x01}, 3841},
		{[2]byte{0xff, 0x0f}, 3855},
		{[2]byte{0xff, 0x0a}, 3850},
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
		input          [2]byte
		expectedResult uint8
	}{
		{[2]byte{0xff, 0xff}, 255},
		{[2]byte{0xff, 0x00}, 0},
		{[2]byte{0x00, 0xff}, 255},
		{[2]byte{0xff, 0x01}, 1},
		{[2]byte{0xff, 0x0f}, 15},
		{[2]byte{0xff, 0x0a}, 10},
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
		input          [2]byte
		expectedResult uint8
	}{
		{[2]byte{0xff, 0xff}, 15},
		{[2]byte{0x00, 0xff}, 0},
		{[2]byte{0xff, 0x01}, 15},
		{[2]byte{0xfa, 0x0f}, 10},
		{[2]byte{0x01, 0x0a}, 1},
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
		input          [2]byte
		expectedResult uint8
	}{
		{[2]byte{0xff, 0xff}, 15},
		{[2]byte{0x00, 0xff}, 15},
		{[2]byte{0xff, 0x01}, 0},
		{[2]byte{0xfa, 0x0f}, 0},
		{[2]byte{0x01, 0xa0}, 10},
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
		input          [2]byte
		expectedResult uint8
	}{
		{[2]byte{0xff, 0xff}, 15},
		{[2]byte{0x00, 0xff}, 0},
		{[2]byte{0xff, 0x01}, 15},
		{[2]byte{0xfa, 0x0f}, 15},
		{[2]byte{0x01, 0xa0}, 0},
		{[2]byte{0xaf, 0xa0}, 10},
	}

	for _, tt := range tests {
		actualValue := ReadHighByteHighNibble(tt.input)

		if actualValue != tt.expectedResult {
			t.Errorf("Wrong value, expected=%d, got=%d", tt.expectedResult, actualValue)
		}
	}
}
