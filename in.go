package xian

import "fmt"

// Bit describes In-Filter mask bit
type Bit uint16

// InBuilder creates Bit for In-Filter
type InBuilder struct {
	nextBit Bit
}

// NewInBuilder creates InBuilder
func NewInBuilder() *InBuilder {
	return &InBuilder{nextBit: 1}
}

// NewBit returns a new bit shifted.
func (f *InBuilder) NewBit() Bit {
	if f.nextBit == 0 {
		panic("overflow")
	}

	bit := f.nextBit
	f.nextBit <<= 1

	return bit
}

// Indexes creates indexes for In-Filter with multi-bits
func (f *InBuilder) Indexes(bits ...Bit) (indexes []string) {
	allBits := f.combineBits(bits...)

	for i := Bit(1); i <= f.nextBit-1; i++ {
		if i&allBits != 0 {
			indexes = append(indexes, fmt.Sprintf("%x", i))
		}
	}

	return indexes
}

// Indexes creates indexes for In-Filter
func (f *InBuilder) Filter(bits ...Bit) string {
	allBits := f.combineBits(bits...)
	return fmt.Sprintf("%x", allBits)
}

func (f *InBuilder) combineBits(bits ...Bit) (allBits Bit) {
	for _, bit := range bits {
		allBits |= bit
	}
	return allBits
}
