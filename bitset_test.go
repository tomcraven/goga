package goga_test

import (
	"github.com/tomcraven/goga"
	. "gopkg.in/check.v1"
)

type BitsetSuite struct {
	bitset *goga.Bitset
}

func (s *BitsetSuite) SetUpTest(t *C) {
	s.bitset = &goga.Bitset{}
}
func (s *BitsetSuite) TearDownTest(t *C) {
	s.bitset = nil
}

var _ = Suite(&BitsetSuite{})

func (s *BitsetSuite) TestShouldInstantiate(t *C) {
	// Tested as part of fixture setup
}

func (s *BitsetSuite) TestShouldCreate(t *C) {
	s.bitset.Create(10)
}

func (s *BitsetSuite) TestShouldGetSize(t *C) {
	t.Assert(0, Equals, s.bitset.GetSize())

	s.bitset.Create(10)
	t.Assert(10, Equals, s.bitset.GetSize())

	b := goga.Bitset{}
	b.Create(100)
	t.Assert(100, Equals, b.GetSize())
	t.Assert(b.Set(99, 1), IsTrue)
}

func (s *BitsetSuite) TestShouldSetAndGet(t *C) {
	s.bitset.Create(10)

	index := 0
	value := 1
	t.Assert(s.bitset.Set(index, value), IsTrue)
	t.Assert(value, Equals, s.bitset.Get(index))

	index = 1
	value = 0
	t.Assert(s.bitset.Set(index, value), IsTrue)
	t.Assert(value, Equals, s.bitset.Get(index))
}

func (s *BitsetSuite) TestShouldFailSetAndGetWhenNotCreated(t *C) {
	t.Assert(s.bitset.Set(0, 1), IsFalse)
	t.Assert(s.bitset.Get(0), Equals, -1)

	s.bitset.Create(10)
	t.Assert(s.bitset.Set(10, 1), IsFalse)
}

func (s *BitsetSuite) TestShouldSetAll(t *C) {
	bitsetSize := 10

	s.bitset.Create(bitsetSize)
	for i := 0; i < bitsetSize; i++ {
		t.Assert(s.bitset.Get(i), Equals, 0)
	}

	s.bitset.SetAll(1)
	for i := 0; i < bitsetSize; i++ {
		t.Assert(s.bitset.Get(i), Equals, 1)
	}
}

func (s *BitsetSuite) TestShouldSlice(t *C) {
	s.bitset.Create(10)
	s.bitset.SetAll(1)

	slice := s.bitset.Slice(0, 3)
	t.Assert(slice, FitsTypeOf, goga.Bitset{})
	t.Assert(slice.GetSize(), Equals, 3)

	for i := 0; i < 3; i++ {
		t.Assert(slice.Get(i), Equals, 1)
	}
}

func (s *BitsetSuite) TestShouldGetAll(t *C) {

	const kBitsetSize = 10
	s.bitset.Create(kBitsetSize)
	s.bitset.SetAll(1)

	bits := s.bitset.GetAll()
	t.Assert(len(bits), Equals, kBitsetSize)

	for i := 0; i < kBitsetSize; i++ {
		t.Assert(bits[i], Equals, 1)
	}
}
