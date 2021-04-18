package chip8

import "fmt"

const (
	width  int = 64
	height int = 32
)

type display struct {
	pixels        []bool
	isDirty       bool
	width, height int
}

func NewDisplay() display {
	d := display{width: width, height: height}
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
				px := d.addrOf(d.normalisePixelCoords(x+bit, y+row_offset))

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

func (d *display) EachPixel(fn func(x, y uint16, addr int)) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			a := d.addrOf(x, y)
			fn(uint16(x), uint16(y), a)
		}
	}
}

func (d *display) GetPixel(x int, y int) (bool, error) {
	px := d.addrOf(d.normalisePixelCoords(x, y))

	if px > height*width {
		return false, fmt.Errorf("out of bounds")
	}
	return d.pixels[px], nil
}

func (d *display) addrOf(x int, y int) int {
	return y*width + x
}

func (d *display) normalisePixelCoords(x, y int) (dx, dy int) {
	dx = x % width
	if dx < 0 {
		dx += width
	}

	dy = y % height
	if dy < 0 {
		dy += height
	}

	return
}
