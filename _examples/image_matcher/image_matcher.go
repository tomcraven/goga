package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"math/rand"
	"os"
	"runtime"
	"time"
)

const (
	// Fiddle with these
	numShapes               = 10
	populationSize          = 10
	maxIterations           = 9999999
	bitsPerCoordinateNumber = 9
	parallelSimulations     = 24
	maxCircleRadiusFactor   = 3 // larger == smaller max circle size relative to image dimensions

	// Don't fiddle with these...
	maxBoxCornerCoordinateNumber = (1 << bitsPerCoordinateNumber) - 1
	bitsPerColourChannel         = 8 // 0 - 255
	bitsPerRect                  = (bitsPerCoordinateNumber * 4) + (bitsPerColourChannel * 4)
	bitsPerCircle                = (bitsPerCoordinateNumber * 3) + (bitsPerColourChannel * 4)
	bitsToDescribeWhichShape     = 1
)

type imageMatcherSimulator struct {
	totalIterations int
}

func (simulator *imageMatcherSimulator) OnBeginSimulation() {
}
func (simulator *imageMatcherSimulator) OnEndSimulation() {
	simulator.totalIterations++
}
func (simulator *imageMatcherSimulator) Simulate(g goga.Genome) {
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
func (simulator *imageMatcherSimulator) ExitFunc(g goga.Genome) bool {
	return simulator.totalIterations >= maxIterations
}

type myBitsetCreate struct {
}

func (bc *myBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(numShapes * largestShapeBits)
	for i := 0; i < b.GetSize(); i++ {
		b.Set(i, rand.Intn(2))
	}
	return b
}

var (
	largestShapeBits   int
	totalBitsPerGenome int

	inputImage image.Image

	circleBitsetFormat goga.IBitsetParse
	rectBitsetFormat   goga.IBitsetParse
)

func init() {
	largestShapeBits = bitsToDescribeWhichShape
	if bitsPerRect > bitsPerCircle {
		largestShapeBits += bitsPerRect
	} else {
		largestShapeBits += bitsPerCircle
	}

	totalBitsPerGenome = largestShapeBits * numShapes
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

	rectBitsetFormat = goga.CreateBitsetParse()
	rectBitsetFormat.SetFormat([]int{
		bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber,
		bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel,
	})

	circleBitsetFormat = goga.CreateBitsetParse()
	circleBitsetFormat.SetFormat([]int{
		bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber,
		bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel,
	})

	genAlgo := goga.NewGeneticAlgorithm()
	genAlgo.Simulator = &imageMatcherSimulator{}
	genAlgo.BitsetCreate = &myBitsetCreate{}
	genAlgo.EliteConsumer = &myEliteConsumer{}
	genAlgo.Mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1.0, F: goga.UniformCrossover, UseElite: true},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.Mutate},
		},
	)
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 1.0, F: goga.Roulette},
		},
	)

	genAlgo.Init(populationSize, parallelSimulations)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
