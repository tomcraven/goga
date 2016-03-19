package ga

// IGenome associates a fitness with a bitset
type IGenome interface {
	GetFitness() int
	SetFitness(int)
	GetBits() *Bitset
}

type genome struct {
	fitness int
	bitset  Bitset
}

// NewGenome creates a genome with a bitset and
// a zero'd fitness score
func NewGenome(bitset Bitset) IGenome {
	return &genome{bitset: bitset}
}

func (g *genome) GetFitness() int {
	return g.fitness
}

func (g *genome) SetFitness(fitness int) {
	g.fitness = fitness
}

func (g *genome) GetBits() *Bitset {
	return &g.bitset
}
