package ga

// IBitsetCreate - an interface to a bitset create struct
type IBitsetCreate interface {
	Go() Bitset
}

// NullBitsetCreate - a null implementation of the IBitsetCreate interface
type NullBitsetCreate struct {
}

// Go returns a bitset with no content
func (ngc *NullBitsetCreate) Go() Bitset {
	return Bitset{}
}
