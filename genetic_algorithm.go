package ga

import (
	"sync"
	"time"
)

// GeneticAlgorithm -
// The main component of goga, holds onto the state of the algorithm -
// * Mater - combining evolved genomes
// * EliteConsumer - an optional class that accepts the 'elite' of each population generation
// * Simulator - a simulation component used to score each genome in each generation
// * BitsetCreate - used to create the initial population of genomes
type GeneticAlgorithm struct {
	Mater         IMater
	EliteConsumer IEliteConsumer
	Simulator     ISimulator
	Selector      ISelector
	BitsetCreate  IBitsetCreate

	populationSize          int
	population              []IGenome
	totalFitness            int
	genomeSimulationChannel chan *IGenome
	exitFunc                func(*IGenome) bool
	waitGroup               *sync.WaitGroup
	parallelSimulations     int
}

// NewGeneticAlgorithm returns a new GeneticAlgorithm structure with null implementations of
// EliteConsumer, Mater, Simulator, Selector and BitsetCreate
func NewGeneticAlgorithm() GeneticAlgorithm {
	return GeneticAlgorithm{
		EliteConsumer: &NullEliteConsumer{},
		Mater:         &NullMater{},
		Simulator:     &NullSimulator{},
		Selector:      &NullSelector{},
		BitsetCreate:  &NullBitsetCreate{},
	}
}

func (ga *GeneticAlgorithm) createPopulation() []IGenome {
	ret := make([]IGenome, ga.populationSize)
	for i := 0; i < ga.populationSize; i++ {
		ret[i] = NewGenome(ga.BitsetCreate.Go())
	}
	return ret
}

// Init initialises internal components, sets up the population size
// and number of parallel simulations
func (ga *GeneticAlgorithm) Init(populationSize, parallelSimulations int) {
	ga.populationSize = populationSize
	ga.population = ga.createPopulation()
	ga.parallelSimulations = parallelSimulations

	ga.waitGroup = new(sync.WaitGroup)
}

func (ga *GeneticAlgorithm) beginSimulation() {
	ga.Simulator.OnBeginSimulation()
	ga.EliteConsumer.OnBeginSimulation()
	ga.totalFitness = 0

	ga.genomeSimulationChannel = make(chan *IGenome)

	// todo: make configurable
	for i := 0; i < ga.parallelSimulations; i++ {
		go func(genomeSimulationChannel chan *IGenome,
			waitGroup *sync.WaitGroup, simulator ISimulator) {

			for genome := range genomeSimulationChannel {
				defer waitGroup.Done()
				simulator.Simulate(genome)
			}
		}(ga.genomeSimulationChannel, ga.waitGroup, ga.Simulator)
	}

	ga.waitGroup.Add(ga.populationSize)
}

func (ga *GeneticAlgorithm) onNewGenomeToSimulate(g *IGenome) {
	ga.genomeSimulationChannel <- g
}

func (ga *GeneticAlgorithm) syncSimulatingGenomes() {
	close(ga.genomeSimulationChannel)
	ga.waitGroup.Wait()
}

func (ga *GeneticAlgorithm) getElite() *IGenome {
	var ret *IGenome
	for i := 0; i < ga.populationSize; i++ {
		if ret == nil || ga.population[i].GetFitness() > (*ret).GetFitness() {
			ret = &ga.population[i]
		}
	}
	return ret
}

// SimulateUntil simulates a population until 'exitFunc' returns true
// The 'exitFunc' is passed the elite of each population and should return true
// if the elite reaches a certain criteria (e.g. fitness above a certain threshold)
func (ga *GeneticAlgorithm) SimulateUntil(exitFunc func(*IGenome) bool) bool {
	ga.exitFunc = exitFunc
	return ga.Simulate()
}

func (ga *GeneticAlgorithm) shouldExit(elite *IGenome) bool {
	if ga.exitFunc == nil {
		return ga.Simulator.ExitFunc(elite)
	}
	return ga.exitFunc(elite)
}

// Simulate runs the genetic algorithm
func (ga *GeneticAlgorithm) Simulate() bool {

	if ga.populationSize == 0 {
		return false
	}

	ga.beginSimulation()
	for i := 0; i < ga.populationSize; i++ {
		ga.onNewGenomeToSimulate(&ga.population[i])
	}
	ga.syncSimulatingGenomes()
	ga.Simulator.OnEndSimulation()

	for {
		elite := ga.getElite()
		ga.Mater.OnElite(elite)
		ga.EliteConsumer.OnElite(elite)
		if ga.shouldExit(elite) {
			break
		}

		time.Sleep(1 * time.Microsecond)

		ga.beginSimulation()

		newPopulation := ga.createPopulation()
		for i := 0; i < ga.populationSize; i += 2 {
			g1 := ga.Selector.Go(ga.population, ga.totalFitness)
			g2 := ga.Selector.Go(ga.population, ga.totalFitness)

			g3, g4 := ga.Mater.Go(g1, g2)

			newPopulation[i] = g3
			ga.onNewGenomeToSimulate(&newPopulation[i])

			if (i + 1) < ga.populationSize {
				newPopulation[i+1] = g4
				ga.onNewGenomeToSimulate(&newPopulation[i+1])
			}
		}
		ga.population = newPopulation
		ga.syncSimulatingGenomes()
		ga.Simulator.OnEndSimulation()
	}

	return true
}

// GetPopulation returns the population
func (ga *GeneticAlgorithm) GetPopulation() []IGenome {
	return ga.population
}
