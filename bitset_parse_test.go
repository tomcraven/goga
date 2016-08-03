package ga_test

import (
	"math/rand"

	"github.com/tomcraven/bitset"
	ga "github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
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
	s.bp.SetFormat([]uint{1, 2, 3, 4, 5})
}

func (s *BitsetParseSuite) TestShouldPanicWithMismatchedFormatAndBitsetSize(t *C) {
	inputFormat := []uint{
		1, 5, 3, 9, 10, 1,
	}
	s.bp.SetFormat(inputFormat)

	inputBitset := bitset.Create(1)

	t.Assert(func() { s.bp.Process(inputBitset) }, Panics, "Input format does not match bitset size")
}

func (s *BitsetParseSuite) TestShouldNotPanicWithCorrectFormatAndBitsetSize(t *C) {
	inputFormat := []uint{
		uint(rand.Intn(10)),
		uint(rand.Intn(10)),
		uint(rand.Intn(10)),
		uint(rand.Intn(10)),
		uint(rand.Intn(10)),
		uint(rand.Intn(10)),
	}
	s.bp.SetFormat(inputFormat)

	bitsetSize := uint(0)
	for _, i := range inputFormat {
		bitsetSize += i
	}

	inputBitset := bitset.Create(bitsetSize)
	s.bp.Process(inputBitset)
}

func (s *BitsetParseSuite) TestShouldProcessSingleFormat(t *C) {
	inputFormat := []uint{
		16,
	}

	s.bp.SetFormat(inputFormat)

	inputBitset := bitset.Create(16)
	for i := uint(0); i < 16; i++ {
		inputBitset.Set(i)
	}
	t.Assert(s.bp.Process(inputBitset), DeepEquals, []uint64{65535})

	for i := uint(0); i < 16; i++ {
		inputBitset.Clear(i)
	}
	t.Assert(s.bp.Process(inputBitset), DeepEquals, []uint64{0})
}

func (s *BitsetParseSuite) TestShouldProcessMultipleFormat(t *C) {
	inputFormat := []uint{
		8, 8,
	}

	s.bp.SetFormat(inputFormat)

	inputBitset := bitset.Create(16)
	for i := uint(0); i < 8; i++ {
		inputBitset.Set(i)
	}
	for i := uint(8); i < 16; i++ {
		inputBitset.Clear(i)
	}
	t.Assert(s.bp.Process(inputBitset), DeepEquals, []uint64{255, 0})

	for i := uint(0); i < 8; i++ {
		inputBitset.Clear(i)
	}
	for i := uint(8); i < 16; i++ {
		inputBitset.Set(i)
	}
	t.Assert(s.bp.Process(inputBitset), DeepEquals, []uint64{0, 255})
}
