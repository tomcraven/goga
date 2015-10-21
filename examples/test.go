package main

import (
	"image/draw"
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	newImage := image.NewRGBA( image.Rect( 0, 0, 100, 100 ) )
	// draw.Draw(newImage, newImage.Bounds(), &image.Uniform{ color.RGBA{ 0, 0, 0, 255 } }, image.ZP, draw.Over)	

	colour := color.RGBA{0, 0, 255, 255}
	colour2 := color.RGBA{255, 255, 255, 255 / 2}

	draw.DrawMask(newImage, image.Rect( 10, 10, 50, 50 ), 
		&image.Uniform{ colour }, image.ZP, &image.Uniform{ colour2 }, image.ZP, draw.Over)

	outputImageFile, _ := os.Create( "elite.png" )
    defer outputImageFile.Close()
    png.Encode( outputImageFile, newImage )
}
