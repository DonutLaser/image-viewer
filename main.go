package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

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

func loadImage(path string, renderer *sdl.Renderer) *sdl.Texture {
	image, err := img.Load(path)
	checkError(err)

	texture, err := renderer.CreateTextureFromSurface(image)
	checkError(err)

	image.Free()

	return texture
}

func render(renderer *sdl.Renderer, image *sdl.Texture) {
	renderer.Clear()
	renderer.Copy(image, nil, nil)
	renderer.Present()
}

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

	supportedImages := []string{".png", ".jpg", ".bmp"} // Must have dots
	currentImage := 0
	images := findImagesInDir("D:/Wallpapers", supportedImages)

	renderer.SetDrawColor(0, 0, 0, 255)

	image := loadImage(images[currentImage], renderer)
	defer image.Destroy()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym

				switch keyCode {
				case sdl.K_RIGHT:
					if t.Repeat > 0 || t.State == sdl.RELEASED {
						break
					}

					nextImage := clamp(currentImage+1, 0, len(images)-1)
					if nextImage != currentImage {
						currentImage = nextImage
						image.Destroy()
						image = loadImage(images[currentImage], renderer)
					}
				case sdl.K_LEFT:
					if t.Repeat > 0 || t.State == sdl.RELEASED {
						break
					}

					nextImage := clamp(currentImage-1, 0, len(images)-1)
					if nextImage != currentImage {
						currentImage = nextImage
						image.Destroy()
						image = loadImage(images[currentImage], renderer)
					}
				}
			}
		}

		render(renderer, image)
	}
}
