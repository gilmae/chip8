package chip8

import "testing"

func TestConstruction(t *testing.T) {
	d := NewDisplay()

	if len(d.pixels) != 32 {
		t.Errorf("screen not high enough, want=%d, got=%d", 32, len(d.pixels))
	}

	if len(d.pixels[0]) != 64 {
		t.Errorf("screen not wide enough, want=%d, got=%d", 64, len(d.pixels))
	}
}

func TestDrawPixel(t *testing.T) {
	tests := []struct {
		x              int
		y              int
		input          []byte
		expectedPixels [][]int
	}{
		{
			0,
			0,
			[]byte{0xF0, 0x90, 0x90, 0x90, 0xF0},
			[][]int{{0, 1, 2, 3}, {0, 3}, {0, 3}, {0, 3}, {0, 1, 2, 3}},
		},
		{
			62,
			0,
			[]byte{0xF0, 0x90, 0x90, 0x90, 0xF0},
			[][]int{{62, 63, 0, 1}, {62, 1}, {62, 1}, {62, 1}, {62, 63, 0, 1}},
		},
	}

	for _, tt := range tests {
		d := NewDisplay()
		d.DrawPixel(tt.input, tt.x, tt.y)

		total_expected_pixel_count := 0
		total_actual_pixel_count := 0

		for y, row := range tt.expectedPixels {
			for _, x := range row {
				total_expected_pixel_count++
				if !d.pixels[tt.y+y][x] {
					t.Errorf("expected pixel %d,%d to be on", x, tt.y+y)
				}
			}
		}

		for y, row := range d.pixels {
			for x, _ := range row {
				if d.pixels[y][x] {
					total_actual_pixel_count++
				}
			}
		}

		if total_actual_pixel_count != total_expected_pixel_count {
			t.Errorf("expected %d pixels on, got %d", total_expected_pixel_count, total_actual_pixel_count)
		}
	}

}

func TestCollision(t *testing.T) {
	input := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0}
	d := NewDisplay()
	d.DrawPixel(input, 0, 0)
	collision := d.DrawPixel(input, 10, 10)

	if collision {
		t.Errorf("expected no collision")
	}

	collision = d.DrawPixel(input, 1, 1)

	if !collision {
		t.Errorf("expected collision")
	}

}

func TestClear(t *testing.T) {
	input := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0}
	d := NewDisplay()
	d.DrawPixel(input, 0, 0)
	d.Clear()

	total_actual_pixel_count := 0
	for y, row := range d.pixels {
		for x, _ := range row {
			if d.pixels[y][x] {
				total_actual_pixel_count++
			}
		}
	}

	if total_actual_pixel_count != 0 {
		t.Errorf("expected %d pixels on, got %d", 0, total_actual_pixel_count)
	}
}
