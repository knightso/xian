package xian

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/pkg/errors"
)

const (
	// IndexNoFilters is an index to be used for no-filters.
	IndexNoFilters = "__NoFilters__"
	// MaxIndexesSize is maximum size of indexes.
	MaxIndexesSize = 512
	// MaxCompositeIndexLabels maximum number of labels for composite index.
	MaxCompositeIndexLabels = 8
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

var timeKind = reflect.TypeOf(time.Time{}).Kind()

// ValidateConfig validates Config fields.
func ValidateConfig(conf *Config) (*Config, error) {
	if len(conf.CompositeIdxLabels) > MaxCompositeIndexLabels {
		return nil, errors.Errorf("CompositeIdxLabels size exceeds %d", MaxCompositeIndexLabels)
	}
	return conf, nil
}

// MustValidateConfig validates fields and panics if it's invalid.
func MustValidateConfig(conf *Config) *Config {
	conf, err := ValidateConfig(conf)
	if err != nil {
		panic(err)
	}
	return conf
}

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
// forFilters is used for Filters.
func createCompositeIndexes(labels []string, m indexesMap, forFilters bool) ([]string, error) {

	if len(labels) > MaxCompositeIndexLabels {
		return nil, errors.Errorf("CompositeIdxLabels size exceeds %d", MaxCompositeIndexLabels)
	}

	indexes := make([]string, 0, 64)

	f := func(combi uint8, index string, someNew bool) {
		if forFilters && !someNew {
			return
		}
		indexes = append(indexes, fmt.Sprintf("%d %s", combi, index))
	}

	// used indexes sets for filters
	usedIndexes := make(indexesMap)

	// generate combination indexes with bit oparation
	// mapping each labels to each bits.

	// construct recursive funcs at first.
	// reverse loop for labels so that the first label will be right-end bit.
	var combiForFilter uint8
	for i := len(labels) - 1; i >= 0; i-- {
		i := i
		prevF := f
		idxLabel := labels[i]

		if len(m[idxLabel]) > 0 {
			combiForFilter |= 1 << uint(i)
		}

		f = func(combi uint8, index string, someNew bool) {
			// check process bit for the combi.
			if combi&(1<<uint(i)) != 0 {
				tokens := make([]string, 0, len(m[idxLabel]))
				for token := range m[idxLabel] {
					tokens = append(tokens, token)
				}
				for _, token := range tokens {
					combiIndex := appendCombinationIndex(index, token)

					// check if the token is already used for filters
					if _, ok := usedIndexes[idxLabel]; !ok {
						usedIndexes[idxLabel] = make(map[string]struct{})
					}
					if _, ok := usedIndexes[idxLabel][token]; !ok {
						usedIndexes[idxLabel][token] = struct{}{}
						someNew = true
					}

					prevF(combi, combiIndex, someNew) // recursive call
				}
			} else {
				// no process bit for the combi.
				prevF(combi, index, someNew)
			}
		}
	}

	// now generate indexes.
	if forFilters {
		f(combiForFilter, "", false)
	} else {
		for i := 3; i < (1 << uint(len(labels))); i++ {
			if (i & (i - 1)) == 0 {
				// do not save single index
				continue
			}
			f(uint8(i), "", false)
		}
	}

	return indexes, nil
}

func appendCombinationIndex(indexes, index string) string {
	buf := bytes.NewBufferString(indexes)
	if buf.Len() > 0 {
		buf.WriteString(combiIndexSeperator)
	}
	buf.WriteString(index)
	return buf.String()
}
