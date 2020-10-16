package xian

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Filters is filters builder for extra indexes.
type Filters struct {
	m    indexesMap // key=label, value=index set
	conf *Config
}

// NewFilters creates and initializes a new Filters.
func NewFilters(conf *Config) *Filters {
	if conf == nil {
		conf = DefaultConfig
	}
	return &Filters{
		m:    make(indexesMap),
		conf: conf,
	}
}

// Add adds new filters with a label.
func (filters *Filters) Add(label string, indexes ...string) *Filters {
	for _, idx := range indexes {
		if filters.conf.IgnoreCase {
			idx = strings.ToLower(idx)
		}

		if _, ok := filters.m[label]; !ok {
			filters.m[label] = make(map[string]struct{})
		}

		filters.m[label][idx] = struct{}{}
	}
	return filters
}

// AddBigrams adds new bigram filters with a label.
func (filters *Filters) AddBigrams(label string, s string) *Filters {
	// same filter as biunigrams'
	filters.AddBiunigrams(label, s)
	return filters
}

// AddBiunigrams adds new biunigram filters with a label.
func (filters *Filters) AddBiunigrams(label string, s string) *Filters {
	if runeLen := utf8.RuneCountInString(s); runeLen == 1 {
		filters.Add(label, s)
	} else if runeLen > 1 {
		filters.Add(label, Bigrams(s)...)
	}
	return filters
}

// AddPrefix adds a new prefix filter with a label.
func (filters *Filters) AddPrefix(label string, s string) *Filters {
	// don't need to split prefixes on filters
	return filters.Add(label, s)
}

// AddSuffix adds a new suffix filter with a label.
func (filters *Filters) AddSuffix(label string, s string) *Filters {
	// don't need to split suffixes on filters
	return filters.Add(label, s)
}

// AddSomething adds new indexes with a label.
// The indexes can be a slice or a string convertible value.
func (filters *Filters) AddSomething(label string, indexes interface{}) *Filters {

	v := reflect.ValueOf(indexes)

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			filters.Add(label, fmt.Sprintf("%v", v.Index(i).Interface()))
		}
	case timeKind:
		unix := v.Interface().(time.Time).UnixNano()
		filters.Add(label, fmt.Sprintf("%d", unix))
	default:
		filters.Add(label, fmt.Sprintf("%v", v.Interface()))
	}

	return filters
}

// Build builds filters to save.
func (filters *Filters) Build() ([]string, error) {

	built := buildIndexes(filters.m, filters.conf.CompositeIdxLabels)

	if len(filters.conf.CompositeIdxLabels) > 1 {
		cis, err := createCompositeIndexes(filters.conf.CompositeIdxLabels, filters.m, true)
		if err != nil {
			return nil, err
		}
		built = append(built, cis...)
	}

	if filters.conf.SaveNoFiltersIndex && len(built) == 0 {
		built = append(built, IndexNoFilters)
	}

	if len(built) > MaxIndexesSize {
		return nil, errors.Errorf("index size exceeds %d", MaxIndexesSize)
	}

	return built, nil
}

// MustBuild builds filters to save and panics with error.
func (filters Filters) MustBuild() []string {
	built, err := filters.Build()
	if err != nil {
		panic(err)
	}
	return built
}
