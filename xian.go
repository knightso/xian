package xian

import (
	"fmt"
	"sort"
)

const (
	// IndexNoFilters is an index to be used for no-filters
	IndexNoFilters = "__NoFilters__"
)

// Config describe extra indexes configuration.
type Config struct {
	// CompositeIdxLabels is a label list which defines composit indexes to improve the search performance
	CompositeIdxLabels []string
	// IgnoreCase defines whether to ignore case on search
	IgnoreCase bool
	// SaveNoFilterIndex defines whether to save IndexNoFilters index.
	SaveNoFilterIndex bool
}

// DefaultConfig is default configuration.
var DefaultConfig = &Config{}

// buildIndexes builds indexes from m.
// ms is map[label]tokens.
func buildIndexes(m map[string][]string) []string {
	idxSet := make(map[string]struct{})

	for label, tokens := range m {
		for _, t := range tokens {
			idxSet[fmt.Sprintf("%s %s", label, t)] = struct{}{}
		}
	}

	// TODO: coposite index

	built := make([]string, 0, len(idxSet))

	for idx := range idxSet {
		built = append(built, idx)
	}

	sort.Strings(built)

	return built
}
