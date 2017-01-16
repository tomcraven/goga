package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/tomcraven/bitset"
	ga "github.com/tomcraven/goga"
)

type stringMaterSimulator struct {
}

func (sms *stringMaterSimulator) OnBeginSimulation() {
}
func (sms *stringMaterSimulator) OnEndSimulation() {
}
func (sms *stringMaterSimulator) Simulate(g *ga.Genome) {
	sms.simulateScoreCharacters(g)
}
func (sms *stringMaterSimulator) simulateScoreCharacters(g *ga.Genome) {
	bits := (*g).GetBits()
	for i, targetCharacter := range targetString {
		var evolvedCharacter uint8
		beginIndex := uint(i * 8)
		bits.Slice(beginIndex, beginIndex+8).BuildUint8(&evolvedCharacter)

		if targetCharacter == rune(evolvedCharacter) {
			(*g).SetFitness((*g).GetFitness() + 1)
		}
	}
}
func (sms *stringMaterSimulator) simulateScoreBits(g *ga.Genome) {
	bits := (*g).GetBits()
	for i, character := range targetString {
		for j := 0; j < 8; j++ {
			targetBit := character & (1 << uint(j))
			set := bits.Get(uint((i * 8) + j))
			if targetBit != 0 && set {
				(*g).SetFitness((*g).GetFitness() + 1)
			} else if targetBit == 0 && !set {
				(*g).SetFitness((*g).GetFitness() + 1)
			}
		}
	}
}
func (sms *stringMaterSimulator) ExitFunc(g *ga.Genome) bool {
	return (*g).GetFitness() == (len(targetString) * 8)
}

type myBitsetCreate struct {
}

func (bc *myBitsetCreate) Go() bitset.Bitset {
	b := bitset.Create(targetLength)
	for i := uint(0); i < targetLength; i++ {
		b.SetTo(i, rand.Intn(2) == 0)
	}
	return b
}

type myEliteConsumer struct {
	currentIter int
}

func (ec *myEliteConsumer) OnElite(g *ga.Genome) {
	gBits := (*g).GetBits()
	ec.currentIter++
	var genomeString string
	for i := uint(0); i < gBits.Size(); i += 8 {
		c := int(0)
		for j := uint(0); j < 8; j++ {
			bit := gBits.Get(i + j)
			if bit {
				c |= 1 << uint(j)
			}
		}
		genomeString += string(c)
	}

	fmt.Println(ec.currentIter, "\t", genomeString, "\t", (*g).GetFitness())
}

const (
	populationSize = 600
)

var (
	targetString = "abcdefghijklmnopqrstuvwxyz"
	targetLength uint
)

func init() {
	if len(os.Args) > 1 {
		targetString = os.Args[1]
	}
	targetLength = uint(len(targetString) * 8)
}

func main() {

	numThreads := 4
	runtime.GOMAXPROCS(numThreads)

	genAlgo := ga.NewGeneticAlgorithm()

	genAlgo.Simulator = &stringMaterSimulator{}
	genAlgo.BitsetCreate = &myBitsetCreate{}
	genAlgo.EliteConsumer = &myEliteConsumer{}
	genAlgo.Mater = ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: ga.TwoPointCrossover},
			{P: 1.0, F: ga.Mutate},
			{P: 1.0, F: ga.UniformCrossover, UseElite: true},
		},
	)
	genAlgo.Selector = ga.NewSelector(
		[]ga.SelectorFunctionProbability{
			{P: 1.0, F: ga.Roulette},
		},
	)

	genAlgo.Init(populationSize, numThreads)

	startTime := time.Now()
	genAlgo.Simulate()
	fmt.Println(time.Since(startTime))
}
