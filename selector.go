package goga

import (
	"math/rand"
)

// Selector - a selector interface used to pick 2 genomes to mate
type Selector interface {
	Go([]Genome, int) Genome
}

// NullSelector - a null implementation of the Selector interface
type NullSelector struct {
}

// Go - a null implementation of Selector's 'go'
func (ns *NullSelector) Go(genomes []Genome, totalFitness int) Genome {
	return genomes[0]
}

// SelectorFunctionProbability -
// Contains a selector function and a probability
// where selector function 'F' is called with probability 'P'
// where 'P' is a value between 0 and 1
// 0 = never called, 1 = called every time we need a new genome to mate
type SelectorFunctionProbability struct {
	P float32
	F func([]Genome, int) Genome
}

type selector struct {
	selectorConfig []SelectorFunctionProbability
}

// NewSelector returns an instance of an ISelector with several SelectorFunctionProbabiities
func NewSelector(selectorConfig []SelectorFunctionProbability) Selector {
	return &selector{
		selectorConfig: selectorConfig,
	}
}

// Go - cycles through the selector function probabilities until one returns a genome
func (s *selector) Go(genomeArray []Genome, totalFitness int) Genome {
	for {
		for _, config := range s.selectorConfig {
			if rand.Float32() < config.P {
				return config.F(genomeArray, totalFitness)
			}
		}
	}
}

// Roulette is a selection function that selects a genome where genomes that have a higher fitness are more likely to be picked
func Roulette(genomeArray []Genome, totalFitness int) Genome {

	if len(genomeArray) == 0 {
		panic("genome array contains no elements")
	}

	if totalFitness == 0 {
		randomIndex := rand.Intn(len(genomeArray))
		return genomeArray[randomIndex]
	}

	randomFitness := rand.Intn(totalFitness)
	for i := range genomeArray {
		randomFitness -= genomeArray[i].GetFitness()
		if randomFitness <= 0 {
			return genomeArray[i]
		}
	}

	panic("total fitness is too large")
}
