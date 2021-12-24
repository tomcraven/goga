package goga

// BitsetParse - an interface to an object that is able
// to parse a bitset into an array of uint64s
type BitsetParse interface {
	SetFormat([]int)
	Process(*Bitset) []uint64
}

type bitsetParse struct {
	expectedBitsetSize int
	format             []int
}

// CreateBitsetParse returns an instance of a bitset parser
func CreateBitsetParse() BitsetParse {
	return &bitsetParse{}
}

func (bp *bitsetParse) SetFormat(format []int) {
	bp.expectedBitsetSize = 0
	for _, i := range format {
		bp.expectedBitsetSize += i
	}
	bp.format = format
}

func (bp *bitsetParse) Process(bitset *Bitset) []uint64 {
	if bitset.GetSize() != bp.expectedBitsetSize {
		panic("Input format does not match bitset size")
	}

	ret := make([]uint64, len(bp.format))
	runningBits := 0
	for retIndex, numBits := range bp.format {
		ret[retIndex] = 0

		for i := 0; i < numBits; i++ {
			if bitset.Get(i+runningBits) == 1 {
				ret[retIndex] |= (1 << uint(i))
			}
		}

		runningBits += numBits
	}
	return ret
}
