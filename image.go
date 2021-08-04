package main

import (
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
		ratio := float64(i.Height) / float64(i.Width)
		newWidth := min(screenWidth, i.Width)
		newHeight := floor(float64(newWidth) * ratio)

		if newHeight > screenHeight {
			inverseRatio := float64(i.Width) / float64(i.Height)
			newHeight = screenHeight
			newWidth = floor(float64(screenHeight) * inverseRatio)
		}

		finalRect = &sdl.Rect{
			X: max((screenWidth-newWidth)/2, 0),
			Y: screenHeight/2 - newHeight/2,
			W: newWidth,
			H: newHeight,
		}
	} else {
		ratio := float64(i.Width) / float64(i.Height)
		newHeight := min(screenHeight, i.Height)
		newWidth := floor(float64(newHeight) * ratio)

		if newWidth > screenWidth {
			inverseRatio := float64(i.Height) / float64(i.Width)
			newWidth = screenWidth
			newHeight = floor(float64(screenWidth) * inverseRatio)
		}

		finalRect = &sdl.Rect{
			X: screenWidth/2 - newWidth/2,
			Y: max((screenHeight-newHeight)/2, 0),
			W: newWidth,
			H: newHeight,
		}
	}

	renderer.Copy(i.Data, nil, finalRect)
}
