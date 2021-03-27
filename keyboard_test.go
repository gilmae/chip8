package chip8

import (
	"strings"
	"testing"
)

func TestReadKey(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput byte
		expectedOk     bool
	}{
		{
			"1",
			1,
			true,
		},
		{
			"2",
			2,
			true,
		},
		{
			"3",
			3,
			true,
		},
		{
			"4",
			0xc,
			true,
		},
		{
			"q",
			4,
			true,
		},
		{
			"w",
			5,
			true,
		},
		{
			"e",
			6,
			true,
		},
		{
			"r",
			0xd,
			true,
		},
		{
			"a",
			7,
			true,
		},
		{
			"s",
			8,
			true,
		},
		{
			"d",
			9,
			true,
		},
		{
			"f",
			0xe,
			true,
		},
		{
			"z",
			0xa,
			true,
		},
		{
			"x",
			0,
			true,
		},
		{
			"c",
			0xb,
			true,
		},
		{
			"v",
			0xf,
			true,
		},
		{
			"p1",
			1,
			true,
		},
		{
			"p",
			0,
			false,
		},
		{
			"",
			0,
			false,
		},
	}

	for _, tt := range tests {
		input := strings.NewReader(tt.input)
		k := NewKeyboard(input)

		actualOutput, ok := k.readKey()

		if ok != tt.expectedOk {
			t.Errorf("unexpected read result, want=%t, got=%t", tt.expectedOk, ok)
		}

		if actualOutput != tt.expectedOutput {
			t.Errorf("unexpected output, want=%d, got=%d", tt.expectedOutput, actualOutput)
		}

	}
}
