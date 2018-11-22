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

func toBigrams(value string) map[bigram]bool {
	result := make(map[bigram]bool)
	var prev rune
	for i, r := range strings.ToLower(value) {
		if i > 0 && prev != ' ' && r != ' ' {
			result[bigram{prev, r}] = true
		}
		prev = r
	}
	return result
}

func toUnigrams(value string) map[rune]bool {
	result := make(map[rune]bool)
	for _, r := range strings.ToLower(value) {
		if r == ' ' {
			continue
		}
		result[r] = true
	}
	return result
}
