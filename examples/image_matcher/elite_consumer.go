package main

import (
	ga "github.com/tomcraven/goga"
	"os"
	_ "image/jpeg"
	"image/png"
	"image/draw"
	"image"
	"image/color"
	"fmt"
	"time"
)

type myEliteConsumer struct {
	currentIter     int
	previousFitness int
}

func (ec *myEliteConsumer) OnElite(g *ga.IGenome) {
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
