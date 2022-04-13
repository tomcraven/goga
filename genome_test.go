package goga_test

import (
	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
)

type GenomeSuite struct {
	genome goga.Genome
}

func (s *GenomeSuite) SetUpTest(t *C) {
	s.genome = goga.NewGenome(goga.Bitset{})
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
	b := goga.Bitset{}
	b.Create(10)
	b.Set(1, 1)
	b.Set(9, 1)

	g := goga.NewGenome(b)
	t.Assert(&b, DeepEquals, g.GetBits())
}
