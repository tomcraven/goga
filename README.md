# goga [![Build Status](https://travis-ci.org/tomcraven/goga.svg?branch=master)](https://travis-ci.org/tomcraven/goga) [![Coverage Status](https://coveralls.io/repos/tomcraven/goga/badge.svg?branch=master&service=github)](https://coveralls.io/github/tomcraven/goga?branch=master)

Golang implementation of a genetic algorithm. See ./examples for info on how to use the library.

## Usage
Goga is configured by injecting different behaviours into the main genetic algorithm object. The main injectable components are the simulator, selector and mater.

The simulator provides a function that accepts a single genome and assigns a fitness score to it. The higher the fitness, the better the genome has done in the simulation. A genome can be simulated by however the application sees fit as long as it can be encoded into a bitset of 0s and 1s. A simulator also provides a function to tell the algorithm when to stop.

The selector object takes a popualtion of genomes and the total fitness and returns a genome from the population that it has chosen. A common implementation is roulette in which a random value between 0..totalFitness is generated and the genomes are cycled through subtracting their fitness away from this random number. Then this number goes below 0 then a genome has been 'selected'. The idea is that a genome with a higher fitness will be more likely to be chosen.

A mater accepts two genomes and combines them to produce two others. The idea is that a 

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
