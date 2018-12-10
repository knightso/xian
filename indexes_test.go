package xian

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

func TestAddIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.Add("label2", "abc debch iJあdeN", "sample")

	built := idx.MustBuild()
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

	built := idx.MustBuild()
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

	built := idx.MustBuild()
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

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestAddSomethingIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.AddSomething("label1", []string{"abc dあいbCh", "abc debch iJあdeN"})
	idx.AddSomething("label2", 123)

	built := idx.MustBuild()
	assertBuiltIndex(t, built, []string{
		"label1 abc dあいbCh",
		"label1 abc debch iJあdeN",
		"label2 123",
	})
}

func TestAddAllIndex(t *testing.T) {
	idx := NewIndexes(nil)
	idx.Add("label1", "abc dあいbCh", "sample")
	idx.AddBigrams("label2", "abc dあいbCh")
	idx.AddBiunigrams("label3", "abc dあいbCh")
	idx.AddPrefixes("label4", "abc dあいbCh")
	idx.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

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

	// AddSomething
	expected = append(expected, "label5 abc dあいbCh")
	expected = append(expected, "label5 AbcdeF")

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)
}

func TestBuildIndex(t *testing.T) {
	t.Run("ConfigのCompositeIdxLabelsがMaxCompositeIndexLabelsより大きい場合", func(t *testing.T) {
		labels := make([]string, MaxCompositeIndexLabels+1)
		for i := 0; i < len(labels); i++ {
			labels[i] = string('a' + i)
		}

		idx := NewIndexes(&Config{CompositeIdxLabels: labels})
		if _, err := idx.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})

	t.Run("Buildの結果件数がMaxIndexesSizeより大きい場合", func(t *testing.T) {
		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		if _, err := idx.Build(); err == nil {
			t.Error("error = nil, wants != nil")
		}
	})

	t.Run("Config,結果件数が正常な場合", func(t *testing.T) {
		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built, err := idx.Build()
		if err != nil {
			t.Errorf("error = %s, wants = nil", err)
		}
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})
}

func TestMustBuildIndex(t *testing.T) {
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

		idx := NewIndexes(&Config{CompositeIdxLabels: labels})
		idx.MustBuild()
	})

	t.Run("Buildの結果件数がMaxIndexesSizeより大きい場合", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected:panic, was:not panic\n")
			}
		}()

		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize+1; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		idx.MustBuild()
	})

	t.Run("Config,結果件数が正常な場合", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("expected:not panic, was:panic = [%v]\n", r)
			}
		}()

		idx := NewIndexes(nil)
		for i := 0; i < MaxIndexesSize; i++ {
			idx.Add(fmt.Sprintf("label%d", i), "abc")
		}

		built := idx.MustBuild()
		if built == nil {
			t.Error("built = nil, wants != nil")
		}
	})
}

func TestIndexConfigCompositeIdxLabels(t *testing.T) {
	idx := NewIndexes(&Config{CompositeIdxLabels: []string{"label1", "label2", "label3"}})
	idx.Add("label1", "a")
	idx.Add("label2", "b")
	idx.Add("label3", "c")
	idx.Add("label4", "d")

	built := idx.MustBuild()

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

func TestIndexConfigIgnoreCase(t *testing.T) {
	idx := NewIndexes(&Config{IgnoreCase: true})

	idx.Add("label1", "abc dあいbCh", "saMPle")
	idx.AddBigrams("label2", "abc dあいbCh")
	idx.AddBiunigrams("label3", "abc dあいbCh")
	idx.AddPrefixes("label4", "abc dあいbCh")
	idx.AddSomething("label5", []string{"abc dあいbCh", "AbcdeF"})

	var expected []string

	// Add
	expected = append(expected, "label1 abc dあいbch")
	expected = append(expected, "label1 sample")

	// AddBigrams
	for _, s := range Bigrams("abc dあいbch") {
		expected = append(expected, "label2 "+s)
	}

	// AddBiunigrams
	for _, s := range Biunigrams("abc dあいbch") {
		expected = append(expected, "label3 "+s)
	}

	// AddPrefixes
	for _, s := range Prefixes("abc dあいbch") {
		expected = append(expected, "label4 "+s)
	}

	// AddSomething
	expected = append(expected, "label5 abc dあいbch")
	expected = append(expected, "label5 abcdef")

	built := idx.MustBuild()
	assertBuiltIndex(t, built, expected)

}

func TestIndexConfigSaveNoFiltersIndex(t *testing.T) {
	idx := NewIndexes(&Config{SaveNoFiltersIndex: true})
	idx.Add("label1", "a")

	built := idx.MustBuild()
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
