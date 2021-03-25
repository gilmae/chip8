package chip8

const (
	width  int = 64
	height int = 32
)

type display struct {
	pixels [][]bool
}

func NewDisplay() display {
	d := display{}
	d.Clear()
	return d
}

func (d *display) Clear() {
	d.pixels = make([][]bool, height)
	for idx := range d.pixels {
		d.pixels[idx] = make([]bool, width)
	}
}

func (d *display) DrawPixel(pixel []byte, x int, y int) bool {
	collision_detected := false

	for row_offset, sprite_row := range pixel {
		bit := 7
		for sprite_row > 0 {
			if sprite_row%2 == 1 {
				collision_detected = collision_detected || d.pixels[y+row_offset][(x+bit)%width]
				d.pixels[y+row_offset][(x+bit)%width] = true
			}
			sprite_row = sprite_row >> 1
			bit--
		}

	}
	return collision_detected
}
