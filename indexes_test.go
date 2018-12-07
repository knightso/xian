package xian

import (
	"reflect"
	"sort"
	"testing"
)

func TestAddIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.Add("label2", "abc debch iJあdeN", "sample")

	built := idx.Build()
	assertBuiltIndex(t, built, []string{
		"label1 abc dあいbCh",
		"label1 sample",
		"label2 abc debch iJあdeN",
		"label2 sample",
	})
}

func TestAddBigrams(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBigrams("label1", "abc dあいbCh")
	idx.AddBigrams("label2", "abc debch iJあdeN")

	var expected []string
	for _, s := range Bigrams("abc dあいbCh") {
		expected = append(expected, "label1 "+s)
	}
	for _, s := range Bigrams("abc debch iJあdeN") {
		expected = append(expected, "label2 "+s)
	}

	built := idx.Build()
	assertBuiltIndex(t, built, expected)
}

func TestAddBiunigramsIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBiunigrams("label1", "abc dあいbCh")
	idx.AddBiunigrams("label2", "abc debch iJあdeN")

	var expected []string
	for _, s := range Biunigrams("abc dあいbCh") {
		expected = append(expected, "label1 "+s)
	}
	for _, s := range Biunigrams("abc debch iJあdeN") {
		expected = append(expected, "label2 "+s)
	}

	built := idx.Build()
	assertBuiltIndex(t, built, expected)
}

func TestAddPrefixesIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddPrefixes("label1", "abc dあいbCh")
	idx.AddPrefixes("label2", "abc debch iJあdeN")

	var expected []string
	for _, s := range Prefixes("abc dあいbCh") {
		expected = append(expected, "label1 "+s)
	}
	for _, s := range Prefixes("abc debch iJあdeN") {
		expected = append(expected, "label2 "+s)
	}

	built := idx.Build()
	assertBuiltIndex(t, built, expected)
}

func TestAddAllIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.AddBigrams("label2", "abc dあいbCh")
	idx.AddBiunigrams("label3", "abc dあいbCh")
	idx.AddPrefixes("label4", "abc dあいbCh")

	var expected []string

	// Add
	expected = append(expected, "label1 abc dあいbCh")
	expected = append(expected, "label1 sample")

	// AddBigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected = append(expected, "label2 "+s)
	}

	// AddBiunigrams
	for _, s := range Biunigrams("abc dあいbCh") {
		expected = append(expected, "label3 "+s)
	}

	// AddPrefixes
	for _, s := range Prefixes("abc dあいbCh") {
		expected = append(expected, "label4 "+s)
	}

	built := idx.Build()
	assertBuiltIndex(t, built, expected)
}

func TestIndexConfigCompositeIdxLabels(t *testing.T) {
	idx := NewIndexes(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	idx.Add("label1", "a")
	idx.Add("label2", "b")
	idx.Add("label3", "c")
	idx.Add("label4", "d")

	built := idx.Build()

	//   c b a
	//  ------
	//3  0 1 1
	//5  1 0 1
	//6  1 1 0
	//7  1 1 1
	assertBuiltIndex(t, built, []string{
		"label1 a",
		"label2 b",
		"label3 c",
		"label4 d",
		"3 a;b",
		"5 a;c",
		"6 b;c",
		"7 a;b;c",
	})
}

func TestIndexConfigSaveNoFiltersIndex(t *testing.T) {
	idx := NewIndexes(&Config{SaveNoFiltersIndex: true})
	idx.Add("label1", "a")

	built := idx.Build()
	assertBuiltIndex(t, built, []string{
		"label1 a",
		"__NoFilters__",
	})
}

func assertBuiltIndex(t *testing.T, actual, expected []string) {
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", actual, expected)
	}
}
