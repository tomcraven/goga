package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/tomcraven/goga"
)

type stringMaterSimulator struct {
}

func (sms *stringMaterSimulator) OnBeginSimulation() {
}
func (sms *stringMaterSimulator) OnEndSimulation() {
}
func (sms *stringMaterSimulator) Simulate(g goga.Genome) {
	bits := g.GetBits()
	for i, character := range targetString {
		for j := 0; j < 8; j++ {
			targetBit := character & (1 << uint(j))
			bit := bits.Get((i * 8) + j)
			if targetBit != 0 && bit == 1 {
				g.SetFitness(g.GetFitness() + 1)
			} else if targetBit == 0 && bit == 0 {
				g.SetFitness(g.GetFitness() + 1)
			}
		}
	}
}
func (sms *stringMaterSimulator) ExitFunc(g goga.Genome) bool {
	return g.GetFitness() == targetLength
}

type myBitsetCreate struct {
}

func (bc *myBitsetCreate) Go() goga.Bitset {
	b := goga.Bitset{}
	b.Create(targetLength)
	for i := 0; i < targetLength; i++ {
		b.Set(i, rand.Intn(2))
	}
	return b
}

type myEliteConsumer struct {
	currentIter int
}

func (ec *myEliteConsumer) OnElite(g goga.Genome) {
	gBits := g.GetBits()
	ec.currentIter++
	var genomeString string
	for i := 0; i < gBits.GetSize(); i += 8 {
		c := int(0)
		for j := 0; j < 8; j++ {
			bit := gBits.Get(i + j)
			if bit != 0 {
				c |= 1 << uint(j)
			}
		}
		genomeString += string(c)
	}

	fmt.Println(ec.currentIter, "\t", genomeString, "\t", g.GetFitness())
}

const (
	populationSize = 600
)

var (
	targetString = "abcdefghijklmnopqrstuvwxyz"
	targetLength int
)

func init() {
	if len(os.Args) > 1 {
		targetString = os.Args[1]
	}
	targetLength = len(targetString) * 8
}

func main() {

	numThreads := 4
	runtime.GOMAXPROCS(numThreads)

	genAlgo := goga.NewGeneticAlgorithm()

	genAlgo.Simulator = &stringMaterSimulator{}
	genAlgo.BitsetCreate = &myBitsetCreate{}
	genAlgo.EliteConsumer = &myEliteConsumer{}
	genAlgo.Mater = goga.NewMater(
		[]goga.MaterFunctionProbability{
			{P: 1.0, F: goga.TwoPointCrossover},
			{P: 1.0, F: goga.Mutate},
			{P: 1.0, F: goga.UniformCrossover, UseElite: true},
		},
	)
	genAlgo.Selector = goga.NewSelector(
		[]goga.SelectorFunctionProbability{
			{P: 1.0, F: goga.Roulette},
		},
	)

	genAlgo.Init(populationSize, numThreads)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
