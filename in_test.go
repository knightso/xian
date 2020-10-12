package xian

import (
	"reflect"
	"testing"
)

func TestInBuilderBit(t *testing.T) {
	inBuilder := NewInBuilder()

	assert := func(actual, expected Bit) {
		t.Helper()

		if actual != expected {
			t.Errorf("unexpected, actual: `%v`, expected: `%v`", actual, expected)
		}
	}

	assert(inBuilder.NewBit(), 1)
	assert(inBuilder.NewBit(), 2)
	assert(inBuilder.NewBit(), 4)
	assert(inBuilder.NewBit(), 8)
	assert(inBuilder.NewBit(), 16)
	assert(inBuilder.NewBit(), 32)
	assert(inBuilder.NewBit(), 64)
	assert(inBuilder.NewBit(), 128)
	assert(inBuilder.NewBit(), 256)

	uintSize := 16

	for i := 9; i < uintSize-1; i++ {
		assert(inBuilder.NewBit(), Bit(1<<uint(i)))
	}

	assert(inBuilder.NewBit(), 1<<(uint(uintSize)-1))

	// overflow
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Errorf("panic expected")
			}
		}()

		inBuilder.NewBit()
	}()
}

func TestInBuilderIndexes(t *testing.T) {
	inBuilder := NewInBuilder()

	a := inBuilder.NewBit()
	b := inBuilder.NewBit()
	c := inBuilder.NewBit()
	d := inBuilder.NewBit()

	idxs := inBuilder.Indexes()

	if len(idxs) != 0 {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", len(idxs), 0)
	}

	idxs = inBuilder.Indexes(a, c)

	expected := []string{"1", "3", "4", "5", "6", "7", "9", "b", "c", "d", "e", "f"}
	if !reflect.DeepEqual(idxs, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", idxs, expected)
	}

	idxs = inBuilder.Indexes(b, d)

	expected = []string{"2", "3", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	if !reflect.DeepEqual(idxs, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", idxs, expected)
	}
}

func TestInBuilderFilters(t *testing.T) {
	inBuilder := NewInBuilder()

	a := inBuilder.NewBit()
	b := inBuilder.NewBit()
	c := inBuilder.NewBit()
	d := inBuilder.NewBit()

	filter := inBuilder.Filter(a, c)
	expected := "5"

	if filter != expected {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", filter, expected)
	}

	filter = inBuilder.Filter(b, d)
	expected = "a"

	if filter != expected {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", filter, expected)
	}
}
