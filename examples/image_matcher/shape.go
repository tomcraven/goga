package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	ga "github.com/tomcraven/goga"
)

type shapeFactory struct {
	inputImageBounds image.Rectangle
}

type shape interface {
	render(*image.RGBA)
	width() int
	height() int
}

type rectShape struct {
	x1, y1, x2, y2 int
	colour color.RGBA
	alpha color.RGBA
}

type circleShape struct {
	x, y, r int
	colour color.RGBA
	alpha uint8
	inputImageBounds image.Rectangle
}

func(sf *shapeFactory) create(bits *ga.Bitset) shape {
	shapeType := bits.Get(0)
	if shapeType == 0 {
		return sf.createRect(bits.Slice(1, bitsPerRect))
	} else {
		return sf.createCircle(bits.Slice(1, bitsPerCircle))
	}
}

func (sf *shapeFactory) createRect(bits ga.Bitset) shape {
	parsedBits := rectBitsetFormat.Process(&bits)

	colour := color.RGBA{
		uint8(parsedBits[4]),
		uint8(parsedBits[5]),
		uint8(parsedBits[6]),
		255,
	}

	alpha := color.RGBA{
		255, 255, 255,
		uint8(parsedBits[7]),
	}

	x1 := int((float64(parsedBits[0]) / float64(maxBoxCornerCoordinateNumber)) * float64(sf.inputImageBounds.Max.X))
	y1 := int((float64(parsedBits[1]) / float64(maxBoxCornerCoordinateNumber)) * float64(sf.inputImageBounds.Max.Y))
	x2 := int((float64(parsedBits[2]) / float64(maxBoxCornerCoordinateNumber)) * float64(sf.inputImageBounds.Max.X))
	y2 := int((float64(parsedBits[3]) / float64(maxBoxCornerCoordinateNumber)) * float64(sf.inputImageBounds.Max.Y))

	return &rectShape {
		x1: x1,
		y1: y1,
		x2: x2,
		y2: y2,
		colour: colour,
		alpha: alpha,
	}
}

func (sf *shapeFactory) createCircle(bits ga.Bitset) shape {
	parsedBits := circleBitsetFormat.Process(&bits)

	colour := color.RGBA{
		uint8(parsedBits[3]),
		uint8(parsedBits[4]),
		uint8(parsedBits[5]),
		255,
	}

	normalisedX := float64(parsedBits[0]) / float64(maxBoxCornerCoordinateNumber)
	normalisedY := float64(parsedBits[1]) / float64(maxBoxCornerCoordinateNumber)

	xMin := float64(-sf.inputImageBounds.Max.X)
	yMin := float64(-sf.inputImageBounds.Max.Y)

	xMax := float64(sf.inputImageBounds.Max.X + sf.inputImageBounds.Max.X)
	yMax := float64(sf.inputImageBounds.Max.Y + sf.inputImageBounds.Max.Y)

	xRange := xMax - xMin
	yRange := yMax - yMin

	x := int(xMin + (normalisedX * xRange))
	y := int(yMin + (normalisedY * yRange))

	normalisedR := float64(parsedBits[2]) / float64(maxBoxCornerCoordinateNumber)
	maxR := math.Max(float64(sf.inputImageBounds.Max.X), float64(sf.inputImageBounds.Max.Y))
	r := int(normalisedR*maxR) / maxCircleRadiusFactor

	// TODO: we're casting from uint64 to uint8 here
	// should be fine, but should check too
	alpha := uint8(parsedBits[6])

	return &circleShape{
		x: x,
		y: y,
		r: r,
		alpha: alpha,
		colour: colour,
		inputImageBounds: sf.inputImageBounds,
	}
}

func (r *rectShape) render(newImage *image.RGBA) {
	draw.DrawMask(newImage, image.Rect(r.x1, r.y1, r.x2, r.y2),
		&image.Uniform{r.colour}, image.ZP,
		&image.Uniform{r.alpha}, image.ZP,
		draw.Over)
}

func intAbs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func (r *rectShape) width() int {
	return intAbs(r.x1 - r.x2)
}

func (r *rectShape) height() int {
	return intAbs(r.y1 - r.y2)
}

func (c *circleShape) render(newImage *image.RGBA) {
	mask := circle{image.Point{c.x, c.y}, c.r, c.alpha}
	draw.DrawMask(newImage, c.inputImageBounds,
		&image.Uniform{c.colour}, image.ZP,
		&mask, image.ZP,
		draw.Over)
}

func (c *circleShape) width() int {
	return c.diameter()
}

func (c *circleShape) height() int {
	return c.diameter()
}

func (c *circleShape) diameter() int {
	return 2 * c.r
}
