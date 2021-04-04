package chip8

import "github.com/veandco/go-sdl2/sdl"

type renderer interface {
	Close()
	Render(d *display) error
}

type nullRenderer struct{}

func NewNullRenderer() *nullRenderer {
	return &nullRenderer{}
}

func (n *nullRenderer) Close() {

}

func (n *nullRenderer) Render(d *display) error {
	return nil
}

// type termboxRenderer struct {
// 	fg, bg termbox.Attribute
// }

// func NewTermboxRenderer(fg, bg termbox.Attribute) (*termboxRenderer, error) {
// 	t := &termboxRenderer{fg: fg, bg: bg}

// 	return t, t.init()
// }

// func (t *termboxRenderer) init() error {
// 	if err := termbox.Init(); err != nil {
// 		return err
// 	}

// 	termbox.HideCursor()

// 	if err := termbox.Clear(t.bg, t.bg); err != nil {
// 		return err
// 	}

// 	return termbox.Flush()
// }

// func (t *termboxRenderer) Close() {
// 	termbox.Close()
// }

// func (t *termboxRenderer) Render(d *display) error {
// 	d.EachPixel(func(x, y uint16, addr int) {
// 		v := ' '

// 		if d.pixels[addr] {
// 			v = '█'
// 		}

// 		termbox.SetCell(
// 			int(x),
// 			int(y),
// 			v,
// 			t.fg,
// 			t.bg,
// 		)
// 	})

// 	return termbox.Flush()
// }

type sdlRenderer struct {
	window  *sdl.Window
	surface *sdl.Surface
	scale   int32
}

func NewSdlRenderer(scale int32, w *sdl.Window, surface *sdl.Surface) *sdlRenderer {
	return &sdlRenderer{scale: scale, window: w, surface: surface}
}

func (s *sdlRenderer) Close() {

}

func (s *sdlRenderer) Render(d *display) error {
	s.surface.FillRect(nil, 0)
	d.EachPixel(func(x, y uint16, addr int) {
		if d.pixels[addr] {
			rect := sdl.Rect{X: int32(x) * s.scale, Y: int32(y) * s.scale, W: s.scale, H: s.scale}
			s.surface.FillRect(&rect, 0xffffffff)
		}

	})
	s.window.UpdateSurface()

	return nil
}
