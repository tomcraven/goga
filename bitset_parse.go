package ga

import "github.com/tomcraven/bitset"

// IBitsetParse - an interface to an object that is able
// to parse a bitset into an array of uint64s
type IBitsetParse interface {
	SetFormat([]uint)
	Process(bitset.Bitset) []uint64
}

type bitsetParse struct {
	expectedBitsetSize uint
	format             []uint
}

// CreateBitsetParse returns an instance of a bitset parser
func CreateBitsetParse() IBitsetParse {
	return &bitsetParse{}
}

func (bp *bitsetParse) SetFormat(format []uint) {
	bp.expectedBitsetSize = 0
	for _, i := range format {
		bp.expectedBitsetSize += i
	}
	bp.format = format
}

func (bp *bitsetParse) Process(bitset bitset.Bitset) []uint64 {
	if bitset.Size() != bp.expectedBitsetSize {
		panic("Input format does not match bitset size")
	}

	ret := make([]uint64, len(bp.format))
	runningBits := uint(0)
	for retIndex, numBits := range bp.format {
		ret[retIndex] = 0

		for i := uint(0); i < numBits; i++ {
			if bitset.Get(i + runningBits) {
				ret[retIndex] |= (1 << uint(i))
			}
		}

		runningBits += numBits
	}
	return ret
}
