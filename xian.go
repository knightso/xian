package xian

import (
	"fmt"
	"strings"
)

// Indexes is extra indexes for datastore query.
type Indexes map[string][]string

// NewIndexes returns a new Indexes.
func NewIndexes() Indexes {
	return make(Indexes)
}

// Add adds new indexes with a lavel.
func (idx Indexes) Add(label string, indexes ...string) {
	idx[label] = append(idx[label], indexes...)
}

// Biunigrams returns bigram and unigram tokens from s.
func Biunigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}
	for unigram := range toUnigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c", unigram))
	}

	return tokens
}

// Bigrams returns bigram tokens from s.
func Bigrams(s string) []string {
	tokens := make([]string, 0, 32)

	for bigram := range toBigrams(strings.ToLower(s)) {
		tokens = append(tokens, fmt.Sprintf("%c%c", bigram.a, bigram.b))
	}

	return tokens
}

// Prefixes returns prefix tokens from s.
func Prefixes(s string) []string {
	prefixes := make(map[string]struct{})

	runes := make([]rune, 0, 64)

	for _, w := range strings.Split(strings.ToLower(s), " ") {
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
