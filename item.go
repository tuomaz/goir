package main

import "github.com/veandco/go-sdl2/sdl"

type stringFunc func() string
type colorFunc func() sdl.Color

type item interface {
	getSurfaceAndRect() (*sdl.Surface, *sdl.Rect, error)
}

type textItem struct {
	text    stringFunc
	color   colorFunc
	bgcolor colorFunc
	x       int32
	y       int32
}

func (ti textItem) getSurfaceAndRect() (*sdl.Surface, *sdl.Rect, error) {
	var surface *sdl.Surface
	var err error
	if ti.bgcolor == nil {
		surface, err = font.RenderUTF8Blended(ti.text(), ti.color())
	} else {
		surface, err = font.RenderUTF8Shaded(ti.text(), ti.color(), ti.bgcolor())
	}
	rect := &sdl.Rect{
		H: surface.H,
		W: surface.W,
		X: ti.x,
		Y: ti.y,
	}

	return surface, rect, err
}
