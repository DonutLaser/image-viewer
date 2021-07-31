package main

import (
	"math"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Image struct {
	Data   *sdl.Texture
	Width  int32
	Height int32
}

func loadImage(path string, renderer *sdl.Renderer) (result Image) {
	image, err := img.Load(path)
	checkError(err)

	texture, err := renderer.CreateTextureFromSurface(image)
	checkError(err)

	result = Image{
		Data:   texture,
		Width:  image.W,
		Height: image.H,
	}

	image.Free()

	return
}

func (i *Image) Render(renderer *sdl.Renderer, screenWidth int32, screenHeight int32) {
	var finalRect *sdl.Rect
	if i.Width >= i.Height {
		newWidth := screenWidth
		newHeight := int32(math.Floor(float64(i.Height) / float64(i.Width) * float64(newWidth)))
		// @TODO do something about these casts

		finalRect = &sdl.Rect{
			X: 0,
			Y: screenHeight/2 - newHeight/2,
			W: newWidth,
			H: newHeight,
		}
	} else {
		newHeight := screenHeight
		newWidth := int32(math.Floor(float64(i.Width) / float64(i.Height) * float64(newHeight)))
		// @TODO do something about these casts

		finalRect = &sdl.Rect{
			X: screenWidth/2 - newWidth/2,
			Y: 0,
			W: newWidth,
			H: newHeight,
		}
	}

	renderer.Copy(i.Data, nil, finalRect)
}
