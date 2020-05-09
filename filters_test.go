package xian

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestAddFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.Add("label2", "abc debch iJあdeN", "sample")

	built := filter.MustBuild()
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

		built := filter.MustBuild()
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

		built := filter.MustBuild()
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

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddPrefixFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddPrefix("label1", "abc dあいbCh")
	filter.AddPrefix("label2", "abc debch iJあdeN")

	built := filter.MustBuild()
	assertBuiltFilter(t, built, []string{
		"label1 abc dあいbCh",
		"label2 abc debch iJあdeN",
	})
}

func TestAddSomethingFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.AddSomething("label1", []string{"abc dあいbCh", "abc debch iJあdeN"})
	filter.AddSomething("label2", 123)
	now := time.Now()
	filter.AddSomething("label3", now)

	built := filter.MustBuild()
	assertBuiltFilter(t, built, []string{
		"label1 abc dあいbCh",
		"label1 abc debch iJあdeN",
		"label2 123",
		fmt.Sprintf("label3 %d", now.Unix()),
	})
}

func TestAddAllFilter(t *testing.T) {
	filter := NewFilters(nil)
	filter.Add("label1", "abc dあいbCh", "sample")
	filter.AddBigrams("label2", "abc dあいbCh")
	filter.AddBiunigrams("label3", "abc dあいbCh")
	filter.AddPrefix("label4", "abc dあいbCh")
	filter.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

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

	// AddPrefix
	expected = append(expected, "label4 abc dあいbCh")

	// AddSomething
	expected = append(expected, "label5 abc dあいbCh")
	expected = append(expected, "label5 AbcdeF")

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestBuildFilter(t *testing.T) {
	t.Run("ConfigのCompositeIdxLabelsがMaxCompositeIndexLabelsより大きい場合", func(t *testing.T) {
		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string('a' + i)
		}

		filter := NewFilters(&Config{CompositeIdxLabels: labels})
		if _, err := filter.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})

	t.Run("Buildの結果件数がMaxIndexesSizeより大きい場合", func(t *testing.T) {
		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		if _, err := filter.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})

	t.Run("Config,結果件数が正常な場合", func(t *testing.T) {
		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built, err := filter.Build()
		if err != nil {
			t.Errorf("error = %s, wants = nil", err)
		}
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})
}

func TestMustBuildFilter(t *testing.T) {
	t.Run("ConfigのCompositeIdxLabelsがMaxCompositeIndexLabelsより大きい場合", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected:panic, was:not panic\n")
			}
		}()

		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string('a' + i)
		}

		filter := NewFilters(&Config{CompositeIdxLabels: labels})
		filter.MustBuild()
	})

	t.Run("Buildの結果件数がMaxIndexesSizeより大きい場合", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected:panic, was:not panic\n")
			}
		}()

		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		filter.MustBuild()
	})

	t.Run("Config,結果件数が正常な場合", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("expected:not panic, was:panic = [%v]\n", r)
			}
		}()

		filter := NewFilters(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			filter.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built := filter.MustBuild()
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})
}

func TestFilterConfigCompositeIdxLabels(t *testing.T) {
	filter := NewFilters(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	filter.Add("label1", "a")
	filter.Add("label2", "b")
	filter.Add("label3", "c")

	built := filter.MustBuild()

	//   c b a
	//  ------
	//3  0 1 1
	//5  1 0 1
	//6  1 1 0
	//7  1 1 1
	assertBuiltIndex(t, built, []string{
		// "3 a;b",
		// "5 a;c",
		// "6 b;c",

		// now needs only indexes with all specified labels
		"7 a;b;c",
	})
}

func TestFilterConfigIgnoreCase(t *testing.T) {
	filter := NewFilters(&Config{IgnoreCase: true})
	filter.Add("label1", "abc dあいbCh", "saMPle")
	filter.AddBigrams("label2", "abc dあいbCh")
	filter.AddBiunigrams("label3", "abc dあいbCh")
	filter.AddPrefix("label4", "abc dあいbCh")
	filter.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	var expected []string

	// Add
	expected = append(expected, "label1 abc dあいbch")
	expected = append(expected, "label1 sample")

	// AddBigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected = append(expected, "label2 "+s)
	}

	// AddBiunigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected = append(expected, "label3 "+s)
	}

	// AddPrefix
	expected = append(expected, "label4 abc dあいbch")

	// AddSomething
	expected = append(expected, "label5 abc dあいbch")
	expected = append(expected, "label5 abcdef")

	built := filter.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestFilterConfigSaveNoFiltersIndex(t *testing.T) {
	t.Run("1つもfilterをAddしていない場合", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})

		built := filter.MustBuild()

		assertBuiltIndex(t, built, []string{
			"__NoFilters__",
		})
	})

	t.Run("1つ以上filterをAddした場合", func(t *testing.T) {
		filter := NewFilters(&Config{SaveNoFiltersIndex: true})
		filter.Add("label1", "a")

		built := filter.MustBuild()

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
