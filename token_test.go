package xian

import (
	"testing"
)

func TestString(t *testing.T) {
	b := bigram{'a', 'あ'}
	if b.String() != "aあ" {
		t.Errorf("exected:%s, but was:%s\n", "aあ", b.String())
	}
}

func TestToUnigrams(t *testing.T) {
	result := toUnigrams("abc dあいbCh")
	if len(result) != 8 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 8, len(result))
	}

	if !result['a'] {
		t.Errorf("Unigram notfound. 'a'")
	}
	if !result['b'] {
		t.Errorf("Unigrbm notfound. 'b'")
	}
	if !result['c'] {
		t.Errorf("Unigrbm notfound. 'c'")
	}
	if !result['C'] {
		t.Errorf("Unigrbm notfound. 'C'")
	}
	if !result['d'] {
		t.Errorf("Unigrbm notfound. 'd'")
	}
	if !result['あ'] {
		t.Errorf("Unigrbm notfound. 'あ'")
	}
	if !result['い'] {
		t.Errorf("Unigrbm notfound. 'い'")
	}
	if !result['h'] {
		t.Errorf("Unigrbm notfound. 'h'")
	}
}

func TestToBigrams(t *testing.T) {
	result := toBigrams("abc debch iJあdeN")
	if len(result) != 9 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 9, len(result))
	}

	assertBigram(t, result, bigram{'a', 'b'})
	assertBigram(t, result, bigram{'b', 'c'})
	assertBigram(t, result, bigram{'d', 'e'})
	assertBigram(t, result, bigram{'e', 'b'})
	assertBigram(t, result, bigram{'c', 'h'})
	assertBigram(t, result, bigram{'i', 'J'})
	assertBigram(t, result, bigram{'J', 'あ'})
	assertBigram(t, result, bigram{'あ', 'd'})
	assertBigram(t, result, bigram{'e', 'N'})
}

func assertBigram(t *testing.T, set map[bigram]bool, bigram bigram) {
	if !set[bigram] {
		t.Errorf("Bigram notfound. %v\n", bigram)
	}
}
