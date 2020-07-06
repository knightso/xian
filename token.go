package xian

import (
	"fmt"
	"strings"
)

type bigram struct {
	a, b rune
}

func (b *bigram) String() string {
	return fmt.Sprintf("%c%c", b.a, b.b)
}

// Biunigrams returns bigram and unigram tokens from s.
func Biunigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}
	for unigram := range toUnigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c", unigram))
	}

	return tokens
}

// Bigrams returns bigram tokens from s.
func Bigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(s) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	return tokens
}

// Prefixes returns prefix tokens from s.
func Prefixes(s string) []string {
	prefixes := make(map[string]struct{})

	runes := make([]rune, 0, 64)

	for _, w := range strings.Split(s, " ") {
		if w == "" {
			continue
		}

		runes = runes[0:0]

		for _, c := range w {
			runes = append(runes, c)
			prefixes[string(runes)] = struct{}{}
		}
	}

	tokens := make([]string, 0, 32)

	for pref := range prefixes {
		tokens = append(tokens, pref)
	}

	return tokens
}

// Suffixes returns suffix tokens from s.
func Suffixes(s string) []string {
	suffixes := make(map[string]struct{})

	runes := make([]rune, 0, 64)

	for _, w := range strings.Split(s, " ") {
		if w == "" {
			continue
		}
		w = reverse(w)

		runes = runes[0:0]

		for _, c := range w {
			runes = append(runes, c)
			suffixes[string(runes)] = struct{}{}
		}
	}

	tokens := make([]string, 0, 32)

	for suf := range suffixes {
		tokens = append(tokens, suf)
	}

	return tokens
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func toBigrams(value string) map[bigram]bool {
	result := make(map[bigram]bool)
	var prev rune
	for i, r := range value {
		if i > 0 && prev != ' ' && r != ' ' {
			result[bigram{prev, r}] = true
		}
		prev = r
	}
	return result
}

func toUnigrams(value string) map[rune]bool {
	result := make(map[rune]bool)
	for _, r := range value {
		if r == ' ' {
			continue
		}
		result[r] = true
	}
	return result
}
