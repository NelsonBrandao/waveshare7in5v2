package waveshare7in5v2

import (
	"image"
	"image/color"
)

// The Canvas implements draw.Image and thus allows to use any compatible
// package to draw on the screen.
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

// Flushes any changes done locally and updates the display
func (c *Canvas) Refresh() {
	c.e.DisplayImage(c.img)
}

// Clear the buffer and updates the screen right away.
func (c *Canvas) Clear() {
	c.e.Clear()
}
