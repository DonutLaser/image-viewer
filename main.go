package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	checkError(err)
	defer sdl.Quit()

	err = ttf.Init()
	checkError(err)
	defer ttf.Quit()

	window, err := sdl.CreateWindow("Image Viewer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_RESIZABLE)
	checkError(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	checkError(err)
	defer renderer.Destroy()

	app := Init(renderer)

	running := true
	for running {
		lastLeftState := app.MouseState.LeftButton
		lastRightState := app.MouseState.RightButton

		app.MouseState.LeftButton = None
		app.MouseState.RightButton = None
		app.KeyboardState.LeftArrow = None
		app.KeyboardState.RightArrow = None

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseMotionEvent:
				app.MouseState.X = t.X
				app.MouseState.Y = t.Y
			case *sdl.MouseButtonEvent:
				if t.Button == sdl.BUTTON_LEFT {
					if t.State == sdl.RELEASED {
						app.MouseState.LeftButton = Released
					} else if t.State == sdl.PRESSED && lastLeftState == JustPressed {
						app.MouseState.LeftButton = Pressed
					} else {
						app.MouseState.LeftButton = JustPressed
					}
				} else if t.Button == sdl.BUTTON_RIGHT {
					if t.State == sdl.RELEASED {
						app.MouseState.RightButton = Released
					} else if t.State == sdl.PRESSED && lastRightState == JustPressed {
						app.MouseState.RightButton = Pressed
					} else {
						app.MouseState.RightButton = JustPressed
					}
				}
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym

				switch keyCode {
				case sdl.K_RIGHT:
					if t.State == sdl.RELEASED {
						app.KeyboardState.RightArrow = Released
					} else if t.Repeat > 0 {
						app.KeyboardState.RightArrow = Pressed
					} else {
						app.KeyboardState.RightArrow = JustPressed
					}
				case sdl.K_LEFT:
					if t.State == sdl.RELEASED {
						app.KeyboardState.LeftArrow = Released
					} else if t.Repeat > 0 {
						app.KeyboardState.LeftArrow = Pressed
					} else {
						app.KeyboardState.LeftArrow = JustPressed
					}
				}
			}
		}

		app.WindowWidth, app.WindowHeight = window.GetSize()
		app.Tick()
	}

	app.Close()
}
