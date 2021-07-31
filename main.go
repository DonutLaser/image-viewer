package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func loadImagesInDir(dir string, renderer *sdl.Renderer) (result []Image) {
	supportedTypes := []string{".png", ".jpg", ".bmp"} // Must have dots
	images := findImagesInDir(dir, supportedTypes)

	for _, image := range images {
		result = append(result, loadImage(image, renderer))
	}

	return
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

func render(window *sdl.Window, renderer *sdl.Renderer, image Image) {
	screenWidth, screenHeight := window.GetSize()

	renderer.Clear()
	image.Render(renderer, screenWidth, screenHeight)
	renderer.Present()
}

// @TODO (!important) slideshow

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	checkError(err)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Image Viewer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_RESIZABLE)
	checkError(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkError(err)
	defer renderer.Destroy()

	var currentImageIndex int32 = 0
	// @TODO (!important) somehow make it not lag on load
	images := loadImagesInDir("D:/Wallpapers", renderer)
	currentImage := images[currentImageIndex]

	renderer.SetDrawColor(0, 0, 0, 255)

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym

				// @TODO (!important) `key up and key down to zoom in/out`
				switch keyCode {
				case sdl.K_RIGHT:
					if t.Repeat > 0 || t.State == sdl.RELEASED {
						break
					}

					nextImage := clamp(currentImageIndex+1, 0, int32(len(images)-1))
					if nextImage != currentImageIndex {
						currentImageIndex = nextImage
						currentImage = images[currentImageIndex]
					}
				case sdl.K_LEFT:
					if t.Repeat > 0 || t.State == sdl.RELEASED {
						break
					}

					nextImage := clamp(currentImageIndex-1, 0, int32(len(images)-1))
					if nextImage != currentImageIndex {
						currentImageIndex = nextImage
						currentImage = images[currentImageIndex]
					}
				}
			}
		}

		render(window, renderer, currentImage)
	}
}
