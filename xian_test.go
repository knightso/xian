package xian

import (
	"testing"
)

func TestBiunigrams(t *testing.T) {
	result := Biunigrams("abc dあいbCh")
	if len(result) != 13 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 13, len(result))
	}

	expectedGrams := []string{
		"ab", "bc", "dあ", "あい", "いb", "ch",
		"a", "b", "c", "d", "あ", "い", "h"}
	for _, gram := range result {
		flag := false
		for _, exp := range expectedGrams {
			if gram == exp {
				flag = true
				break
			}
		}
		if !flag {
			t.Errorf("unexpected gram: %s", gram)
		}
	}
}

func TestBigrams(t *testing.T) {
	result := Bigrams("abc dあいbCh")
	if len(result) != 6 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 6, len(result))
	}

	expectedGrams := []string{
		"ab", "bc", "dあ", "あい", "いb", "ch"}
	for _, gram := range result {
		flag := false
		for _, exp := range expectedGrams {
			if gram == exp {
				flag = true
				break
			}
		}
		if !flag {
			t.Errorf("unexpected gram: %s", gram)
		}
	}
}

func TestPrefixes(t *testing.T) {
	result := Prefixes("abc dあいbCh")
	if len(result) != 2 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 2, len(result))
	}
	// TODO: [a ab abc d d? d? dあ dあ? dあ? dあい dあいb dあいbC dあいbCh]
}
