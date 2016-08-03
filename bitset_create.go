package ga

import "github.com/tomcraven/bitset"

// BitsetCreate - an interface to a bitset create struct
type BitsetCreate interface {
	Go() bitset.Bitset
}

// NullBitsetCreate - a null implementation of the BitsetCreate interface
type NullBitsetCreate struct {
}

// Go returns a bitset with no content
func (ngc *NullBitsetCreate) Go() bitset.Bitset {
	return bitset.Create(0)
}
