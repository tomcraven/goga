package ga

import (
	"math/rand"

	"github.com/tomcraven/bitset"
)

// Mater - an interface to a mater object
type Mater interface {
	Go(*Genome, *Genome) (Genome, Genome)
	OnElite(*Genome)
}

// NullMater - null implementation of the Mater interface
type NullMater struct {
}

// Go - null implementation of the Mater go func
func (nm *NullMater) Go(a, b *Genome) (Genome, Genome) {
	return NewGenome((*a).GetBits()), NewGenome((*b).GetBits())
}

// OnElite - null implementation of the Mater OnElite func
func (nm *NullMater) OnElite(a *Genome) {
}

// MaterFunctionProbability -
// An implementation of Mater that has a function and a probability
// where mater function 'F' is called with a probability of 'P'
// where 'P' is a value between 0 and 1
// 0 = never called, 1 = called for every genome
type MaterFunctionProbability struct {
	P        float32
	F        func(*Genome, *Genome) (Genome, Genome)
	UseElite bool
}

type mater struct {
	materConfig []MaterFunctionProbability
	elite       *Genome
}

// NewMater returns an instance of an Mater with several MaterFuncProbabilities
func NewMater(materConfig []MaterFunctionProbability) Mater {
	return &mater{
		materConfig: materConfig,
	}
}

// Go cycles through, and applies, the configures mater functions in the
// MaterFunctionProbability array
func (m *mater) Go(g1, g2 *Genome) (Genome, Genome) {

	newG1 := NewGenome((*g1).GetBits())
	newG2 := NewGenome((*g2).GetBits())
	for _, config := range m.materConfig {
		if rand.Float32() < config.P {
			if config.UseElite {
				newG1, newG2 = config.F(&newG1, m.elite)
			} else {
				newG1, newG2 = config.F(&newG1, &newG2)
			}
		}
	}

	return newG1, newG2
}

// OnElite -
func (m *mater) OnElite(elite *Genome) {
	m.elite = elite
}

func max(a, b uint) uint {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint) uint {
	if a < b {
		return a
	}
	return b
}

// OnePointCrossover -
// Accepts 2 genomes and combines them to create 2 new genomes using one point crossover
// i.e.
// input genomes of:
// 000000 and 111111
// could produce output genomes of:
// 000111 and 111000
func OnePointCrossover(g1, g2 *Genome) (Genome, Genome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	g1Size := g1Bits.Size()
	g2Size := g2Bits.Size()

	b1 := bitset.Create(g1Size)
	b2 := bitset.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)
	randIndex := uint(rand.Intn(int(minSize-1)) + 1)

	for i := uint(0); i < randIndex; i++ {
		b1.SetTo(i, g1Bits.Get(i))
		b2.SetTo(i, g2Bits.Get(i))
	}

	for i := randIndex; i < minSize; i++ {
		b2.SetTo(i, g1Bits.Get(i))
		b1.SetTo(i, g2Bits.Get(i))
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.SetTo(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.SetTo(i, g2Bits.Get(i))
		}
	}

	return NewGenome(b1), NewGenome(b2)
}

// TwoPointCrossover -
// Accepts 2 genomes and combines them to create 2 new genomes using two point crossover
// i.e.
// input genomes of:
// 000000 and 111111
// could produce output genomes of:
// 001100 and 110011
func TwoPointCrossover(g1, g2 *Genome) (Genome, Genome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	g1Size := g1Bits.Size()
	g2Size := g2Bits.Size()

	b1 := bitset.Create(g1Size)
	b2 := bitset.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)
	randIndex1 := uint(rand.Intn(int(minSize-1)) + 1)
	randIndex2 := randIndex1

	for randIndex1 == randIndex2 {
		randIndex2 = uint(rand.Intn(int(minSize-1)) + 1)
	}

	// Note: cannot be same value
	if randIndex1 > randIndex2 {
		randIndex1, randIndex2 = randIndex2, randIndex1
	}

	for i := uint(0); i < randIndex1; i++ {
		b1.SetTo(i, g1Bits.Get(i))
		b2.SetTo(i, g2Bits.Get(i))
	}

	for i := randIndex1; i < randIndex2; i++ {
		b2.SetTo(i, g1Bits.Get(i))
		b1.SetTo(i, g2Bits.Get(i))
	}

	for i := randIndex2; i < minSize; i++ {
		b1.SetTo(i, g1Bits.Get(i))
		b2.SetTo(i, g2Bits.Get(i))
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.SetTo(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.SetTo(i, g2Bits.Get(i))
		}
	}
	return NewGenome(b1), NewGenome(b2)
}

// UniformCrossover -
// Accepts 2 genomes and combines them to create 2 new genomes using uniform crossover
// i.e.
// input genomes of:
// 000000 and 111111
// could produce output genomes of:
// 101010 and 010101
func UniformCrossover(g1, g2 *Genome) (Genome, Genome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	g1Size := g1Bits.Size()
	g2Size := g2Bits.Size()

	b1 := bitset.Create(g1Size)
	b2 := bitset.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)

	for i := uint(0); i < minSize; i++ {
		if rand.Float32() > 0.5 {
			b1.SetTo(i, g1Bits.Get(i))
			b2.SetTo(i, g2Bits.Get(i))
		} else {
			b2.SetTo(i, g1Bits.Get(i))
			b1.SetTo(i, g2Bits.Get(i))
		}
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.SetTo(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.SetTo(i, g2Bits.Get(i))
		}
	}

	return NewGenome(b1), NewGenome(b2)
}

// Mutate -
// Accepts 2 genomes and mutates a single bit in the first to create a new
// very slightly different genome
// i.e.
// input genomes of:
// 000000 and 111111
// could produce output genomes of:
// 001000 and 111111
func Mutate(g1, g2 *Genome) (Genome, Genome) {

	g1BitsOrig := (*g1).GetBits()
	g1Bits := g1BitsOrig.CreateCopy()
	randomBit := uint(rand.Intn(int(g1Bits.Size())))
	g1Bits.SetTo(randomBit, !g1Bits.Get(randomBit))

	return NewGenome(g1Bits), NewGenome((*g2).GetBits())
}
