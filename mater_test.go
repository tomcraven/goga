package ga_test

import (
	. "gopkg.in/check.v1"

	"github.com/tomcraven/goga"
	// "fmt"
)

type MaterSuite struct {
	mater ga.IMater
}

func (s *MaterSuite) SetUpTest(t *C) {
	s.mater = ga.NewMater([]ga.MaterFunctionProbability{})
}
func (s *MaterSuite) TearDownTest(t *C) {
	s.mater = nil
}

var _ = Suite(&MaterSuite{})

func (s *MaterSuite) TestShouldInstantiate(t *C) {
	// Tested as part of fixture
}

func (s *MaterSuite) TestGoShouldAccept2GenomePointers(t *C) {
	g1, g2 := ga.NewGenome(ga.Bitset{}), ga.NewGenome(ga.Bitset{})
	s.mater.Go(&g1, &g2)
}

func (s *MaterSuite) TestGoShouldReturn2NewGenomes(t *C) {

	b1, b2 := ga.Bitset{}, ga.Bitset{}
	b1.Create(10)
	b2.Create(10)
	b1.SetAll(0)
	b2.SetAll(1)

	g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
	t.Assert(g1, Not(DeepEquals), g2)

	c1, c2 := s.mater.Go(&g1, &g2)

	var iGenome ga.IGenome
	t.Assert(c1, Implements, &iGenome)
	t.Assert(c2, Implements, &iGenome)

	t.Assert(g1, DeepEquals, c1)
	t.Assert(g2, DeepEquals, c2)
}

func (s *MaterSuite) TestShouldOnePointCrossover_DifferentBitset(t *C) {

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize := 10
		b1.Create(genomeSize)
		b2.Create(genomeSize)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		var iGenome ga.IGenome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		t.Assert(g1, Not(DeepEquals), c1)
		t.Assert(g1, Not(DeepEquals), c2)
		t.Assert(g2, Not(DeepEquals), c1)
		t.Assert(g2, Not(DeepEquals), c2)

		crossoverPoints := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		startingBitValue := c1Bits.Get(0)
		for i := 0; i < genomeSize; i++ {
			t.Assert(c1Bits.Get(i), Not(Equals), c2Bits.Get(i), Commentf("Index [%v]", i))

			// Find the crossover point
			if startingBitValue != c1Bits.Get(i) {
				crossoverPoints++
				startingBitValue = c1Bits.Get(i)
			}
		}
		t.Assert(crossoverPoints, Equals, 1)
	}
}

func (s *MaterSuite) TestShouldTwoPointCrossOver_DifferentBitset(t *C) {

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize := 10
		b1.Create(genomeSize)
		b2.Create(genomeSize)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		var iGenome ga.IGenome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		t.Assert(g1, Not(DeepEquals), c1)
		t.Assert(g1, Not(DeepEquals), c2)
		t.Assert(g2, Not(DeepEquals), c1)
		t.Assert(g2, Not(DeepEquals), c2)

		crossoverPoints := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		startingBitValue := c1Bits.Get(0)
		for i := 0; i < genomeSize; i++ {
			t.Assert(c1Bits.Get(i), Not(Equals), c2Bits.Get(i), Commentf("Index [%v]", i))

			// Find the crossover point
			if startingBitValue != c1Bits.Get(i) {
				crossoverPoints++
				startingBitValue = c1Bits.Get(i)
			}
		}
		t.Assert(crossoverPoints, Equals, 2)
	}
}

func (s *MaterSuite) TestShouldUniformCrossover(t *C) {

	for i := 0; i < 10; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize := 1000
		b1.Create(genomeSize)
		b2.Create(genomeSize)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		crossoverPoints := 0
		startingBitValue := c1Bits.Get(0)
		for i := 0; i < genomeSize; i++ {
			t.Assert(c1Bits.Get(i), Not(Equals), c2Bits.Get(i), Commentf("Index [%v]", i))

			// Find the crossover point
			if startingBitValue != c1Bits.Get(i) {
				crossoverPoints++
				startingBitValue = c1Bits.Get(i)
			}
		}

		sixtyPercent := (genomeSize / 100) * 60
		fourtyPercent := (genomeSize / 100) * 40
		t.Assert(crossoverPoints > fourtyPercent, IsTrue)
		t.Assert(crossoverPoints < sixtyPercent, IsTrue)
	}
}

func (s *MaterSuite) TestShouldOnePointCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 10
		genomeSize2 := 5
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.GetSize() == 10) || (c2Bits.GetSize() == 10)
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.GetSize() == 5) || (c2Bits.GetSize() == 5)
		t.Assert(oneHasSize5, IsTrue)
	}

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 5
		genomeSize2 := 10
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.GetSize() == 10) || (c2Bits.GetSize() == 10)
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.GetSize() == 5) || (c2Bits.GetSize() == 5)
		t.Assert(oneHasSize5, IsTrue)
	}
}

func (s *MaterSuite) TestShouldTwoPointCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 10
		genomeSize2 := 5
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()

		oneHasSize10 := (c1Bits.GetSize() == 10) || (c2Bits.GetSize() == 10)
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.GetSize() == 5) || (c2Bits.GetSize() == 5)
		t.Assert(oneHasSize5, IsTrue)
	}

	for i := 0; i < 100; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 5
		genomeSize2 := 10
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()

		oneHasSize10 := (c1Bits.GetSize() == 10) || (c2Bits.GetSize() == 10)
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.GetSize() == 5) || (c2Bits.GetSize() == 5)
		t.Assert(oneHasSize5, IsTrue)
	}
}

func (s *MaterSuite) TestShouldUniformCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 10; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 20
		genomeSize2 := 10
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneSizedGenomeSize1 := (c1Bits.GetSize() == genomeSize1) || (c2Bits.GetSize() == genomeSize1)
		t.Assert(oneSizedGenomeSize1, IsTrue)
		oneSizedGenomeSize2 := (c1Bits.GetSize() == genomeSize2) || (c2Bits.GetSize() == genomeSize2)
		t.Assert(oneSizedGenomeSize2, IsTrue)
	}

	for i := 0; i < 10; i++ {
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		genomeSize1 := 10
		genomeSize2 := 20
		b1.Create(genomeSize1)
		b2.Create(genomeSize2)
		b1.SetAll(0)
		b2.SetAll(1)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneSizedGenomeSize1 := (c1Bits.GetSize() == genomeSize1) || (c2Bits.GetSize() == genomeSize1)
		t.Assert(oneSizedGenomeSize1, IsTrue)
		oneSizedGenomeSize2 := (c1Bits.GetSize() == genomeSize2) || (c2Bits.GetSize() == genomeSize2)
		t.Assert(oneSizedGenomeSize2, IsTrue)
	}
}

func (s *MaterSuite) TestShouldConfig_Single(t *C) {

	numCalls := 0
	myFunc := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		numCalls++
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: myFunc},
		},
	)

	numIterations := 100
	b1, b2 := ga.Bitset{}, ga.Bitset{}
	b1.Create(10)
	b2.Create(10)
	for i := 0; i < numIterations; i++ {
		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		m.Go(&g1, &g2)
	}

	t.Assert(numCalls, Equals, numIterations)
}

func (s *MaterSuite) TestShouldConfig_Multiple(t *C) {

	for i := 0; i < 100; i++ {
		numCalls1 := 0
		numCalls2 := 0
		myFunc1 := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
			numCalls1++
			return *a, *b
		}
		myFunc2 := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
			numCalls2++
			return *a, *b
		}

		m := ga.NewMater(
			[]ga.MaterFunctionProbability{
				{P: 0.5, F: myFunc1},
				{P: 0.5, F: myFunc2},
			},
		)

		numIterations := 1000
		b1, b2 := ga.Bitset{}, ga.Bitset{}
		b1.Create(10)
		b2.Create(10)
		for i := 0; i < numIterations; i++ {
			g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
			m.Go(&g1, &g2)
		}

		sixtyPercent := (numIterations / 100) * 60
		fourtyPercent := (numIterations / 100) * 40
		t.Assert(numCalls1 < sixtyPercent, IsTrue, Commentf("Num calls [%v] sixty percent [%v]", numCalls1, sixtyPercent))
		t.Assert(numCalls2 < sixtyPercent, IsTrue, Commentf("Num calls [%v] sixty percent [%v]", numCalls2, sixtyPercent))
		t.Assert(numCalls1 > fourtyPercent, IsTrue, Commentf("Num calls [%v] fourty percent [%v]", numCalls1, fourtyPercent))
		t.Assert(numCalls2 > fourtyPercent, IsTrue, Commentf("Num calls [%v] fourty percent [%v]", numCalls2, fourtyPercent))
	}
}

func (s *MaterSuite) TestShouldMutate(t *C) {

	genomeSize := 10
	for i := 0; i < 100; i++ {
		b1 := ga.Bitset{}
		b1.Create(genomeSize)
		b1.SetAll(0)

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b1)
		c1, c2 := ga.Mutate(&g1, &g2)

		var iGenome ga.IGenome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		differringPointsC1 := 0
		differringPointsC2 := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		for i := 0; i < genomeSize; i++ {
			if c1Bits.Get(i) != b1.Get(i) {
				differringPointsC1++
			}
			if c2Bits.Get(i) != b1.Get(i) {
				differringPointsC2++
			}
		}
		oneIsDifferent := (differringPointsC1 == 1) || (differringPointsC2 == 1)
		bothDiffernet := (differringPointsC1 == 1) && (differringPointsC2 == 1)
		t.Assert(oneIsDifferent, IsTrue)
		t.Assert(bothDiffernet, IsFalse)
	}
}

func (s *MaterSuite) TestShouldUseEliteFromConfigSettings(t *C) {

	elite := ga.NewGenome(ga.Bitset{})
	myFunc := func(a, b *ga.IGenome) (ga.IGenome, ga.IGenome) {
		t.Assert(b, Equals, &elite)
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: myFunc, UseElite: true},
		},
	)

	m.OnElite(&elite)

	g1, g2 := ga.NewGenome(ga.Bitset{}), ga.NewGenome(ga.Bitset{})
	m.Go(&g1, &g2)
}
