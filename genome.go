package ga

import "github.com/tomcraven/bitset"

// Genome associates a fitness with a bitset
type Genome interface {
	GetFitness() int
	SetFitness(int)
	GetBits() bitset.Bitset
}

type genome struct {
	fitness int
	bitset  bitset.Bitset
}

// NewGenome creates a genome with a bitset and
// a zero'd fitness score
func NewGenome(bitset bitset.Bitset) Genome {
	return &genome{bitset: bitset}
}

func (g *genome) GetFitness() int {
	return g.fitness
}

func (g *genome) SetFitness(fitness int) {
	g.fitness = fitness
}

func (g *genome) GetBits() bitset.Bitset {
	return g.bitset
}
