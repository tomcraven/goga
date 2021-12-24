package main

import (
	"image"
	"image/color"
)

// http://blog.golang.org/go-imagedraw-package
type circle struct {
	p     image.Point
	r     int
	alpha uint8
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{c.alpha}
	}
	return color.Alpha{0}
}
