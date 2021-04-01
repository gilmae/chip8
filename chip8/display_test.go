package chip8

import "testing"

func TestConstruction(t *testing.T) {
	d := NewDisplay()

	if len(d.pixels) != 64*32 {
		t.Errorf("screen does not have enough pixels, want=%d, got=%d", 32*64, len(d.pixels))
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
		d.DrawSprite(tt.input, tt.x, tt.y)

		total_expected_pixel_count := 0
		total_actual_pixel_count := 0

		for y, row := range tt.expectedPixels {
			for _, x := range row {
				total_expected_pixel_count++
				px, err := d.GetPixel(x, tt.y+y)
				if err != nil {
					t.Errorf("pixel at %d,%d could not be resolved, got %s", x, tt.y+y, err)
				}

				if !px {
					t.Errorf("expected pixel %d,%d to be on", x, tt.y+y)
				}
			}
		}

		for _, px := range d.pixels {
			if px {
				total_actual_pixel_count++
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
	d.DrawSprite(input, 0, 0)
	collision := d.DrawSprite(input, 10, 10)

	if collision {
		t.Errorf("expected no collision")
	}

	collision = d.DrawSprite(input, 3, 3)

	if !collision {
		t.Errorf("expected collision")
	}

}

func TestPixelsXored(t *testing.T) {
	input := []byte{0x80}
	d := NewDisplay()
	d.DrawSprite(input, 0, 0)
	d.DrawSprite(input, 0, 0)

	total_actual_pixel_count := 0
	for _, px := range d.pixels {
		if px {
			total_actual_pixel_count++
		}
	}

	if total_actual_pixel_count != 0 {
		t.Errorf("expected %d pixels on, got %d", 0, total_actual_pixel_count)
	}
}

func TestClear(t *testing.T) {
	input := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0}
	d := NewDisplay()
	d.DrawSprite(input, 0, 0)
	d.Clear()

	total_actual_pixel_count := 0
	for _, px := range d.pixels {
		if px {
			total_actual_pixel_count++
		}
	}

	if total_actual_pixel_count != 0 {
		t.Errorf("expected %d pixels on, got %d", 0, total_actual_pixel_count)
	}
}
