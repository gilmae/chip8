package chip8

type display struct {
	pixels [][]bool
}

func NewDisplay() display {
	d := display{}
	d.Clear()
	return d
}

func (d *display) Clear() {
	d.pixels = make([][]bool, 32)
	for idx := range d.pixels {
		d.pixels[idx] = make([]bool, 64)
	}
}

func (d *display) DrawPixel(pixel []byte, x int, y int) bool {
	collision_detected := false

	for row_offset, sprite_row := range pixel {
		bit := 7
		for sprite_row > 0 {
			if sprite_row%2 == 1 {
				collision_detected = collision_detected || d.pixels[y+row_offset][x+bit]
				d.pixels[y+row_offset][x+bit] = true
			}
			sprite_row = sprite_row >> 1
			bit--
		}

	}
	return collision_detected
}
