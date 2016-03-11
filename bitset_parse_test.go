package ga_test

import (
	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
	"math/rand"
)

type BitsetParseSuite struct {
	bp ga.IBitsetParse
}

func (s *BitsetParseSuite) SetUpTest(t *C) {
	s.bp = ga.CreateBitsetParse()
}
func (s *BitsetParseSuite) TearDownTest(t *C) {
	s.bp = nil
}

var _ = Suite(&BitsetParseSuite{})

func (s *BitsetParseSuite) TestShouldInstantiate(t *C) {
	// tested as part of suite setup
}

func (s *BitsetParseSuite) TestShouldSetFormat(t *C) {
	s.bp.SetFormat([]int{1, 2, 3, 4, 5})
}

func (s *BitsetParseSuite) TestShouldPanicWithMismatchedFormatAndBitsetSize(t *C) {
	inputFormat := []int{
		1, 5, 3, 9, 10, 1,
	}
	s.bp.SetFormat(inputFormat)

	inputBitset := ga.Bitset{}
	inputBitset.Create(1) // Size should equal sum of all formats

	t.Assert(func() { s.bp.Process(&inputBitset) }, Panics, "Input format does not match bitset size")
}

func (s *BitsetParseSuite) TestShouldNotPanicWithCorrectFormatAndBitsetSize(t *C) {
	inputFormat := []int{
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
	}
	s.bp.SetFormat(inputFormat)

	bitsetSize := 0
	for _, i := range inputFormat {
		bitsetSize += i
	}

	inputBitset := ga.Bitset{}
	inputBitset.Create(bitsetSize) // Size should equal sum of all formats

	s.bp.Process(&inputBitset)
}

func (s *BitsetParseSuite) TestShouldProcessSingleFormat(t *C) {
	inputFormat := []int{
		16,
	}

	s.bp.SetFormat(inputFormat)

	inputBitset := ga.Bitset{}
	inputBitset.Create(16)
	for i := 0; i < 16; i++ {
		inputBitset.Set(i, 1)
	}
	t.Assert(s.bp.Process(&inputBitset), DeepEquals, []uint64{65535})

	for i := 0; i < 16; i++ {
		inputBitset.Set(i, 0)
	}
	t.Assert(s.bp.Process(&inputBitset), DeepEquals, []uint64{0})
}

func (s *BitsetParseSuite) TestShouldProcessMultipleFormat(t *C) {
	inputFormat := []int{
		8, 8,
	}

	s.bp.SetFormat(inputFormat)

	inputBitset := ga.Bitset{}
	inputBitset.Create(16)
	for i := 0; i < 8; i++ {
		inputBitset.Set(i, 1)
	}
	for i := 8; i < 16; i++ {
		inputBitset.Set(i, 0)
	}
	t.Assert(s.bp.Process(&inputBitset), DeepEquals, []uint64{255, 0})

	for i := 0; i < 8; i++ {
		inputBitset.Set(i, 0)
	}
	for i := 8; i < 16; i++ {
		inputBitset.Set(i, 1)
	}
	t.Assert(s.bp.Process(&inputBitset), DeepEquals, []uint64{0, 255})
}
