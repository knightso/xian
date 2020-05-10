package xian

import (
	"sort"
	"testing"
)

func TestValidateConfig(t *testing.T) {
	labels := make([]string, MaxCompositeIndexLabels+1)
	for i := 0; i < len(labels); i++ {
		labels[i] = string('a' + i)
	}

	t.Run("len(CompositeIdxLabels)<=MaxCompositeIndexLabels", func(t *testing.T) {
		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels]}
		if _, err := ValidateConfig(conf); err != nil {
			t.Errorf("expected:error = nil, but was:[%v]\n", err)
		}
		if validated, _ := ValidateConfig(conf); validated != conf {
			t.Errorf("validated, _ := ValidateConfig(conf) expected:validated = conf, but was:validated = %#v\n", validated)
		}
	})

	t.Run("len(CompositeIdxLabels)>MaxCompositeIndexLabels", func(t *testing.T) {
		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels+1]}
		if _, err := ValidateConfig(conf); err == nil {
			t.Error("CompositeIdxLabels > MaxCompositeIndexLabels expected:err != nil, but was:err = nil\n")
		}
	})

	t.Run("ValidateConfig(DefaultConfig)", func(t *testing.T) {
		if _, err := ValidateConfig(DefaultConfig); err != nil {
			t.Errorf("expected:error = nil, but was:error = [%v]\n", err)
		}
	})
}

func TestMustValidateConfig(t *testing.T) {
	labels := make([]string, MaxCompositeIndexLabels+1)
	for i := 0; i < len(labels); i++ {
		labels[i] = string('a' + i)
	}

	t.Run("CompositeIdxLabels<=MaxCompositeIndexLabels", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("expected:not panic, was:panic = [%v]\n", r)
			}
		}()

		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels]}
		MustValidateConfig(conf)
	})

	t.Run("CompositeIdxLabels>MaxCompositeIndexLabels", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected:panic, was:not panic\n")
			}
		}()

		conf := &Config{CompositeIdxLabels: labels[:MaxCompositeIndexLabels+1]}
		MustValidateConfig(conf)
	})
}

func TestBiunigrams(t *testing.T) {
	result := Biunigrams("abc dあいbCh")
	if len(result) != 15 {
		t.Errorf("len(result) expected:%d, but was:%d\n", 13, len(result))
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

func TestSuffixes(t *testing.T) {
	result := Suffixes("abc dあいbCh")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 9, len(result))
	}

	sort.Strings(result)

	assert(t, "result[0]", result[0], "c")
	assert(t, "result[1]", result[1], "cb")
	assert(t, "result[2]", result[2], "cba")
	assert(t, "result[3]", result[3], "h")
	assert(t, "result[4]", result[4], "hC")
	assert(t, "result[5]", result[5], "hCb")
	assert(t, "result[6]", result[6], "hCbい")
	assert(t, "result[7]", result[7], "hCbいあ")
	assert(t, "result[8]", result[8], "hCbいあd")
}

func TestAddIndexAndFilter(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")

	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")

	builtIndexes := idx.MustBuild()
	builtFilters := filter.MustBuild()

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

	builtIndexes := idx.MustBuild()
	builtFilters := filter.MustBuild()

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

	builtIndexes := idx.MustBuild()
	builtFilters := filter.MustBuild()

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
