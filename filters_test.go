package xian

import (
	"reflect"
	"sort"
	"testing"
)

func TestAddFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.Add("label2", "abc debch iJあdeN", "sample")

	built := filter.Build()
	assertBuiltFilter(t, built, []string{
		"label1 abc dあいbCh",
		"label1 sample",
		"label2 abc debch iJあdeN",
		"label2 sample",
	})
}

func TestAddBigramsFilter(t *testing.T) {
	t.Run("indexes が1文字の場合", func(t *testing.T) {
		filter := NewFilters(nil)
		filter.AddBigrams("label1", "a")
		filter.AddBigrams("label2", "b")

		var expected []string
		expected = append(expected, "label1 "+"a")
		expected = append(expected, "label2 "+"b")

		built := filter.Build()
		assertBuiltIndex(t, built, expected)
	})

	t.Run("indexes が2文字以上の場合", func(t *testing.T) {
		filter := NewFilters(nil)
		filter.AddBigrams("label1", "abc dあいbCh")
		filter.AddBigrams("label2", "abc debch iJあdeN")

		var expected []string
		for _, s := range Bigrams("abc dあいbCh") {
			expected = append(expected, "label1 "+s)
		}
		for _, s := range Bigrams("abc debch iJあdeN") {
			expected = append(expected, "label2 "+s)
		}

		built := filter.Build()
		assertBuiltIndex(t, built, expected)
	})
}

func TestAddBiunigramsFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddBiunigrams("label1", "abc dあいbCh")
	filter.AddBiunigrams("label2", "abc debch iJあdeN")

	// Filters の場合は Bigrams を使用する
	var expected []string
	for _, s := range Bigrams("abc dあいbCh") {
		expected = append(expected, "label1 "+s)
	}
	for _, s := range Bigrams("abc debch iJあdeN") {
		expected = append(expected, "label2 "+s)
	}

	built := filter.Build()
	assertBuiltIndex(t, built, expected)
}

func TestAddAllFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.AddBigrams("label2", "abc dあいbCh")
	filter.AddBiunigrams("label3", "abc dあいbCh")

	var expected []string

	// Add
	expected = append(expected, "label1 abc dあいbCh")
	expected = append(expected, "label1 sample")

	// AddBigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected = append(expected, "label2 "+s)
	}

	// AddBiunigrams
	for _, s := range Bigrams("abc dあいbCh") {
		expected = append(expected, "label3 "+s)
	}

	built := filter.Build()
	assertBuiltIndex(t, built, expected)
}

func TestFilterConfigCompositeIdxLabels(t *testing.T) {
	filter := NewFilters(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	filter.Add("label1", "a")
	filter.Add("label2", "b")
	filter.Add("label3", "c")

	built := filter.Build()

	//   c b a
	//  ------
	//3  0 1 1
	//5  1 0 1
	//6  1 1 0
	//7  1 1 1
	assertBuiltIndex(t, built, []string{
		"3 a;b",
		"5 a;c",
		"6 b;c",
		"7 a;b;c",
	})
}

func TestFilterConfigSaveNoFiltersIndex(t *testing.T) {
	t.Run("1つもfilterをAddしていない場合", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})

		built := filter.Build()

		assertBuiltIndex(t, built, []string{
			"__NoFilters__",
		})
	})

	t.Run("1つ以上filterをAddした場合", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})
		filter.Add("label1", "a")

		built := filter.Build()

		assertBuiltIndex(t, built, []string{
			"label1 a",
		})
	})
}

func assertBuiltFilter(t *testing.T, actual, expected []string) {
	sort.Strings(actual)
	sort.Strings(expected)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected, actual: `%v`, expected: `%v`", actual, expected)
	}
}
