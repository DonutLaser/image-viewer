package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type KeyState uint32

const (
	None KeyState = iota
	JustPressed
	Pressed
	Released
)

type Mouse struct {
	X           int32
	Y           int32
	LeftButton  KeyState
	RightButton KeyState
}

type Keyboard struct {
	LeftArrow  KeyState
	RightArrow KeyState
}

type TimeAmount struct {
	Value uint32 // Milliseconds
	Text  string
}

type App struct {
	Renderer Renderer

	MouseState    Mouse
	KeyboardState Keyboard

	WindowWidth  int32
	WindowHeight int32

	CurrentImageIndex int
	Images            []Image

	ShowTimer        bool
	Times            []TimeAmount
	CurrentTimeIndex int
	StartTime        uint32
}

func Init(renderer *sdl.Renderer) App {
	app := App{}
	app.Renderer = Renderer{
		Instance: renderer,
	}

	font, err := ttf.OpenFont("consola.ttf", 16)
	checkError(err)
	app.Renderer.Font = font
	app.Renderer.Instance.SetDrawColor(0, 0, 0, 255)

	app.CurrentTimeIndex = 0
	app.Times = make([]TimeAmount, 7)
	app.Times[0] = TimeAmount{Value: 0, Text: "Slideshow turned off"}
	app.Times[1] = TimeAmount{Value: 60000, Text: "1 minute"}
	app.Times[2] = TimeAmount{Value: 120000, Text: "2 minutes"}
	app.Times[3] = TimeAmount{Value: 300000, Text: "5 minutes"}
	app.Times[4] = TimeAmount{Value: 600000, Text: "10 minutes"}
	app.Times[5] = TimeAmount{Value: 1800000, Text: "30 minutes"}
	app.Times[6] = TimeAmount{Value: 3600000, Text: "1 hour"}

	return app
}

func (app *App) Tick() {
	if len(app.Images) == 0 {
		app.initMode()
	} else {
		app.imagesMode()
	}
}

func (app *App) initMode() {
	messageWidth, messageHeight := app.Renderer.GetStringSize("Click anywhere to choose a directory...")
	messageRect := sdl.Rect{
		X: app.WindowWidth/2 - messageWidth/2,
		Y: app.WindowHeight/2 - messageHeight/2,
		W: messageWidth,
		H: messageHeight,
	}

	message := "Click anywhere to choose a directory..."

	if app.MouseState.LeftButton == Released {
		directory, err := dialog.Directory().Title("Choose directory...").Browse()
		if err == nil {
			app.loadImagesInDir(directory)
			app.CurrentImageIndex = 0
		}
	}

	app.Renderer.Instance.Clear()
	app.Renderer.DrawText(message, &messageRect, sdl.Color{R: 122, G: 122, B: 122, A: 255})
	app.Renderer.Instance.Present()

}

func (app *App) imagesMode() {
	if app.KeyboardState.RightArrow == JustPressed {
		app.changeImage(1)
	} else if app.KeyboardState.LeftArrow == JustPressed {
		app.changeImage(-1)
	}

	if app.MouseState.RightButton == Released {
		app.ShowTimer = !app.ShowTimer
	} else if app.MouseState.LeftButton == JustPressed && app.ShowTimer {
		app.CurrentTimeIndex = wrap(app.CurrentTimeIndex+1, 0, len(app.Times)-1)

		if app.CurrentImageIndex > 0 {
			app.StartTime = sdl.GetTicks()
		}
	}

	if app.CurrentTimeIndex > 0 {
		currentTime := sdl.GetTicks()

		if currentTime >= app.StartTime+app.Times[app.CurrentTimeIndex].Value {
			app.changeImage(1)
			app.StartTime = currentTime
		}
	}

	app.Renderer.Instance.SetDrawColor(0, 0, 0, 255)
	app.Renderer.Instance.Clear()

	app.Images[app.CurrentImageIndex].Render(app.Renderer.Instance, app.WindowWidth, app.WindowHeight)

	if app.ShowTimer {
		text := app.Times[app.CurrentTimeIndex].Text
		textWidth, textHeight := app.Renderer.GetStringSize(text)
		timerRect := sdl.Rect{
			X: app.WindowWidth/2 - textWidth/2 - 50,
			Y: app.WindowHeight - 50 - textHeight - 100,
			W: textWidth + 100,
			H: textHeight + 100,
		}

		textRect := sdl.Rect{
			X: timerRect.X + 50,
			Y: timerRect.Y + 50,
			W: textWidth,
			H: textHeight,
		}

		app.Renderer.DrawRect(&timerRect, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		app.Renderer.DrawRectOutline(&timerRect, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		app.Renderer.DrawText(text, &textRect, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	}

	app.Renderer.Instance.Present()
}

func (app *App) Close() {
	app.Renderer.Font.Close()
}

// @TODO (!important) prevent lag when loading, goroutines?
func (app *App) loadImagesInDir(dir string) {
	supportedTypes := []string{".png", ".jpg", ".bmp"} // Must have dots
	images := findImagesInDir(dir, supportedTypes)

	for _, image := range images {
		app.Images = append(app.Images, loadImage(image, app.Renderer.Instance))
	}
}

func (app *App) changeImage(direction int) {
	app.CurrentImageIndex = wrap(app.CurrentImageIndex+direction, 0, len(app.Images)-1)
}

func findImagesInDir(dir string, supportedTypes []string) (result []string) {
	files, err := ioutil.ReadDir(dir)
	checkError(err)

	for _, file := range files {
		extension := strings.ToLower(path.Ext(file.Name()))
		if containsString(supportedTypes, extension) {
			result = append(result, fmt.Sprintf("%s/%s", dir, file.Name()))
		}
	}

	return
}
