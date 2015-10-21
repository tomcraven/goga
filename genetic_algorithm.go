
package ga

import (
	// "fmt"
	"time"
)

type GeneticAlgorithm struct {	
	Mater IMater
	EliteConsumer IEliteConsumer
	Simulator ISimulator
	Selector ISelector
	BitsetCreate IBitsetCreate

	populationSize int
	population []IGenome
	totalFitness int
	genomeSimulationChannel chan *IGenome
	exitFunc func( *IGenome ) bool
}

func NewGeneticAlgorithm() GeneticAlgorithm {
	return GeneticAlgorithm {
		EliteConsumer : &NullEliteConsumer{},
		Mater : &NullMater{},
		Simulator : &NullSimulator{},
		Selector : &NullSelector{},
		BitsetCreate : &NullBitsetCreate{},
	}
}

func ( ga *GeneticAlgorithm ) createPopulation() []IGenome {
	ret := make( []IGenome, ga.populationSize )
	for i := 0; i < ga.populationSize; i++ {
		ret[i] = NewGenome( ga.BitsetCreate.Go() )
	}
	return ret
}

func ( ga *GeneticAlgorithm ) Init( populationSize int ) {
	ga.populationSize = populationSize
	ga.genomeSimulationChannel = make( chan *IGenome, populationSize )
	ga.population = ga.createPopulation()
}

func ( ga *GeneticAlgorithm ) beginSimulation() {
	ga.Simulator.OnBeginSimulation()
	ga.totalFitness = 0
}

func ( ga *GeneticAlgorithm ) onNewGenomeToSimulate( g *IGenome ) {
	kMultiThreaded := true

	if kMultiThreaded {
		go func( genome *IGenome ) {
			ga.Simulator.Simulate( genome )
			ga.genomeSimulationChannel <- genome
		}( g )
	} else {
		ga.Simulator.Simulate( g )
		ga.genomeSimulationChannel <- g
	}
}

func ( ga *GeneticAlgorithm ) syncSimulatingGenomes() {
	for i := 0; i < ga.populationSize; i++ {
		g := <- ga.genomeSimulationChannel
		ga.totalFitness += (*g).GetFitness()
	}
}

func ( ga *GeneticAlgorithm ) getElite() *IGenome {
	var ret *IGenome = nil
	for i := 0; i < ga.populationSize; i++ {
		if ret == nil || ga.population[i].GetFitness() > (*ret).GetFitness() {
			ret = &ga.population[i]
		}
	}
	return ret
}

func ( ga *GeneticAlgorithm ) SimulateUntil( exitFunc func( *IGenome ) bool ) bool {
	ga.exitFunc = exitFunc
	return ga.Simulate()
}

func ( ga *GeneticAlgorithm ) shouldExit( elite *IGenome ) bool {
	if ga.exitFunc == nil {
		return ga.Simulator.ExitFunc( elite )
	}
	return ga.exitFunc( elite )
}

func ( ga *GeneticAlgorithm ) Simulate() bool {

	if ( ga.populationSize == 0 ) {
		return false
	}

	ga.beginSimulation()
	for i := 0; i < ga.populationSize; i++ {
		ga.onNewGenomeToSimulate( &ga.population[i] )
	}
	ga.syncSimulatingGenomes()
	ga.Simulator.OnEndSimulation()

	for {
		elite := ga.getElite()
		ga.Mater.OnElite( elite )
		ga.EliteConsumer.OnElite( elite )
		if ga.shouldExit( elite ) {
			break
		}

		time.Sleep( 1 * time.Microsecond )

		ga.beginSimulation()

		newPopulation := ga.createPopulation()
		for i := 0; i < ga.populationSize; i += 2 {
			g1 := ga.Selector.Go( ga.population, ga.totalFitness )
			g2 := ga.Selector.Go( ga.population, ga.totalFitness )

			g3, g4 := ga.Mater.Go(g1, g2)

			newPopulation[i] = g3
			ga.onNewGenomeToSimulate( &newPopulation[i] )

			if ( i + 1 ) < ga.populationSize {
				newPopulation[i + 1] = g4
				ga.onNewGenomeToSimulate( &newPopulation[i + 1] )
			}
		}
		ga.population = newPopulation
		ga.syncSimulatingGenomes()
		ga.Simulator.OnEndSimulation()
	}

	return true
}

func ( ga *GeneticAlgorithm ) GetPopulation() ( []IGenome ) {
	return ga.population
}