package waveshare7in5v2

import (
	"image"
	"image/color"
)

type Canvas struct {
	e *Epd

	img *image.Gray
}

func NewCanvas(e *Epd) *Canvas {
	img := image.NewGray(e.bounds)

	c := &Canvas{
		e: e,

		img: img,
	}

	return c
}

func (c *Canvas) At(x, y int) color.Color {
	return c.img.At(x, y)
}

func (c *Canvas) Bounds() image.Rectangle {
	return c.img.Bounds()
}

func (c *Canvas) ColorModel() color.Model {
	return c.img.ColorModel()
}

func (c *Canvas) Set(x, y int, color color.Color) {
	c.img.Set(x, y, color)
}

func (c *Canvas) Refresh() {
	buffer := c.e.GetBuffer(c.img)

	c.e.UpdateFrameAndRefresh(buffer)
}

func (c *Canvas) Init() {
	c.e.Init()
}

func (c *Canvas) Clear() {
	c.e.Clear()
}

func (c *Canvas) Sleep() {
	c.e.Sleep()
}
