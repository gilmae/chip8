package chip8

import "fmt"

const (
	width  int = 64
	height int = 32
)

type display struct {
	pixels  []bool
	isDirty bool
}

func NewDisplay() display {
	d := display{}
	d.Clear()
	return d
}

func (d *display) Clear() {
	d.pixels = make([]bool, width*height)
	d.isDirty = true
}

func (d *display) DrawSprite(pixel []byte, x int, y int) bool {
	collision_detected := false

	for row_offset, sprite_row := range pixel {
		bit := 7
		for sprite_row > 0 {
			if sprite_row%2 == 1 {
				px := d.addrOf(((x + bit) % width), y+row_offset)

				collision_detected = collision_detected || d.pixels[px]

				d.pixels[px] = !d.pixels[px]
			}
			sprite_row = sprite_row >> 1
			bit--
		}

	}
	d.isDirty = true
	return collision_detected
}

func (d *display) GetPixel(x int, y int) (bool, error) {
	px := d.addrOf(x, y)

	if px > height*width {
		return false, fmt.Errorf("out of bounds")
	}
	return d.pixels[px], nil
}

func (d *display) addrOf(x int, y int) int {
	return y*width + x
}
