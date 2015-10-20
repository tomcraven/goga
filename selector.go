
package ga

import (
	"math/rand"
	// "fmt"
)

type ISelector interface {
	Go( []IGenome, int ) *IGenome
}

type NullSelector struct {
}
func ( ns *NullSelector ) Go( genomes []IGenome, totalFitness int) *IGenome {
	return &genomes[0]
}

type SelectorFunctionProbability struct {
	P float32
	F func( []IGenome, int ) *IGenome
}

type selector struct {
	selectorConfig []SelectorFunctionProbability
}

func NewSelector( selectorConfig []SelectorFunctionProbability ) ISelector {
	return &selector {
		selectorConfig : selectorConfig,
	}
}

func ( s *selector ) Go( genomeArray []IGenome, totalFitness int ) *IGenome {
	for {
		for _, config := range s.selectorConfig {
			if rand.Float32() < config.P {
				return config.F( genomeArray, totalFitness )
			}
		}
	}
	return nil
}

func Roulette( genomeArray []IGenome, totalFitness int ) *IGenome {

	if len( genomeArray ) == 0 {
		panic( "genome array contains no elements" )
	}

	if ( totalFitness == 0 ) {
		randomIndex := rand.Intn( len( genomeArray ) )
		return &genomeArray[randomIndex]
	}

	randomFitness := rand.Intn( totalFitness )
	for i, _ := range genomeArray {
		randomFitness -= genomeArray[i].GetFitness()
		if randomFitness <= 0 {
			return &genomeArray[i]
		}
	}

	panic( "total fitness is too large" )
}