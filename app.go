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

type App struct {
	Renderer          Renderer
	MouseState        Mouse
	KeyboardState     Keyboard
	WindowWidth       int32
	WindowHeight      int32
	CurrentImageIndex int32
	CurrentImage      Image
	Images            []Image
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
			app.CurrentImage = app.Images[app.CurrentImageIndex]
		}
	}

	app.Renderer.Instance.Clear()
	app.Renderer.DrawText(message, &messageRect, sdl.Color{R: 122, G: 122, B: 122, A: 255})
	app.Renderer.Instance.Present()

}

func (app *App) imagesMode() {
	if app.KeyboardState.RightArrow == JustPressed {
		nextImage := clamp(app.CurrentImageIndex+1, 0, int32(len(app.Images)-1))
		if nextImage != app.CurrentImageIndex {
			app.CurrentImageIndex = nextImage
			app.CurrentImage = app.Images[app.CurrentImageIndex]
		}
	} else if app.KeyboardState.LeftArrow == JustPressed {
		nextImage := clamp(app.CurrentImageIndex-1, 0, int32(len(app.Images)-1))
		if nextImage != app.CurrentImageIndex {
			app.CurrentImageIndex = nextImage
			app.CurrentImage = app.Images[app.CurrentImageIndex]
		}
	}

	app.Renderer.Instance.Clear()
	app.CurrentImage.Render(app.Renderer.Instance, app.WindowWidth, app.WindowHeight)
	app.Renderer.Instance.Present()
}

func (app *App) Close() {
	app.Renderer.Font.Close()
}

func (app *App) loadImagesInDir(dir string) {
	supportedTypes := []string{".png", ".jpg", ".bmp"} // Must have dots
	images := findImagesInDir(dir, supportedTypes)

	for _, image := range images {
		app.Images = append(app.Images, loadImage(image, app.Renderer.Instance))
	}
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
