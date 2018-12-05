package xian

import (
	"unicode/utf8"
)

// Filters is filters builder for extra indexes.
type Filters struct {
	m    map[string][]string // key=label, value=indexes
	conf *Config
}

// NewFilters creates and initializes a new Filters.
func NewFilters(config *Config) *Filters {
	if config == nil {
		config = DefaultConfig
	}
	return &Filters{
		m:    make(map[string][]string),
		conf: config,
	}
}

// Add adds new filters with a label.
func (filters *Filters) Add(label string, indexes ...string) *Filters {
	filters.m[label] = append(filters.m[label], indexes...)
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

// AddSomething adds new indexes with a label.
// The indexes can be a slice or a string convertible value.
func (filters *Filters) AddSomething(label string, indexes interface{}) *Filters {
	panic("under construction")
}

// Build builds indexes to save.
func (filters *Filters) Build() []string {

	built := buildIndexes(filters.m)

	if filters.conf.SaveNoFilterIndex && len(built) == 0 {
		built = append(built, IndexNoFilters)
	}

	return built
}
