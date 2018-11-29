package main

import "github.com/veandco/go-sdl2/sdl"

type stringFunc func() string
type colorFunc func() sdl.Color

type item interface {
	getSurfaceAndRect() (*sdl.Surface, *sdl.Rect, error)
}

type textItem struct {
	text  stringFunc
	color colorFunc
	x     int32
	y     int32
}

func (ti textItem) getSurfaceAndRect() (*sdl.Surface, *sdl.Rect, error) {
	surface, err := font.RenderUTF8Blended(ti.text(), ti.color())

	rect := &sdl.Rect{
		H: surface.H,
		W: surface.W,
		X: ti.x,
		Y: ti.y,
	}

	return surface, rect, err
}
