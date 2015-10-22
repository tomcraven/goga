
package main

import (
	"github.com/tomcraven/goga"
	"fmt"
	"math/rand"
	"time"
	"image"
	"image/color"
	"image/png"
	"image/draw"
	"os"
	_ "image/jpeg"
	"math"
	"runtime"
)

func swap( a, b *int ) {
	temp := *a

	*a = *b
	*b = temp
}

func createImageFromBitset( bits *ga.Bitset, bitsetFormat ga.IBitsetParse ) image.Image {
	inputImageBounds := inputImage.Bounds()

	newImage := image.NewRGBA( inputImageBounds )
	// draw.Draw(newImage, newImage.Bounds(), &image.Uniform{ color.RGBA{ 0, 0, 0, 255 } }, image.ZP, draw.Over)	

	for i := 0; i < bits.GetSize() / kBitsPerBox; i++ {
		boxBitset := bits.Slice( i * kBitsPerBox, kBitsPerBox )
		parsedBits := bitsetFormat.Process( &boxBitset )

		colour := color.RGBA{
			uint8( parsedBits[4] ),
			uint8( parsedBits[5] ),
			uint8( parsedBits[6] ),
			255,
		}

		alpha := color.RGBA{
			255, 255, 255,
			uint8( parsedBits[7] ),
		}

		x1 := int( ( float64( parsedBits[0] ) / float64( kMaxBoxCornerCoordinate ) ) * float64( inputImageBounds.Max.X ) )
		y1 := int( ( float64( parsedBits[1] ) / float64( kMaxBoxCornerCoordinate ) ) * float64( inputImageBounds.Max.Y ) )
		x2 := int( ( float64( parsedBits[2] ) / float64( kMaxBoxCornerCoordinate ) ) * float64( inputImageBounds.Max.X ) )
		y2 := int( ( float64( parsedBits[3] ) / float64( kMaxBoxCornerCoordinate ) ) * float64( inputImageBounds.Max.Y ) )

		// if x1 > x2 {
		// 	swap( &x1, &x2 )
		// }

		// if ( y1 > y2 ) {
		// 	swap( &y1, &y2 )
		// }

		draw.DrawMask(newImage, image.Rect( x1, y1, x2, y2 ),
			&image.Uniform{ colour }, image.ZP, 
			&image.Uniform{ alpha }, image.ZP, 
			draw.Over)
	}

	return newImage
}

func minFloat64( a, b float64 ) float64 {
	if ( a < b ) {
		return a
	}

	return b
}

func maxFloat64( a, b float64 ) float64 {
	if ( a > b ) {
		return a
	}

	return b
}

func calculateHue( r, g, b uint32 ) float64 {
	normalisedR := float64( r ) / 0xffff
	normalisedG := float64( g ) / 0xffff
	normalisedB := float64( b ) / 0xffff

	min := minFloat64( normalisedR, minFloat64( normalisedG, normalisedB ) )
	max := maxFloat64( normalisedR, maxFloat64( normalisedG, normalisedB ) )

	if( ( max - min ) == 0 ) {
		return 0.0
	}

	hue := 0.0
	if ( ( normalisedR > normalisedG ) && ( normalisedR > normalisedB ) ) {
		hue = ( normalisedG - normalisedB ) / ( max - min )
	} else if ( ( normalisedG > normalisedR ) && ( normalisedG > normalisedB ) ) {
		hue = 2.0 + ( ( normalisedB - normalisedR ) / ( max - min ) )
	} else {
		hue = 4.0 + ( ( normalisedR - normalisedG ) / ( max - min ) )
	}

	hue *= 60
	if ( hue < 0.0 ) {
		hue += 360
	}

	return hue
}
type ImageMatcherSimulator struct {
	BitsetFormat ga.IBitsetParse
	totalIterations int
}
func ( simulator *ImageMatcherSimulator ) OnBeginSimulation() {
}
func ( simulator *ImageMatcherSimulator ) OnEndSimulation() {
	simulator.totalIterations++
}
func ( simulator *ImageMatcherSimulator ) Simulate( g *ga.IGenome ) {

	bits := (*g).GetBits()
	newImage := createImageFromBitset( bits, simulator.BitsetFormat )

	inputImageBounds := inputImage.Bounds()
	fitness := 0.0
	for y := 0; y < inputImageBounds.Max.Y; y++ {
		for x := 0; x < inputImageBounds.Max.X; x++ {
			inputR, inputG, inputB, inputA := inputImage.At( x, y ).RGBA()
			createdR, createdG, createdB, createdA := newImage.At( x, y ).RGBA()

			inputHue := calculateHue( inputR, inputG, inputB )
			createdHue := calculateHue( createdR, createdG, createdB )

			hueDifference := math.Abs( inputHue - createdHue )
			alphaDifference := math.Abs( float64( createdA ) - float64( inputA ) )

			fitness += ( 360.0 - hueDifference ) + ( 0xFFFF - alphaDifference )
		}
	}

	(*g).SetFitness( int( fitness ) )
}
func ( simulator *ImageMatcherSimulator ) ExitFunc( g *ga.IGenome ) bool {
	return simulator.totalIterations >= kMaxIterations
}

type MyBitsetCreate struct {
}
func ( bc *MyBitsetCreate ) Go() ga.Bitset {
	b := ga.Bitset{}
	b.Create( kNumBoxes * kBitsPerBox )
	for i := 0; i < b.GetSize(); i++ {
		b.Set( i, rand.Intn( 2 ) )
	}
	return b
}

type MyEliteConsumer struct {
	currentIter int
	BitsetFormat ga.IBitsetParse
	previousFitness int
}
func ( ec *MyEliteConsumer ) OnElite( g *ga.IGenome ) {
	bits := (*g).GetBits()
	newImage := createImageFromBitset( bits, ec.BitsetFormat )

	outputImageFile, _ := os.Create( "elite.png" )
    defer outputImageFile.Close()
    png.Encode( outputImageFile, newImage )

	ec.currentIter++
	fitness := (*g).GetFitness()
	fmt.Println( ec.currentIter, "\t", fitness, "\t", fitness - ec.previousFitness )

	ec.previousFitness = fitness
}

const (
	// Fiddle with these
	kNumBoxes = 4
	kPopulationSize = 10
	kMaxIterations = 99999999

	// Don't fiddle with these...
	kBitsPerCorner = 3
	kMaxBoxCornerCoordinate = ( 1 << kBitsPerCorner ) - 1
	kBitsPerColourChannel = 8	// 0 - 255
	kBitsPerBox = ( kBitsPerCorner * 4 ) + ( kBitsPerColourChannel * 4 )
	kTotalBitsPerGenome = kBitsPerBox * kNumBoxes
)

var (
	inputImage image.Image
)

func getImageFromFile( filename string ) image.Image {
	inputImageFile, _ := os.Open( filename )
	defer inputImageFile.Close()
	inputImage, _, _ := image.Decode( inputImageFile )
	return inputImage
}

func main() {

	runtime.GOMAXPROCS( 4 )

	// Get the input image
	inputImage = getImageFromFile( os.Args[ 1 ] )

	genAlgo := ga.NewGeneticAlgorithm()

	imageMatcherSimulator := ImageMatcherSimulator{}
	imageMatcherSimulator.BitsetFormat = ga.CreateBitsetParse()
	imageMatcherSimulator.BitsetFormat.SetFormat( []int{
			kBitsPerCorner, kBitsPerCorner, kBitsPerCorner, kBitsPerCorner,
			kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel, kBitsPerColourChannel,
		})

	genAlgo.Simulator = &imageMatcherSimulator
	genAlgo.BitsetCreate = &MyBitsetCreate{}

	eliteConsumer := MyEliteConsumer{}
	eliteConsumer.BitsetFormat = imageMatcherSimulator.BitsetFormat
	genAlgo.EliteConsumer = &eliteConsumer
	genAlgo.Mater = ga.NewMater( 
		[]ga.MaterFunctionProbability{
			{ P : 1.0, F : ga.UniformCrossover, UseElite : true },
			{ P : 1.0, F : ga.Mutate },
			{ P : 1.0, F : ga.Mutate },
			{ P : 1.0, F : ga.Mutate },
			{ P : 1.0, F : ga.Mutate },
			{ P : 1.0, F : ga.Mutate },
		},
	)
	genAlgo.Selector = ga.NewSelector(
		[]ga.SelectorFunctionProbability {
			{ P : 1.0, F : ga.Roulette },
		},
	)

	genAlgo.Init( kPopulationSize )

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println( time.Since( startTime ) )
}