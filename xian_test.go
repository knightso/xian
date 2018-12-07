package xian

import (
	"sort"
	"testing"
)

func TestBiunigrams(t *testing.T) {
	result := Biunigrams("abc dあいbCh")
	if len(result) != 15 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 13, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "C")
	assert(t, "result[1]", result[1], "Ch")
	assert(t, "result[2]", result[2], "a")
	assert(t, "result[3]", result[3], "ab")
	assert(t, "result[4]", result[4], "b")
	assert(t, "result[5]", result[5], "bC")
	assert(t, "result[6]", result[6], "bc")
	assert(t, "result[7]", result[7], "c")
	assert(t, "result[8]", result[8], "d")
	assert(t, "result[9]", result[9], "dあ")
	assert(t, "result[10]", result[10], "h")
	assert(t, "result[11]", result[11], "あ")
	assert(t, "result[12]", result[12], "あい")
	assert(t, "result[13]", result[13], "い")
	assert(t, "result[14]", result[14], "いb")
}

func TestBigrams(t *testing.T) {
	result := Bigrams("abc dあいbCh")
	if len(result) != 7 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 6, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "Ch")
	assert(t, "result[1]", result[1], "ab")
	assert(t, "result[2]", result[2], "bC")
	assert(t, "result[3]", result[3], "bc")
	assert(t, "result[4]", result[4], "dあ")
	assert(t, "result[5]", result[5], "あい")
	assert(t, "result[6]", result[6], "いb")
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
	assert(t, "result[7]", result[7], "dあいbC")
	assert(t, "result[8]", result[8], "dあいbCh")
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
