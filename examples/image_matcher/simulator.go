package main

import (
	"time"
	"fmt"
	ga "github.com/tomcraven/goga"
)

type imageMatcherSimulator struct {
	totalIterations int
}

func (simulator *imageMatcherSimulator) OnBeginSimulation() {
}
func (simulator *imageMatcherSimulator) OnEndSimulation() {
	simulator.totalIterations++
}
func (simulator *imageMatcherSimulator) Simulate(g *ga.Genome) {
	bits := (*g).GetBits()

	currentTime := time.Now()

	newImage := createImageFromBitset(bits)

	fmt.Printf("%v\t", time.Since(currentTime))
	currentTime = time.Now()

	maxColourDifference := 0.0	
	inputImageBounds := inputImage.Bounds()
	fitness := 0.0
	totalSeconds := 0.0
	maxFitness := 32768.0 // calculated through trial and error
	for y := 0; y < inputImageBounds.Max.Y; y+=simulationAccuracy {
		for x := 0; x < inputImageBounds.Max.X; x+=simulationAccuracy {
			
			inputR, inputG, inputB, _ := inputImage.At(x, y).RGBA()
			createdR, createdG, createdB, _ := newImage.At(x, y).RGBA()

			startDifferenceTime := time.Now()
			colourDifference := calculateColourDifference(inputR, inputG, inputB, createdR, createdG, createdB)
			totalSeconds += time.Since(startDifferenceTime).Seconds()

			if colourDifference > maxColourDifference {
				maxColourDifference = colourDifference
			}
			
			fitness += (maxFitness  - colourDifference)
		}
	}

	numPixels := float64((inputImageBounds.Max.Y/simulationAccuracy) * (inputImageBounds.Max.X/simulationAccuracy))
	fmt.Printf("%v(%v)\t", totalSeconds / numPixels, totalSeconds)
	fmt.Printf("%v %v\t", numPixels, maxColourDifference)

	// Award fitness for smaller shapes
	sf := shapeFactory{
		inputImageBounds: inputImageBounds,
	}
	for i := 0; i < numShapes; i++ {
		shapeBitset := bits.Slice(i*largestShapeBits, largestShapeBits)
		shape := sf.create(&shapeBitset)
		fitness += shapeSizeMultiplier * float64(inputImageBounds.Max.X - shape.width())
		fitness += shapeSizeMultiplier * float64(inputImageBounds.Max.Y - shape.height())
	}

	fmt.Println(time.Since(currentTime))

	(*g).SetFitness(int(fitness))
}
func (simulator *imageMatcherSimulator) ExitFunc(g *ga.Genome) bool {
	return simulator.totalIterations >= maxIterations
}
