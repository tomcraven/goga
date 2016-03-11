package ga

type IBitsetCreate interface {
	Go() Bitset
}

type NullBitsetCreate struct {
}

func (ngc *NullBitsetCreate) Go() Bitset {
	return Bitset{}
}
