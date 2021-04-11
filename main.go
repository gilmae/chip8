package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gilmae/chip8/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	winHeight int32 = 320
	winWidth  int32 = 640
)

func main() {
	//renderer := chip8.NewNullRenderer()
	//renderer, err := chip8.NewTermboxRenderer(termbox.ColorWhite, termbox.ColorBlack)

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	window, err := sdl.CreateWindow("Chip-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	defer r.Destroy()

	tex, err := r.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(winWidth), int32(winHeight))
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()
	// surface, err := window.GetSurface()
	// if err != nil {
	// 	panic(err)
	// }
	// surface.FillRect(nil, 0)

	renderer := chip8.NewSdlRenderer(int32(winWidth), winHeight, window, r, tex)

	defer renderer.Close()

	keyboard := chip8.NewKeyboard()
	cpu := chip8.NewCpu(keyboard, renderer)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}

	filepath := os.Args[1]
	program, err := ioutil.ReadFile(filepath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(3)
	}

	_, err = cpu.LoadBytes(program)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(4)
	}

	cpu.Run()

}
