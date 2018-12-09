package xian

import (
	"bytes"
	"fmt"
	"sort"
)

const (
	// IndexNoFilters is an index to be used for no-filters
	IndexNoFilters = "__NoFilters__"
)

const (
	combiIndexSeperator = ";"
)

// Config describe extra indexes configuration.
type Config struct {
	// CompositeIdxLabels is a label list which defines composit indexes to improve the search performance
	CompositeIdxLabels []string
	// IgnoreCase defines whether to ignore case on search
	IgnoreCase bool
	// SaveNoFiltersIndex defines whether to save IndexNoFilters index.
	SaveNoFiltersIndex bool
}

// DefaultConfig is default configuration.
var DefaultConfig = &Config{}

// common indexes map
// key=label, value=index set
type indexesMap map[string]map[string]struct{}

// buildIndexes builds indexes from m.
// m is map[label]tokens.
func buildIndexes(m indexesMap, labelsToExclude []string) []string {
	idxSet := make(map[string]struct{})

	excludeSet := make(map[string]struct{})
	for _, l := range labelsToExclude {
		excludeSet[l] = struct{}{}
	}

	for label, tokens := range m {
		if _, ok := excludeSet[label]; ok {
			continue
		}
		for t := range tokens {
			idxSet[fmt.Sprintf("%s %s", label, t)] = struct{}{}
		}
	}

	built := make([]string, 0, len(idxSet))

	for idx := range idxSet {
		built = append(built, idx)
	}

	sort.Strings(built)

	return built
}

// createCompositeIndexes creates composite indexes of labels from m.
// It reduces zig-zag merge join latency.
// m is indexesMap.
// onlyUsedLabels is used for Filters.
func createCompositeIndexes(labels []string, m indexesMap, onlyUsedLabels bool) []string {

	indexes := make([]string, 0, 64)

	f := func(combi int, index string) {
		indexes = append(indexes, fmt.Sprintf("%d %s", combi, index))
	}

	// generate combination indexes with bit oparation
	// mapping each labels to each bits.

	// construct recursive funcs at first.
	// reverse loop for labels so that the first label will be right-end bit.
	for i := len(labels) - 1; i >= 0; i-- {
		i := i
		prevF := f
		idxLabel := labels[i]

		f = func(combi int, index string) {
			// check process bit for the combi.
			if combi&(1<<uint(i)) != 0 {
				tokens := m[idxLabel]
				for token := range tokens {
					combiIndex := appendCombinationIndex(index, token)
					prevF(combi, combiIndex) // recursive call
				}
			} else {
				// no process bit for the combi.
				prevF(combi, index)
			}
		}
	}

	// now generate indexes.
	for i := 3; i < (1 << uint(len(labels))); i++ {
		if (i & (i - 1)) == 0 {
			// do not save single index
			continue
		}
		f(i, "")
	}

	return indexes
}

func appendCombinationIndex(indexes, index string) string {
	buf := bytes.NewBufferString(indexes)
	if buf.Len() > 0 {
		buf.WriteString(combiIndexSeperator)
	}
	buf.WriteString(index)
	return buf.String()
}
