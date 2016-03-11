package main

import (
	"fmt"
	"github.com/tomcraven/goga"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const (
	// Fiddle with these
	kNumShapes               = 100
	kPopulationSize          = 100
	kMaxIterations           = 9999999
	kBitsPerCoordinateNumber = 9
	kParallelSimulations     = 24
	kMaxCircleRadiusFactor   = 3 // larger == smaller max circle size relative to image dimensions

	// Don't fiddle with these...
	kMaxBoxCornerCoordinateNumber = (1 << kBitsPerCoordinateNumber) - 1
	kBitsPerColourChannel         = 8 // 0 - 255
	kBitsPerRect                  = (kBitsPerCoordinateNumber * 4) + (kBitsPerColourChannel * 4)
	kBitsPerCircle                = (kBitsPerCoordinateNumber * 3) + (kBitsPerColourChannel * 4)
	kBitsToDescribeWhichShape     = 1
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

func createImageFromBitset(bits *ga.Bitset) draw.Image {
	inputImageBounds := inputImage.Bounds()

	newImage := image.NewRGBA(inputImageBounds)
	draw.Draw(newImage, newImage.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Over)

	for i := 0; i < bits.GetSize()/kLargestShapeBits; i++ {
		shapeBitset := bits.Slice(i*kLargestShapeBits, kLargestShapeBits)

		shapeType := shapeBitset.Get(0)
		if shapeType == 0 {
			rectBitset := shapeBitset.Slice(1, kBitsPerRect)
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

			x1 := int((float64(parsedBits[0]) / float64(kMaxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.X))
			y1 := int((float64(parsedBits[1]) / float64(kMaxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.Y))
			x2 := int((float64(parsedBits[2]) / float64(kMaxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.X))
			y2 := int((float64(parsedBits[3]) / float64(kMaxBoxCornerCoordinateNumber)) * float64(inputImageBounds.Max.Y))

			draw.DrawMask(newImage, image.Rect(x1, y1, x2, y2),
				&image.Uniform{colour}, image.ZP,
				&image.Uniform{alpha}, image.ZP,
				draw.Over)

		} else {
			circleBitset := shapeBitset.Slice(1, kBitsPerCircle)
			parsedBits := circleBitsetFormat.Process(&circleBitset)

			colour := color.RGBA{
				uint8(parsedBits[3]),
				uint8(parsedBits[4]),
				uint8(parsedBits[5]),
				255,
			}

			normalisedX := float64(parsedBits[0]) / float64(kMaxBoxCornerCoordinateNumber)
			normalisedY := float64(parsedBits[1]) / float64(kMaxBoxCornerCoordinateNumber)

			xMin := float64(-inputImageBounds.Max.X)
			yMin := float64(-inputImageBounds.Max.Y)

			xMax := float64(inputImageBounds.Max.X + inputImageBounds.Max.X)
			yMax := float64(inputImageBounds.Max.Y + inputImageBounds.Max.Y)

			xRange := xMax - xMin
			yRange := yMax - yMin

			x := int(xMin + (normalisedX * xRange))
			y := int(yMin + (normalisedY * yRange))

			normalisedR := float64(parsedBits[2]) / float64(kMaxBoxCornerCoordinateNumber)
			maxR := math.Max(float64(inputImageBounds.Max.X), float64(inputImageBounds.Max.Y))
			r := int(normalisedR*maxR) / kMaxCircleRadiusFactor

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

type ImageMatcherSimulator struct {
	totalIterations int
}

func (simulator *ImageMatcherSimulator) OnBeginSimulation() {
}
func (simulator *ImageMatcherSimulator) OnEndSimulation() {
	simulator.totalIterations++
}
func (simulator *ImageMatcherSimulator) Simulate(g *ga.IGenome) {

	bits := (*g).GetBits()
	newImage := createImageFromBitset(bits)

	inputImageBounds := inputImage.Bounds()
	fitness := 0.0
	for y := 0; y < inputImageBounds.Max.Y; y++ {
		for x := 0; x < inputImageBounds.Max.X; x++ {
			inputR, inputG, inputB, _ := inputImage.At(x, y).RGBA()
			createdR, createdG, createdB, _ := newImage.At(x, y).RGBA()

			colourDifference := calculateColourDifference(inputR, inputG, inputB, createdR, createdG, createdB)

			fitness += (500.0 - colourDifference)
		}
	}

	(*g).SetFitness(int(fitness))
}
func (simulator *ImageMatcherSimulator) ExitFunc(g *ga.IGenome) bool {
	return simulator.totalIterations >= kMaxIterations
}

type MyBitsetCreate struct {
}

func (bc *MyBitsetCreate) Go() ga.Bitset {
	b := ga.Bitset{}
	b.Create(kNumShapes * kLargestShapeBits)
	for i := 0; i < b.GetSize(); i++ {
		b.Set(i, rand.Intn(2))
	}
	return b
}

type MyEliteConsumer struct {
	currentIter     int
	previousFitness int
}

func (ec *MyEliteConsumer) OnElite(g *ga.IGenome) {
	bits := (*g).GetBits()
	newImage := createImageFromBitset(bits)

	// Output elite
	outputImageFile, _ := os.Create("elite.png")
	png.Encode(outputImageFile, newImage)
	outputImageFile.Close()

	// Output elite with input image blended over the top
	outputImageFileAlphaBlended, _ := os.Create("elite_with_original.png")
	draw.DrawMask(newImage, newImage.Bounds(),
		inputImage, image.ZP,
		&image.Uniform{color.RGBA{0, 0, 0, 255 / 4}}, image.ZP,
		draw.Over)
	png.Encode(outputImageFileAlphaBlended, newImage)
	outputImageFileAlphaBlended.Close()

	ec.currentIter++
	fitness := (*g).GetFitness()
	fmt.Println(ec.currentIter, "\t", fitness, "\t", fitness-ec.previousFitness)

	ec.previousFitness = fitness

	time.Sleep(10 * time.Millisecond)
}

var (
	kLargestShapeBits   int
	kTotalBitsPerGenome int

	inputImage image.Image

	circleBitsetFormat ga.IBitsetParse
	rectBitsetFormat   ga.IBitsetParse
)

func init() {
	kLargestShapeBits = kBitsToDescribeWhichShape
	if kBitsPerRect > kBitsPerCircle {
		kLargestShapeBits += kBitsPerRect
	} else {
		kLargestShapeBits += kBitsPerCircle
	}

	kTotalBitsPerGenome = kLargestShapeBits * kNumShapes
}

func getImageFromFile(filename string) image.Image {
	inputImageFile, _ := os.Open(filename)
	defer inputImageFile.Close()
	inputImage, _, _ := image.Decode(inputImageFile)
	return inputImage
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	// Get the input image
	inputImage = getImageFromFile(os.Args[1])

	rectBitsetFormat = ga.CreateBitsetParse()
	rectBitsetFormat.SetFormat([]int{
		kBitsPerCoordinateNumber, kBitsPerCoordinateNumber, kBitsPerCoordinateNumber, kBitsPerCoordinateNumber,
		kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel,
	})

	circleBitsetFormat = ga.CreateBitsetParse()
	circleBitsetFormat.SetFormat([]int{
		kBitsPerCoordinateNumber, kBitsPerCoordinateNumber, kBitsPerCoordinateNumber,
		kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel,
	})

	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Simulator = &ImageMatcherSimulator{}
	genAlgo.BitsetCreate = &MyBitsetCreate{}
	genAlgo.EliteConsumer = &MyEliteConsumer{}
	genAlgo.Mater = ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: ga.UniformCrossover, UseElite: true},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.Mutate},
		},
	)
	genAlgo.Selector = ga.NewSelector(
		[]ga.SelectorFunctionProbability{
			{P: 1.0, F: ga.Roulette},
		},
	)

	genAlgo.Init(kPopulationSize, kParallelSimulations)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
