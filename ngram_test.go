package xian

import (
	"log"
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
	if len(result) != 7 {
		t.Errorf("len(result) exected:%d, but was:%d\n", 7, len(result))
	}

	log.Println("1111")
	log.Println(result)

	if !result['a'] {
		t.Errorf("Unigram notfound. 'a'")
	}
	if !result['b'] {
		t.Errorf("Unigrbm notfound. 'b'")
	}
	if !result['c'] {
		t.Errorf("Unigrbm notfound. 'c'")
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

	log.Println("2222")
	log.Println(result)

	assertBigram(t, result, bigram{'a', 'b'})
	assertBigram(t, result, bigram{'b', 'c'})
	assertBigram(t, result, bigram{'d', 'e'})
	assertBigram(t, result, bigram{'e', 'b'})
	assertBigram(t, result, bigram{'c', 'h'})
	assertBigram(t, result, bigram{'i', 'j'})
	assertBigram(t, result, bigram{'j', 'あ'})
	assertBigram(t, result, bigram{'あ', 'd'})
	assertBigram(t, result, bigram{'e', 'n'})
}

func assertBigram(t *testing.T, set map[bigram]bool, bigram bigram) {
	if !set[bigram] {
		t.Errorf("Bigram notfound. %v\n", bigram)
	}
}
