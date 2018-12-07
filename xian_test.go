package xian

import (
	"sort"
	"testing"
)

func TestBiunigrams(t *testing.T) {
	result := Biunigrams("abc dあいbCh")
	if len(result) != 13 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 13, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "a")
	assert(t, "result[1]", result[1], "ab")
	assert(t, "result[2]", result[2], "b")
	assert(t, "result[3]", result[3], "bc")
	assert(t, "result[4]", result[4], "c")
	assert(t, "result[5]", result[5], "ch")
	assert(t, "result[6]", result[6], "d")
	assert(t, "result[7]", result[7], "dあ")
	assert(t, "result[8]", result[8], "h")
	assert(t, "result[9]", result[9], "あ")
	assert(t, "result[10]", result[10], "あい")
	assert(t, "result[11]", result[11], "い")
	assert(t, "result[12]", result[12], "いb")
}

func TestBigrams(t *testing.T) {
	result := Bigrams("abc dあいbCh")
	if len(result) != 6 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 6, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "ab")
	assert(t, "result[1]", result[1], "bc")
	assert(t, "result[2]", result[2], "ch")
	assert(t, "result[3]", result[3], "dあ")
	assert(t, "result[4]", result[4], "あい")
	assert(t, "result[5]", result[5], "いb")
}

func TestPrefixes(t *testing.T) {
	result := Prefixes("abc dあいbCh")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 9, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "a")
	assert(t, "result[1]", result[1], "ab")
	assert(t, "result[2]", result[2], "abc")
	assert(t, "result[3]", result[3], "d")
	assert(t, "result[4]", result[4], "dあ")
	assert(t, "result[5]", result[5], "dあい")
	assert(t, "result[6]", result[6], "dあいb")
	assert(t, "result[7]", result[7], "dあいbc")
	assert(t, "result[8]", result[8], "dあいbch")
}

func TestAddIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")

	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")

	builtIndexes := idx.Build()
	builtFilters := filter.Build()

	// filter の内容が全て index に存在すること
	for _, builtFilter := range builtFilters {
		if !containsString(builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func TestAddBigramsIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBigrams("label1", "abc dあいbCh")

	filter := NewFilters(nil)
	filter.AddBigrams("label1", "dあいb") // idx の中間一致

	builtIndexes := idx.Build()
	builtFilters := filter.Build()

	// filter の内容が全て index に存在すること
	for _, builtFilter := range builtFilters {
		if !containsString(builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func TestAddBiunigramsIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddBiunigrams("label1", "abc dあいbCh")

	filter := NewFilters(nil)
	filter.AddBiunigrams("label1", "dあいb") // idx の中間一致

	builtIndexes := idx.Build()
	builtFilters := filter.Build()

	// filter の内容が全て index に存在すること
	for _, builtFilter := range builtFilters {
		if !containsString(builtIndexes, builtFilter) {
			t.Errorf("filter: %s not contains", builtFilter)
		}
	}
}

func assert(t *testing.T, title string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%s : unexpected, actual: `%v`, expected: `%v`", title, actual, expected)
	}
}

func containsString(list []string, target string) bool {
	for _, s := range list {
		if s == target {
			return true
		}
	}
	return false
}
