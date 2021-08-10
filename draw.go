package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Renderer struct {
	Instance *sdl.Renderer
	Font     *ttf.Font
}

func (r *Renderer) GetStringSize(text string) (int32, int32) {
	width, height, err := r.Font.SizeUTF8(text)
	checkError(err)

	return int32(width), int32(height)
}

func (r *Renderer) DrawText(text string, rect *sdl.Rect, color sdl.Color) {
	surface, err := r.Font.RenderUTF8Blended(text, color)
	checkError(err)

	texture, err := r.Instance.CreateTextureFromSurface(surface)
	checkError(err)

	r.Instance.Copy(texture, nil, rect)

	surface.Free()
}

func (r *Renderer) DrawRectOutline(rect *sdl.Rect, color sdl.Color) {
	r.Instance.SetDrawColor(color.R, color.G, color.B, color.A)
	r.Instance.DrawRect(rect)
}

func (r *Renderer) DrawRect(rect *sdl.Rect, color sdl.Color) {
	r.Instance.SetDrawColor(color.R, color.G, color.B, color.A)
	r.Instance.FillRect(rect)
}
