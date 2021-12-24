package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func createImageFromBitset(bits *goga.Bitset) draw.Image {
	inputImageBounds := inputImage.Bounds()

	newImage := image.NewRGBA(inputImageBounds)
	draw.Draw(newImage, newImage.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Over)

	for i := 0; i < bits.GetSize()/largestShapeBits; i++ {
		shapeBitset := bits.Slice(i*largestShapeBits, largestShapeBits)

		shapeType := shapeBitset.Get(0)
		if shapeType == 0 {
			rectBitset := shapeBitset.Slice(1, bitsPerRect)
			parsedBits := rectBitsetFormat.Process(&rectBitset)

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

			x1 := int((float64(parsedBits[0]) / float64(maxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.X))
			y1 := int((float64(parsedBits[1]) / float64(maxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.Y))
			x2 := int((float64(parsedBits[2]) / float64(maxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.X))
			y2 := int((float64(parsedBits[3]) / float64(maxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.Y))

			draw.DrawMask(newImage, image.Rect(x1, y1, x2, y2),
				&image.Uniform{colour}, image.ZP,
				&image.Uniform{alpha}, image.ZP,
				draw.Over)

		} else {
			circleBitset := shapeBitset.Slice(1, bitsPerCircle)
			parsedBits := circleBitsetFormat.Process(&circleBitset)

			colour := color.RGBA{
				uint8(parsedBits[3]),
				uint8(parsedBits[4]),
				uint8(parsedBits[5]),
				255,
			}

			normalisedX := float64(parsedBits[0]) / float64(maxBoxCornerCoordinateNumber)
			normalisedY := float64(parsedBits[1]) / float64(maxBoxCornerCoordinateNumber)

			xMin := float64(-inputImageBounds.Max.X)
			yMin := float64(-inputImageBounds.Max.Y)

			xMax := float64(inputImageBounds.Max.X + inputImageBounds.Max.X)
			yMax := float64(inputImageBounds.Max.Y + inputImageBounds.Max.Y)

			xRange := xMax - xMin
			yRange := yMax - yMin

			x := int(xMin + (normalisedX * xRange))
			y := int(yMin + (normalisedY * yRange))

			normalisedR := float64(parsedBits[2]) / float64(maxBoxCornerCoordinateNumber)
			maxR := math.Max(float64(inputImageBounds.Max.X), float64(inputImageBounds.Max.Y))
			r := int(normalisedR*maxR) / maxCircleRadiusFactor

			c := circle{image.Point{x, y}, r, uint8(parsedBits[6])}

			draw.DrawMask(newImage, inputImageBounds,
				&image.Uniform{colour}, image.ZP,
				&c, image.ZP,
				draw.Over)
		}
	}

	return newImage
}

// http://www.easyrgb.com/index.php?X=MATH
func rgbToXyz(r, g, b uint32) (float64, float64, float64) {
	normalizedR := float64(r) / 0xFFFF
	normalizedG := float64(g) / 0xFFFF
	normalizedB := float64(b) / 0xFFFF

	if normalizedR > 0.04045 {
		normalizedR = math.Pow(((normalizedR + 0.055) / 1.055), 2.4)
	} else {
		normalizedR = normalizedR / 12.92
	}

	if normalizedG > 0.04045 {
		normalizedG = math.Pow(((normalizedG + 0.055) / 1.055), 2.4)
	} else {
		normalizedG = normalizedG / 12.92
	}

	if normalizedB > 0.04045 {
		normalizedB = math.Pow(((normalizedB + 0.055) / 1.055), 2.4)
	} else {
		normalizedB = normalizedB / 12.92
	}

	normalizedR *= 100
	normalizedG *= 100
	normalizedB *= 100

	x := normalizedR*0.4124 + normalizedG*0.3576 + normalizedB*0.1805
	y := normalizedR*0.2126 + normalizedG*0.7152 + normalizedB*0.0722
	z := normalizedR*0.0193 + normalizedG*0.1192 + normalizedB*0.9505

	return x, y, z
}

// http://www.easyrgb.com/index.php?X=MATH
func xyzToLabAB(x, y, z float64) (float64, float64, float64) {
	normalizedX := x / 95.047
	normalizedY := y / 100.0
	normalizedZ := z / 108.883

	if normalizedX > 0.008856 {
		normalizedX = math.Pow(normalizedX, (1.0 / 3.0))
	} else {
		normalizedX = (7.787 * normalizedX) + (16.0 / 116.0)
	}

	if normalizedY > 0.008856 {
		normalizedY = math.Pow(normalizedY, (1.0 / 3.0))
	} else {
		normalizedY = (7.787 * normalizedY) + (16.0 / 116.0)
	}

	if normalizedZ > 0.008856 {
		normalizedZ = math.Pow(normalizedZ, (1.0 / 3.0))
	} else {
		normalizedZ = (7.787 * normalizedZ) + (16.0 / 116.0)
	}

	l := (116 * normalizedY) - 16
	a := 500 * (normalizedX - normalizedY)
	b := 200 * (normalizedY - normalizedZ)

	return l, a, b
}

func distance(a, b float64) float64 {
	return (a - b) * (a - b)
}

func calculateColourDifference(red1, green1, blue1, red2, green2, blue2 uint32) float64 {
	// First calculate XYZ
	x1, y1, z1 := rgbToXyz(red1, green1, blue1)
	x2, y2, z2 := rgbToXyz(red2, green2, blue2)

	// Then calculate CIE-L*ab
	l1, a1, b1 := xyzToLabAB(x1, y1, z1)
	l2, a2, b2 := xyzToLabAB(x2, y2, z2)

	// Calculate difference
	differences := distance(l1, l2) + distance(a1, a2) + distance(b1, b2)
	return math.Pow(differences, 0.5)
}
