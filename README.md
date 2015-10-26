# goga [![Build Status](https://travis-ci.org/tomcraven/goga.svg?branch=master)](https://travis-ci.org/tomcraven/goga) [![Coverage Status](https://coveralls.io/repos/tomcraven/goga/badge.svg?branch=master&service=github)](https://coveralls.io/github/tomcraven/goga?branch=master)

Golang implementation of a genetic algorithm. See ./examples for info on how to use the library.

## Overview
Goga is a genetic algorithm solution written in Golang. It is used and configured by injecting different behaviours into the main genetic algorithm object. The main injectable components are the simulator, selector and mater.

The simulator provides a function that accepts a single genome and assigns a fitness score to it. The higher the fitness, the better the genome has done in the simulation. A genome can be simulated by however the application sees fit as long as it can be encoded into a bitset of 0s and 1s. A simulator also provides a function to tell the algorithm when to stop.

The selector object takes a popualtion of genomes and the total fitness and returns a genome from the population that it has chosen. A common implementation is roulette in which a random value between 0..totalFitness is generated and the genomes are cycled through subtracting their fitness away from this random number. Then this number goes below 0 then a genome has been 'selected'. The idea is that a genome with a higher fitness will be more likely to be chosen.

A mater accepts two genomes from the selector and combines them to produce two others. There are some common predefined mating algorithms but the user is also free to define their own.

As genomes that have a fitness are more likely to mate, the program will slowly work its way towards what it thinks is an optimal solution.

## Examples
This section will talk through any example programs using this library.

#### examples/string_matcher.go
To run:
```
cd examples
go run string_matcher.go
```

The string matcher is a program that can be configured with any string and the goga library will attempt to generate a bitset that decodes to this string. Each character is representated by 8 bits in each genome.

There are some configuration constants just above the main function, use this to configure the string you'd like the algorithm to match to.
```
const (
	kTargetString = "abcdefghijklmnopqrstuvwxyz"
	// ...
)
```

The elite consumer is called after each simulation with the best genome of that iteration. It prints it out along with the iteration number and the fitness for that particular genome. A typical output might look something like this:
```
1 	 Øoâ|-7rPKw
                   D( 	 71
2 	 Xnæ=írVÏw
T3 	 74
3 	 oî,ë.2wës
# 	 81
4 	 XOà|m,WOz
D" 	 84
5 	 XOålo,rWorD# 	 89
6 	 XOìlo.0WozD# 	 91
7 	 Xgflm,pWorND# 	 93
8 	 Xgllo,0WorD# 	 96
9 	 Xgllo.0WorNd# 	 97
10 	 Xello.0WorNd# 	 98
11 	 Hello,0WorNd# 	 100
12 	 Hello,0WorNd# 	 100
13 	 Hello, WorNd# 	 101
14 	 Hello, Wornd! 	 103
15 	 Hello, Wornd! 	 103
16 	 Hello, World! 	 104
151.88245ms
```
You can see the string slowly becoming more like the input string as there are more iterations.

#### examples/image_matcher.go
To run:
```
cd examples
go run image_matcher.go <path_to_image>
```

Image matcher takes an input image and attempts to produce an output image that is as close to it as possible only using RGBA coloured rectangles and circles. There are a few parameters at the top of the file that are interesting to fiddle with:
```
kNumShapes = 100
kPopulationSize = 1000
kMaxIterations = 9999999
kBitsPerCoordinateNumber = 9
kParallelSimulations = 24
kMaxCircleRadiusFactor = 3
```
* kNumShapes - the number of shapes that are used when the algorithm re-creates the input image
* kPopulationSize - the number of genomes in each population. Each genome can be decoded into a picture. A high value will mean each iteration takes longer, and usually results in the algorithm finding its optimal solution in less iterations.
* kMaxIterations - the maximum number of simulations/iterations to run. Providing a huge number will essentially run until the algorithm has figured out what it thinks is an optimal solution.
* kBitsPerCoordinateNumber - Each shape is positioned using coordinates. A rect is represented by the top left and bottom right coordinates, and a circle by its centre. A coordinate is made up of two numbers, each number is represented by this many bits. The number generated is used to calculate a percentage of the overall width/height of the image for the coordinate to be positioned at. For example, if kBitsPerCoordinateNumber is 8, that means the maximum value a coordinte can be is ```0b11111111```, or ```255```. To calculate the coordinates number relative to the image's width and height we normalise this and apply the decimal to the pictures dimensions. For example, if our X coordinate produced by the algorithm is 233, and our images width is 120. ```( 233 / 255 ) * 120 == 109 == our X coordinate```. Setting this to a high value means the algorithm has more accuracy when placing shapes. A low value of 2 or 3 also creates some interesting effects.
* kParallelSimulations - The number of simulations to run in parallel. A higher value usually means each iteration takes less time, up to a certain point.
* kMaxCircleRadiusFactor - A number to divide a circles radius by. 1 will mean a circles radius can be anywhere between 1 and max( inputImageWidth, inputImageHeight ). From experimenting, restricting the radius a little creates 'better' images quicker

The script outputs the 'best' genome from each iteration to "elite.png" as well as the original overlayed in "elite_with_original.png".

Here's an example of what can be generated. Input, output and both overlayed over the top of each other:
![input](https://cloud.githubusercontent.com/assets/5236109/10744734/01031bda-7c34-11e5-94ab-795afba114c1.gif)
![output](https://cloud.githubusercontent.com/assets/5236109/10744673/97a7aea8-7c33-11e5-8cfe-ea66489d8d9c.png)
![overlayed](https://cloud.githubusercontent.com/assets/5236109/10744674/9ad23fa8-7c33-11e5-88d9-aff565cca6c4.png)