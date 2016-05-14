package ga

// BitsetCreate - an interface to a bitset create struct
type BitsetCreate interface {
	Go() Bitset
}

// NullBitsetCreate - a null implementation of the BitsetCreate interface
type NullBitsetCreate struct {
}

// Go returns a bitset with no content
func (ngc *NullBitsetCreate) Go() Bitset {
	return Bitset{}
}
