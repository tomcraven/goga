package ga

import (
// "fmt"
)

type Bitset struct {
	size int
	bits []int
}

func (b *Bitset) Create(size int) {
	b.size = size
	b.bits = make([]int, size)
}

func (b *Bitset) GetSize() int {
	return b.size
}

func (b *Bitset) Get(index int) int {
	if index < b.size {
		return b.bits[index]
	}
	return -1
}

func (b *Bitset) GetAll() []int {
	return b.bits
}

func (b *Bitset) setImpl(index, value int) {
	b.bits[index] = value
}

func (b *Bitset) Set(index, value int) bool {
	if index < b.size {
		b.setImpl(index, value)
		return true
	}
	return false
}

func (b *Bitset) SetAll(value int) {
	for i := 0; i < b.size; i++ {
		b.setImpl(i, value)
	}
}

func (b *Bitset) CreateCopy() Bitset {
	newBitset := Bitset{}
	newBitset.Create(b.size)
	for i := 0; i < b.size; i++ {
		newBitset.Set(i, b.Get(i))
	}
	return newBitset
}

func (b *Bitset) Slice(startingBit, size int) Bitset {
	ret := Bitset{}
	ret.Create(size)
	ret.bits = b.bits[startingBit : startingBit+size]
	return ret
}
