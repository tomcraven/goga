
package main

import (
	"github.com/tomcraven/goga"
	"fmt"
	"math/rand"
	"time"
	"runtime"
)

type StringMaterSimulator struct {
}
func ( sms *StringMaterSimulator ) OnBeginSimulation() {
}
func ( sms *StringMaterSimulator ) OnEndSimulation() {
}
func ( sms *StringMaterSimulator ) Simulate( g *ga.IGenome ) {	
	bits := (*g).GetBits()
	for i, character := range kTargetString {
		for j := 0; j < 8; j++ {
			targetBit := character & (1 << uint(j))
			bit := bits.Get(( i * 8 ) + j)
			if targetBit != 0 && bit == 1 {
				(*g).SetFitness( (*g).GetFitness() + 1 )
			} else if targetBit == 0 && bit == 0 {
				(*g).SetFitness( (*g).GetFitness() + 1 )
			}
		}
	}
}
func ( sms *StringMaterSimulator ) ExitFunc( g *ga.IGenome ) bool {
	return (*g).GetFitness() == kTargetLength
}

type MyBitsetCreate struct {
}
func ( bc *MyBitsetCreate ) Go() ga.Bitset {
	b := ga.Bitset{}
	b.Create( kTargetLength )
	for i := 0; i < kTargetLength; i++ {
		b.Set( i, rand.Intn( 2 ) )
	}
	return b
}

type MyEliteConsumer struct {
	currentIter int
}
func ( ec *MyEliteConsumer ) OnElite( g *ga.IGenome ) {
	gBits := (*g).GetBits()
	ec.currentIter++
	var genomeString string
	for i := 0; i < gBits.GetSize(); i += 8 {
		c := int(0)
		for j := 0; j < 8; j++ {
			bit := gBits.Get( i + j )
			if bit != 0 {
				c |= 1 << uint( j )
			}
		}
		genomeString += string( c )
	}

	fmt.Println( ec.currentIter, "\t", genomeString, "\t", (*g).GetFitness() )
}

const (
	kTargetString = "abcdefghijklmnopqrstuvwxyz"
	kTargetLength = len( kTargetString ) * 8
	kPopulationSize = 600
)

func main() {

	runtime.GOMAXPROCS( 1 )
	
	genAlgo := ga.NewGeneticAlgorithm()

	genAlgo.Simulator = &StringMaterSimulator{}
	genAlgo.BitsetCreate = &MyBitsetCreate{}
	genAlgo.EliteConsumer = &MyEliteConsumer{}
	genAlgo.Mater = ga.NewMater( 
		[]ga.MaterFunctionProbability{
			{ P : 1.0, F : ga.TwoPointCrossover },
			{ P : 1.0, F : ga.Mutate },
			{ P : 1.0, F : ga.UniformCrossover, UseElite : true },
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