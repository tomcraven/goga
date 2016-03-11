package ga_test

import (
	. "gopkg.in/check.v1"

	"github.com/tomcraven/goga"
	// "fmt"
	"math/rand"
	"sync"
	"time"
)

const kNumThreads = 4

type GeneticAlgorithmSuite struct {
}

var _ = Suite(&GeneticAlgorithmSuite{})

func helperGenerateExitFunction(numIterations int) func(*ga.IGenome) bool {
	totalIterations := 0
	return func(*ga.IGenome) bool {
		totalIterations++
		if totalIterations >= numIterations {
			return true
		}
		return false
	}
}

func (s *GeneticAlgorithmSuite) TestShouldSimulateUntil(t *C) {

	callCount := 0
	exitFunc := func(g *ga.IGenome) bool {
		callCount++
		return true
	}

	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Init(1, kNumThreads)
	ret := genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 1)

	callCount = 0
	exitFunc2 := func(g *ga.IGenome) bool {
		callCount++
		if callCount >= 2 {
			return true
		}
		return false
	}
	ret = genAlgo.SimulateUntil(exitFunc2)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 2)
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_1(t *C) {

	numCalls1 := 0
	mateFunc1 := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		numCalls1++
		return *a, *b
	}

	numCalls2 := 0
	mateFunc2 := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		numCalls2++
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 0.5, F: mateFunc1},
			{P: 0.75, F: mateFunc2},
		},
	)

	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Init(2, kNumThreads)
	genAlgo.Mater = m

	numIterations := 1000
	ret := genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ret, IsTrue)

	sixtyPercent := (numIterations / 100) * 60
	fourtyPercent := (numIterations / 100) * 40
	t.Assert(numCalls1 < sixtyPercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls1, sixtyPercent))
	t.Assert(numCalls1 > fourtyPercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls1, fourtyPercent))

	sixtyFivePercent := (numIterations / 100) * 65
	eightyFivePercent := (numIterations / 100) * 85
	t.Assert(numCalls2 < eightyFivePercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls2, sixtyPercent))
	t.Assert(numCalls2 > sixtyFivePercent, IsTrue, Commentf("Num calls [%v] percent [%v]", numCalls2, fourtyPercent))
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_2(t *C) {

	numCalls := 0
	mateFunc := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		numCalls++
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo := ga.NewGeneticAlgorithm()
	populationSize := 100
	genAlgo.Init(populationSize, kNumThreads)
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	expectedNumIterations := (numIterations * (populationSize / 2))
	expectedNumIterations -= (populationSize / 2)
	t.Assert(numCalls, Equals, expectedNumIterations)
}

func (s *GeneticAlgorithmSuite) TestShouldCallMaterAppropriately_OddSizedPopulation(t *C) {

	numCalls := 0
	mateFunc := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		numCalls++
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo := ga.NewGeneticAlgorithm()
	populationSize := 99
	genAlgo.Init(populationSize, kNumThreads)
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	resultantPopulationSize := len(genAlgo.GetPopulation())
	t.Assert(resultantPopulationSize, Equals, populationSize)
}

type MyEliteConsumerCounter struct {
	NumCalls int
}

func (ec *MyEliteConsumerCounter) OnElite(g *ga.IGenome) {
	ec.NumCalls++
}

func (s *GeneticAlgorithmSuite) TestShouldCallIntoEliteConsumer(t *C) {

	ec := MyEliteConsumerCounter{}
	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Init(1, kNumThreads)
	genAlgo.EliteConsumer = &ec

	numIterations := 42
	ret := genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ret, IsTrue)
	t.Assert(ec.NumCalls, Equals, numIterations)
}

func (s *GeneticAlgorithmSuite) TestShouldNotSimulateWithNoPopulation(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	callCount := 0
	exitFunc := func(g *ga.IGenome) bool {
		callCount++
		return true
	}
	ret := genAlgo.SimulateUntil(exitFunc)

	t.Assert(ret, IsFalse)
	t.Assert(callCount, Equals, 0)

	genAlgo.Init(0, kNumThreads)
	ret = genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsFalse)
	t.Assert(callCount, Equals, 0)

	genAlgo.Init(1, kNumThreads)
	ret = genAlgo.SimulateUntil(exitFunc)
	t.Assert(ret, IsTrue)
	t.Assert(callCount, Equals, 1)
}

func (s *GeneticAlgorithmSuite) TestShouldGetPopulation(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	t.Assert(genAlgo.GetPopulation(), HasLen, 0)

	genAlgo.Init(1, kNumThreads)
	pop := genAlgo.GetPopulation()
	t.Assert(pop, HasLen, 1)

	g := ga.NewGenome(ga.Bitset{})
	t.Assert(pop[0], FitsTypeOf, g)

	genAlgo.Init(123, kNumThreads)
	t.Assert(genAlgo.GetPopulation(), HasLen, 123)

	p1 := genAlgo.GetPopulation()
	p2 := genAlgo.GetPopulation()
	t.Assert(len(p1), Equals, len(p2))
	for i := 0; i < len(p1); i++ {
		t.Assert(p1[i], Equals, p2[i])
	}
}

type MySimulatorCounter struct {
	NumCalls int
	m        sync.Mutex
}

func (ms *MySimulatorCounter) Simulate(*ga.IGenome) {
	ms.m.Lock()
	ms.NumCalls++
	ms.m.Unlock()
}
func (ms *MySimulatorCounter) OnBeginSimulation() {
}
func (ms *MySimulatorCounter) OnEndSimulation() {
}
func (ms *MySimulatorCounter) ExitFunc(*ga.IGenome) bool {
	return false
}

func (s *GeneticAlgorithmSuite) TestShouldSimulatePopulatonCounter(t *C) {
	genAlgo := ga.NewGeneticAlgorithm()

	ms := MySimulatorCounter{}
	genAlgo.Simulator = &ms
	t.Assert(ms.NumCalls, Equals, 0)

	populationSize := 100
	genAlgo.Init(populationSize, kNumThreads)

	numIterations := 10
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(ms.NumCalls, Equals, numIterations*populationSize)
}

type MySimulatorFitness struct {
	NumIterations     int
	LargestFitnessess []int

	currentLargestFitness int
	m                     sync.Mutex
}

func (ms *MySimulatorFitness) Simulate(g *ga.IGenome) {
	ms.m.Lock()
	randomFitness := rand.Intn(1000)
	if randomFitness > ms.currentLargestFitness {
		ms.currentLargestFitness = randomFitness
	}
	(*g).SetFitness(randomFitness)
	ms.m.Unlock()
}
func (ms *MySimulatorFitness) OnBeginSimulation() {
	ms.currentLargestFitness = 0
}
func (ms *MySimulatorFitness) OnEndSimulation() {
	ms.LargestFitnessess = append(ms.LargestFitnessess, ms.currentLargestFitness)
}
func (ms *MySimulatorFitness) ExitFunc(*ga.IGenome) bool {
	return false
}

type MyEliteConsumerFitness struct {
	EliteFitnesses []int
}

func (ec *MyEliteConsumerFitness) OnElite(g *ga.IGenome) {
	ec.EliteFitnesses = append(ec.EliteFitnesses, (*g).GetFitness())
}

func (s *GeneticAlgorithmSuite) TestShouldSimulatePopulationAndPassEliteToConsumer(t *C) {
	genAlgo := ga.NewGeneticAlgorithm()

	numIterations := 100
	ms := MySimulatorFitness{NumIterations: numIterations}
	genAlgo.Simulator = &ms

	ec := MyEliteConsumerFitness{}
	genAlgo.EliteConsumer = &ec

	genAlgo.Init(100, kNumThreads)

	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))

	t.Assert(ec.EliteFitnesses, DeepEquals, ms.LargestFitnessess)
}

type MySimulatorOrder struct {
	Order []int

	BeginCalled    bool
	SimulateCalled bool
	EndCalled      bool
}

func (ms *MySimulatorOrder) OnBeginSimulation() {
	ms.Order = append(ms.Order, 1)
	ms.SimulateCalled = true
}
func (ms *MySimulatorOrder) Simulate(g *ga.IGenome) {
	ms.Order = append(ms.Order, 2)
	ms.BeginCalled = true
}
func (ms *MySimulatorOrder) OnEndSimulation() {
	ms.Order = append(ms.Order, 3)
	ms.EndCalled = true
}
func (ms *MySimulatorOrder) ExitFunc(*ga.IGenome) bool {
	return false
}

func (s *GeneticAlgorithmSuite) TestShouldCallOnBeginEndSimulation(t *C) {
	genAlgo := ga.NewGeneticAlgorithm()

	ms := MySimulatorOrder{}
	genAlgo.Simulator = &ms

	t.Assert(ms.BeginCalled, Equals, false)
	t.Assert(ms.SimulateCalled, Equals, false)
	t.Assert(ms.Order, HasLen, 0)

	genAlgo.Init(1, kNumThreads)
	genAlgo.SimulateUntil(helperGenerateExitFunction(1))

	// Sleep and give time for threads to start up
	time.Sleep(100 * time.Millisecond)

	t.Assert(ms.BeginCalled, Equals, true)
	t.Assert(ms.SimulateCalled, Equals, true)
	t.Assert(ms.Order, HasLen, 3)
	t.Assert(ms.Order, DeepEquals, []int{1, 2, 3})
}

func (s *GeneticAlgorithmSuite) TestShouldPassEliteToExitFunc(t *C) {
	genAlgo := ga.NewGeneticAlgorithm()

	numIterations := 10
	ms := MySimulatorFitness{NumIterations: numIterations}
	genAlgo.Simulator = &ms

	ec := MyEliteConsumerFitness{}
	genAlgo.EliteConsumer = &ec

	populationSize := 10
	genAlgo.Init(populationSize, kNumThreads)

	passedGenomeFitnesses := make([]int, populationSize)
	callCount := 0
	exitFunc := func(g *ga.IGenome) bool {
		passedGenomeFitnesses[callCount] = (*g).GetFitness()

		callCount++
		if callCount >= numIterations {
			return true
		}
		return false
	}

	genAlgo.SimulateUntil(exitFunc)

	t.Assert(passedGenomeFitnesses, DeepEquals, ms.LargestFitnessess)
}

func (s *GeneticAlgorithmSuite) TestShouldNotCallMaterWithGenomesFromPopulation(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	mateFunc := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		population := genAlgo.GetPopulation()
		aFound, bFound := false, false
		for i := range population {
			if a == &population[i] {
				aFound = true
				if aFound && bFound {
					break
				}
			} else if b == &population[i] {
				bFound = true
				if aFound && bFound {
					break
				}
			}
		}
		t.Assert(aFound, IsFalse)
		t.Assert(bFound, IsFalse)
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1, F: mateFunc},
		},
	)

	genAlgo.Init(10, kNumThreads)
	genAlgo.Mater = m

	numIterations := 1000
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
}

type MySelectorCounter struct {
	CallCount int
}

func (ms *MySelectorCounter) Go(genomes []ga.IGenome, totalFitness int) *ga.IGenome {
	ms.CallCount++
	return &genomes[0]
}

func (s *GeneticAlgorithmSuite) TestShouldCallSelectorAppropriately(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	selector := MySelectorCounter{}
	genAlgo.Selector = &selector

	populationSize := 100
	genAlgo.Init(populationSize, kNumThreads)
	t.Assert(selector.CallCount, Equals, 0)

	numIterations := 100
	genAlgo.SimulateUntil(helperGenerateExitFunction(numIterations))
	t.Assert(selector.CallCount, Equals, (populationSize*numIterations)-populationSize)
}

type MySelectorPassCache struct {
	PassedGenomes []*ga.IGenome
}

func (ms *MySelectorPassCache) Go(genomes []ga.IGenome, totalFitness int) *ga.IGenome {
	randomGenome := &genomes[rand.Intn(len(genomes))]
	ms.PassedGenomes = append(ms.PassedGenomes, randomGenome)
	return randomGenome
}

type MyMaterPassCache struct {
	PassedGenomes []*ga.IGenome
}

func (ms *MyMaterPassCache) Go(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
	ms.PassedGenomes = append(ms.PassedGenomes, a)
	ms.PassedGenomes = append(ms.PassedGenomes, b)
	return *a, *b
}
func (ms *MyMaterPassCache) OnElite(*ga.IGenome) {
}

func (s *GeneticAlgorithmSuite) TestShouldPassSelectedGenomesToMater(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	selector := MySelectorPassCache{}
	genAlgo.Selector = &selector

	mater := MyMaterPassCache{}
	genAlgo.Mater = &mater
	genAlgo.Simulator = &MySimulatorFitness{}

	genAlgo.Init(100, kNumThreads)
	genAlgo.SimulateUntil(helperGenerateExitFunction(100))

	t.Assert(len(mater.PassedGenomes), Equals, len(selector.PassedGenomes))
	t.Assert(mater.PassedGenomes, DeepEquals, selector.PassedGenomes)
}

type MyBitsetCreateCounter struct {
	NumCalls int
}

func (gc *MyBitsetCreateCounter) Go() ga.Bitset {
	gc.NumCalls++
	return ga.Bitset{}
}

func (s *GeneticAlgorithmSuite) TestShouldCallIntoBitsetCreate(t *C) {

	genAlgo := ga.NewGeneticAlgorithm()

	bitsetCreate := MyBitsetCreateCounter{}
	genAlgo.BitsetCreate = &bitsetCreate

	numGenomes := 100
	genAlgo.Init(numGenomes, kNumThreads)

	t.Assert(bitsetCreate.NumCalls, Equals, numGenomes)
}

type MyMaterPassCache2 struct {
	PassedGenomes  []ga.IGenome
	runningFitness int
}

func (ms *MyMaterPassCache2) Go(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {

	g1, g2 := ga.NewGenome(ga.Bitset{}), ga.NewGenome(ga.Bitset{})

	ms.PassedGenomes = append(ms.PassedGenomes, g1)
	ms.PassedGenomes = append(ms.PassedGenomes, g2)

	g1.SetFitness(ms.runningFitness)
	ms.runningFitness++
	g2.SetFitness(ms.runningFitness)
	ms.runningFitness++

	return g1, g2
}
func (ms *MyMaterPassCache2) OnElite(*ga.IGenome) {
}

func (s *GeneticAlgorithmSuite) TestShouldReplaceOldPopulationWithMatedOne(t *C) {

	mater := MyMaterPassCache2{}

	genAlgo := ga.NewGeneticAlgorithm()
	genAlgo.Mater = &mater
	populationSize := 10
	genAlgo.Init(populationSize, kNumThreads)
	genAlgo.SimulateUntil(helperGenerateExitFunction(2))

	genAlgoPopulation := genAlgo.GetPopulation()
	t.Assert(mater.PassedGenomes, HasLen, populationSize)
	t.Assert(genAlgoPopulation, HasLen, populationSize)

	for i := 0; i < populationSize; i++ {
		t.Assert(mater.PassedGenomes[i].GetFitness(), Equals, genAlgoPopulation[i].GetFitness())
		t.Assert(mater.PassedGenomes[i], Equals, genAlgoPopulation[i])
	}
}

type MySimulatorCallTracker struct {
	NumBeginSimulationsUntilExit int

	NumBeginSimulationCalls int
	NumSimulateCalls        int
	m                       sync.Mutex
}

func (ms *MySimulatorCallTracker) Simulate(*ga.IGenome) {
	ms.m.Lock()
	ms.NumSimulateCalls++
	ms.m.Unlock()
}
func (ms *MySimulatorCallTracker) OnBeginSimulation() {
	ms.NumBeginSimulationCalls++
}
func (ms *MySimulatorCallTracker) OnEndSimulation() {
}
func (ms *MySimulatorCallTracker) ExitFunc(*ga.IGenome) bool {
	return (ms.NumBeginSimulationCalls >= ms.NumBeginSimulationsUntilExit)
}

func (s *GeneticAlgorithmSuite) TestShouldSimulateUsingSimulatorExitFunction(t *C) {
	genAlgo := ga.NewGeneticAlgorithm()

	ms := MySimulatorCallTracker{}
	ms.NumBeginSimulationsUntilExit = 5
	genAlgo.Simulator = &ms
	t.Assert(ms.NumBeginSimulationCalls, Equals, 0)
	t.Assert(ms.NumSimulateCalls, Equals, 0)

	populationSize := 100
	genAlgo.Init(populationSize, kNumThreads)
	genAlgo.Simulate()

	t.Assert(ms.NumSimulateCalls, Equals, ms.NumBeginSimulationsUntilExit*populationSize)
	t.Assert(ms.NumBeginSimulationCalls, Equals, ms.NumBeginSimulationsUntilExit)
}
