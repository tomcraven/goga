package ga

import (
	// "math"
	"math/rand"
	// "time"
	// "fmt"
	// "reflect"
)

type IMater interface {
	Go(*IGenome, *IGenome) (IGenome, IGenome)
	OnElite(*IGenome)
}

type NullMater struct {
}

func (nm *NullMater) Go(a, b *IGenome) (IGenome, IGenome) {
	return NewGenome(*(*a).GetBits()), NewGenome(*(*b).GetBits())
}
func (nm *NullMater) OnElite(a *IGenome) {
}

type MaterFunctionProbability struct {
	P        float32
	F        func(*IGenome, *IGenome) (IGenome, IGenome)
	UseElite bool
}

type mater struct {
	materConfig []MaterFunctionProbability
	elite       *IGenome
}

func NewMater(materConfig []MaterFunctionProbability) IMater {
	return &mater{
		materConfig: materConfig,
	}
}

func (m *mater) Go(g1, g2 *IGenome) (IGenome, IGenome) {

	newG1 := NewGenome(*(*g1).GetBits())
	newG2 := NewGenome(*(*g2).GetBits())
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

func (m *mater) OnElite(elite *IGenome) {
	m.elite = elite
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func OnePointCrossover(g1, g2 *IGenome) (IGenome, IGenome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	b1, b2 := Bitset{}, Bitset{}
	g1Size := g1Bits.GetSize()
	g2Size := g2Bits.GetSize()
	b1.Create(g1Size)
	b2.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)
	randIndex := rand.Intn(minSize-1) + 1

	for i := 0; i < randIndex; i++ {
		b1.Set(i, g1Bits.Get(i))
		b2.Set(i, g2Bits.Get(i))
	}

	for i := randIndex; i < minSize; i++ {
		b2.Set(i, g1Bits.Get(i))
		b1.Set(i, g2Bits.Get(i))
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.Set(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.Set(i, g2Bits.Get(i))
		}
	}

	return NewGenome(b1), NewGenome(b2)
}

func TwoPointCrossover(g1, g2 *IGenome) (IGenome, IGenome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	b1, b2 := Bitset{}, Bitset{}
	g1Size := g1Bits.GetSize()
	g2Size := g2Bits.GetSize()
	b1.Create(g1Size)
	b2.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)
	randIndex1 := rand.Intn(minSize-1) + 1
	randIndex2 := randIndex1

	for randIndex1 == randIndex2 {
		randIndex2 = rand.Intn(minSize-1) + 1
	}

	// Note: cannot be same value
	if randIndex1 > randIndex2 {
		randIndex1, randIndex2 = randIndex2, randIndex1
	}

	for i := 0; i < randIndex1; i++ {
		b1.Set(i, g1Bits.Get(i))
		b2.Set(i, g2Bits.Get(i))
	}

	for i := randIndex1; i < randIndex2; i++ {
		b2.Set(i, g1Bits.Get(i))
		b1.Set(i, g2Bits.Get(i))
	}

	for i := randIndex2; i < minSize; i++ {
		b1.Set(i, g1Bits.Get(i))
		b2.Set(i, g2Bits.Get(i))
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.Set(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.Set(i, g2Bits.Get(i))
		}
	}
	return NewGenome(b1), NewGenome(b2)
}

func UniformCrossover(g1, g2 *IGenome) (IGenome, IGenome) {

	g1Bits, g2Bits := (*g1).GetBits(), (*g2).GetBits()

	b1, b2 := Bitset{}, Bitset{}
	g1Size := g1Bits.GetSize()
	g2Size := g2Bits.GetSize()
	b1.Create(g1Size)
	b2.Create(g2Size)

	maxSize := max(g1Size, g2Size)
	minSize := min(g1Size, g2Size)

	for i := 0; i < minSize; i++ {
		if rand.Float32() > 0.5 {
			b1.Set(i, g1Bits.Get(i))
			b2.Set(i, g2Bits.Get(i))
		} else {
			b2.Set(i, g1Bits.Get(i))
			b1.Set(i, g2Bits.Get(i))
		}
	}

	if g1Size > g2Size {
		for i := minSize; i < maxSize; i++ {
			b2.Set(i, g1Bits.Get(i))
		}
	} else {
		for i := minSize; i < maxSize; i++ {
			b1.Set(i, g2Bits.Get(i))
		}
	}

	return NewGenome(b1), NewGenome(b2)
}

func Mutate(g1, g2 *IGenome) (IGenome, IGenome) {

	g1BitsOrig := (*g1).GetBits()
	g1Bits := g1BitsOrig.CreateCopy()
	randomBit := rand.Intn(g1Bits.GetSize())
	g1Bits.Set(randomBit, 1-g1Bits.Get(randomBit))

	return NewGenome(g1Bits), NewGenome(*(*g2).GetBits())
}
