package ga_test

import (
	. "gopkg.in/check.v1"

	"github.com/tomcraven/goga"
)

type GenomeSuite struct {
	genome ga.IGenome
}

func (s *GenomeSuite) SetUpTest(t *C) {
	s.genome = ga.NewGenome(ga.Bitset{})
}
func (s *GenomeSuite) TearDownTest(t *C) {
	s.genome = nil
}

var _ = Suite(&GenomeSuite{})

func (s *GenomeSuite) TestShouldInstantiate(t *C) {
	// Tested as part of fixture setup
}

func (s *GenomeSuite) TestShouldSetGetFitness(t *C) {
	t.Assert(s.genome.GetFitness(), Equals, 0)

	s.genome.SetFitness(100)
	t.Assert(s.genome.GetFitness(), Equals, 100)
}

func (s *GenomeSuite) TestShouldGetBits(t *C) {
	b := ga.Bitset{}
	b.Create(10)
	b.Set(1, 1)
	b.Set(9, 1)

	g := ga.NewGenome(b)
	t.Assert(&b, DeepEquals, g.GetBits())
}
