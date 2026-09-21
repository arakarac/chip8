package main

import (
	//"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	w =        64 
	h =        32
	mod =      10
)

var (
	window *sdl.Window
	renderer *sdl.Renderer
)

func sdlInit() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic("SDL init fail")
	}
	//defer sdl.Quit()

  var err error
	window, err = sdl.CreateWindow("Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w*mod, h*mod, sdl.WINDOW_SHOWN)
  if err != nil {
		panic(err)
	}
	//defer window.Destroy() 

	renderer, err = sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
	//defer renderer.Destroy()
}

func (c *chip8) render() {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
 
	b := c.buffer()

	renderer.SetDrawColor(255,255,255,255)
	for j := 0; j < len(b); j++ {
		for i := 0; i <len(b[j]); i++ {
			if b[j][i] != 0 {
				renderer.FillRect(&sdl.Rect{
					int32(i) * mod,
					int32(j) * mod,
					mod,
					mod,
				})
			}
		}
	}

	renderer.Present()
}

func sdlEnd() {
	window.Destroy()
	renderer.Destroy()
	sdl.Quit()
}

