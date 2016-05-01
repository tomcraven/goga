package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"math/rand"
	"os"
	"runtime"
	"time"
	ga "github.com/tomcraven/goga"
)

const (
	// Fiddle with these
	numShapes               = 100
	populationSize          = 10
	maxIterations           = 9999999
	bitsPerCoordinateNumber = 9 // smaller = more blocky picture (default 9)
	parallelSimulations     = 1 // number of simulations to run in parallel (default 4 (usually))
	maxCircleRadiusFactor   = 3 // larger = smaller max circle size relative to image dimensions (default 3)
	simulationAccuracy      = 1 // smaller = more accurate (min 1) (default 1)
	shapeSizeMultiplier     = 1.0 // larger = more fitness for smaller shapes (default 1.0)

	// Don't fiddle with these...
	maxBoxCornerCoordinateNumber = (1 << bitsPerCoordinateNumber) - 1
	bitsPerColourChannel         = 8 // 0 - 255
	bitsPerRect                  = (bitsPerCoordinateNumber * 4) + (bitsPerColourChannel * 4)
	bitsPerCircle                = (bitsPerCoordinateNumber * 3) + (bitsPerColourChannel * 4)
	bitsToDescribeWhichShape     = 1
)

type myBitsetCreate struct {
}

func (bc *myBitsetCreate) Go() ga.Bitset {
	b := ga.Bitset{}
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

	circleBitsetFormat ga.IBitsetParse
	rectBitsetFormat   ga.IBitsetParse
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

	rectBitsetFormat = ga.CreateBitsetParse()
	rectBitsetFormat.SetFormat([]int{
		bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber,
		bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel,
	})

	circleBitsetFormat = ga.CreateBitsetParse()
	circleBitsetFormat.SetFormat([]int{
		bitsPerCoordinateNumber, bitsPerCoordinateNumber, bitsPerCoordinateNumber,
		bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel, bitsPerColourChannel,
	})

	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Simulator = &imageMatcherSimulator{}
	genAlgo.BitsetCreate = &myBitsetCreate{}
	genAlgo.EliteConsumer = &myEliteConsumer{}
	genAlgo.Mater = ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: ga.UniformCrossover, UseElite: true},
			{P: 1.0, F: ga.Mutate},
		},
	)
	genAlgo.Selector = ga.NewSelector(
		[]ga.SelectorFunctionProbability{
			{P: 1.0, F: ga.Roulette},
		},
	)

	genAlgo.Init(populationSize, parallelSimulations)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
