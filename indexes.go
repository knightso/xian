package xian

import "strings"

// Indexes is extra indexes for datastore query.
type Indexes struct {
	m    indexesMap // key=label, value=indexes
	conf *Config
}

// NewIndexes creates and initializes a new Indexes.
func NewIndexes(config *Config) *Indexes {
	if config == nil {
		config = DefaultConfig
	}
	return &Indexes{
		m:    make(indexesMap),
		conf: config,
	}
}

// Add adds new indexes with a label.
func (idxs *Indexes) Add(label string, indexes ...string) *Indexes {
	for _, idx := range indexes {
		if idxs.conf.IgnoreCase {
			idx = strings.ToLower(idx)
		}

		if _, ok := idxs.m[label]; !ok {
			idxs.m[label] = make(map[string]struct{})
		}

		idxs.m[label][idx] = struct{}{}
	}
	return idxs
}

// AddBigrams adds new bigram indexes with a label.
func (idxs *Indexes) AddBigrams(label string, s string) *Indexes {
	idxs.Add(label, Bigrams(s)...)
	return idxs
}

// AddBiunigrams adds new biunigram indexes with a label.
func (idxs *Indexes) AddBiunigrams(label string, s string) *Indexes {
	idxs.Add(label, Biunigrams(s)...)
	return idxs
}

// AddPrefixes adds new prefix indexes with a label.
func (idxs *Indexes) AddPrefixes(label string, s string) *Indexes {
	idxs.Add(label, Prefixes(s)...)
	return idxs
}

// AddSomething adds new indexes with a label.
// The indexes can be a slice or a string convertible value.
func (idxs *Indexes) AddSomething(label string, indexes interface{}) *Indexes {
	panic("under construction")
}

// Build builds indexes to save.
func (idxs Indexes) Build() []string {

	built := buildIndexes(idxs.m, nil)

	if len(idxs.conf.CompositeIdxLabels) > 1 {
		ci := createCompositeIndexes(idxs.conf.CompositeIdxLabels, idxs.m, false)
		built = append(built, ci...)
	}

	if idxs.conf.SaveNoFiltersIndex {
		built = append(built, IndexNoFilters)
	}

	return built
}
