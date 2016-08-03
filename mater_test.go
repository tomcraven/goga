package ga_test

import (
	. "gopkg.in/check.v1"

	"github.com/tomcraven/bitset"
	ga "github.com/tomcraven/goga"
)

type MaterSuite struct {
	mater ga.Mater
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
	g1 := ga.NewGenome(bitset.Create(0))
	g2 := ga.NewGenome(bitset.Create(0))
	s.mater.Go(&g1, &g2)
}

func (s *MaterSuite) TestGoShouldReturn2NewGenomes(t *C) {

	b1 := bitset.Create(10)
	b2 := bitset.Create(10)
	b1.ClearAll()
	b2.SetAll()

	g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
	t.Assert(g1, Not(DeepEquals), g2)

	c1, c2 := s.mater.Go(&g1, &g2)

	var iGenome ga.Genome
	t.Assert(c1, Implements, &iGenome)
	t.Assert(c2, Implements, &iGenome)

	t.Assert(g1, DeepEquals, c1)
	t.Assert(g2, DeepEquals, c2)
}

func (s *MaterSuite) TestShouldOnePointCrossover_DifferentBitset(t *C) {

	for i := 0; i < 100; i++ {
		genomeSize := uint(10)
		b1 := bitset.Create(genomeSize)
		b2 := bitset.Create(genomeSize)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		var iGenome ga.Genome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		t.Assert(g1, Not(DeepEquals), c1)
		t.Assert(g1, Not(DeepEquals), c2)
		t.Assert(g2, Not(DeepEquals), c1)
		t.Assert(g2, Not(DeepEquals), c2)

		crossoverPoints := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		startingBitValue := c1Bits.Get(0)
		for i := uint(0); i < genomeSize; i++ {
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
		genomeSize := uint(10)
		b1 := bitset.Create(genomeSize)
		b2 := bitset.Create(genomeSize)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		var iGenome ga.Genome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		t.Assert(g1, Not(DeepEquals), c1)
		t.Assert(g1, Not(DeepEquals), c2)
		t.Assert(g2, Not(DeepEquals), c1)
		t.Assert(g2, Not(DeepEquals), c2)

		crossoverPoints := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		startingBitValue := c1Bits.Get(0)
		for i := uint(0); i < genomeSize; i++ {
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
		genomeSize := uint(1000)
		b1 := bitset.Create(genomeSize)
		b2 := bitset.Create(genomeSize)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		crossoverPoints := 0
		startingBitValue := c1Bits.Get(0)

		for i := uint(0); i < genomeSize; i++ {
			t.Assert(c1Bits.Get(i), Not(Equals), c2Bits.Get(i), Commentf("Index [%v]", i))

			// Find the crossover point
			if startingBitValue != c1Bits.Get(i) {
				crossoverPoints++
				startingBitValue = c1Bits.Get(i)
			}
		}

		sixtyPercent := int((genomeSize / 100) * 60)
		fourtyPercent := int((genomeSize / 100) * 40)
		t.Assert(crossoverPoints > fourtyPercent, IsTrue)
		t.Assert(crossoverPoints < sixtyPercent, IsTrue)
	}
}

func (s *MaterSuite) TestShouldOnePointCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 100; i++ {
		genomeSize1 := uint(10)
		genomeSize2 := uint(5)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.Size() == uint(10)) || (c2Bits.Size() == uint(10))
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.Size() == uint(5)) || (c2Bits.Size() == uint(5))
		t.Assert(oneHasSize5, IsTrue)
	}

	for i := 0; i < 100; i++ {
		genomeSize1 := uint(5)
		genomeSize2 := uint(10)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.OnePointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.Size() == uint(10)) || (c2Bits.Size() == uint(10))
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.Size() == uint(5)) || (c2Bits.Size() == uint(5))
		t.Assert(oneHasSize5, IsTrue)
	}
}

func (s *MaterSuite) TestShouldTwoPointCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 100; i++ {
		genomeSize1 := uint(10)
		genomeSize2 := uint(5)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.Size() == uint(10)) || (c2Bits.Size() == uint(10))
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.Size() == uint(5)) || (c2Bits.Size() == uint(5))
		t.Assert(oneHasSize5, IsTrue)
	}

	for i := 0; i < 100; i++ {
		genomeSize1 := uint(5)
		genomeSize2 := uint(10)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.TwoPointCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneHasSize10 := (c1Bits.Size() == uint(10)) || (c2Bits.Size() == uint(10))
		t.Assert(oneHasSize10, IsTrue)
		oneHasSize5 := (c1Bits.Size() == uint(5)) || (c2Bits.Size() == uint(5))
		t.Assert(oneHasSize5, IsTrue)
	}
}

func (s *MaterSuite) TestShouldUniformCrossover_DifferentSizedBitsets(t *C) {

	for i := 0; i < 10; i++ {
		genomeSize1 := uint(20)
		genomeSize2 := uint(10)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneSizedGenomeSize1 := (c1Bits.Size() == genomeSize1) || (c2Bits.Size() == genomeSize1)
		t.Assert(oneSizedGenomeSize1, IsTrue)
		oneSizedGenomeSize2 := (c1Bits.Size() == genomeSize2) || (c2Bits.Size() == genomeSize2)
		t.Assert(oneSizedGenomeSize2, IsTrue)
	}

	for i := 0; i < 10; i++ {
		genomeSize1 := uint(10)
		genomeSize2 := uint(20)

		b1 := bitset.Create(genomeSize1)
		b2 := bitset.Create(genomeSize2)

		b1.ClearAll()
		b2.SetAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b2)
		c1, c2 := ga.UniformCrossover(&g1, &g2)

		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		oneSizedGenomeSize1 := (c1Bits.Size() == genomeSize1) || (c2Bits.Size() == genomeSize1)
		t.Assert(oneSizedGenomeSize1, IsTrue)
		oneSizedGenomeSize2 := (c1Bits.Size() == genomeSize2) || (c2Bits.Size() == genomeSize2)
		t.Assert(oneSizedGenomeSize2, IsTrue)
	}
}

func (s *MaterSuite) TestShouldConfig_Single(t *C) {

	numCalls := 0
	myFunc := func(a, b *ga.Genome) (ga.Genome, ga.Genome) {
		numCalls++
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: myFunc},
		},
	)

	numIterations := 100
	b1 := bitset.Create(10)
	b2 := bitset.Create(10)
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
		myFunc1 := func(a, b *ga.Genome) (ga.Genome, ga.Genome) {
			numCalls1++
			return *a, *b
		}
		myFunc2 := func(a, b *ga.Genome) (ga.Genome, ga.Genome) {
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
		b1 := bitset.Create(10)
		b2 := bitset.Create(10)
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

	genomeSize := uint(10)
	for i := 0; i < 100; i++ {
		b1 := bitset.Create(genomeSize)
		b1.ClearAll()

		g1, g2 := ga.NewGenome(b1), ga.NewGenome(b1)
		c1, c2 := ga.Mutate(&g1, &g2)

		var iGenome ga.Genome
		t.Assert(c1, Implements, &iGenome)
		t.Assert(c2, Implements, &iGenome)

		differringPointsC1 := 0
		differringPointsC2 := 0
		c1Bits, c2Bits := c1.GetBits(), c2.GetBits()
		for i := uint(0); i < genomeSize; i++ {
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

	elite := ga.NewGenome(bitset.Create(0))
	myFunc := func(a, b *ga.Genome) (ga.Genome, ga.Genome) {
		t.Assert(b, Equals, &elite)
		return *a, *b
	}

	m := ga.NewMater(
		[]ga.MaterFunctionProbability{
			{P: 1.0, F: myFunc, UseElite: true},
		},
	)

	m.OnElite(&elite)

	g1 := ga.NewGenome(bitset.Create(0))
	g2 := ga.NewGenome(bitset.Create(0))
	m.Go(&g1, &g2)
}
