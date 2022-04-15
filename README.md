# xian

Package `github.com/knightso/xian` generates Indexes and Filters to search NoSQL Datastores.

It's designed espesially for Google Cloud Datastore or Cloud Firestore in Datastore mode. However it doesn't depend on any specific Datastore APIs so that you may be able to use for some other NoSQL databases which support list-property and merge join.

## Features

* prefix/suffix/partial match search
* IN search
* reduce composite indexes(esp. for Cloud Datastore)

## Note

* Search latency can increase depending on its result-set size and filter condition.
* Index storage size can be bigger especially with long text prefix/suffix/partial match.

## Usage

Code example is for Cloud Datastore.

### Installation

```
$ go get -u github.com/knightso/xian
```

### Configuration

```go
var bookIndexesConfig = xian.MustValidateConfig(&xian.Config{
	IgnoreCase:         true, // search case-insensitive
	SaveNoFiltersIndex: true, // always save 'NoFilters' index
})

// configure IN-filter
var statusInBuilder *xian.InBuilder = xian.NewInBuilder()

var (
    BookStatusUnpublished = statusInBuilder.NewBit()
    BookStatusPublished = statusInBuilder.NewBit()
    BookStatusDiscontinued = statusInBuilder.NewBit()
)
```

This configuration should be used to initialize both Indexes and Filters.

### Label Constants

Define common labels for both Indexes and Filters.  
Constants are not necessaly but recommended.  
Short label names would make index size smaller.

```go
const (
	BookQueryLabelTitlePartial = "ti"
	BookQueryLabelTitlePrefix = "tp"
	BookQueryLabelTitleSuffix = "ts"
	BookQueryLabelIsHobby = "h"
	BookQueryLabelStatusIN = "s"
	BookQueryLabelPriceRange = "pr"
)
```

### Save indexes

```go
idxs := xian.NewIndexes(bookIndexesConfig)

idxs.AddBigrams(BookQueryLabelTitlePartial, book.Title)
idxs.AddBiunigrams(BookQueryLabelTitlePartial, book.Title)
idxs.AddPrefixes(BookQueryLabelTitlePrefix, book.Title)
idxs.AddSuffixes(BookQueryLabelTitleSuffix, book.Title)
idxs.AddSomething(BookQueryLabelIsHobby, book.Category == "sports" || book.Category == "cooking")
idxs.Add(BookQueryLabelStatusIN, statusInBuilder.Indexes(BookStatusUnpublished)...)

switch {
case book.Price < 3000:
	idxs.Add(BookQueryLabelPriceRange, "p<3000")
case book.Price < 5000:
	idxs.Add(BookQueryLabelPriceRange, "3000<=p<5000")
case book.Price < 10000:
	idxs.Add(BookQueryLabelPriceRange, "5000<=p<10000")
default:
	idxs.Add(BookQueryLabelPriceRange, "10000<=p")
}

// build and set indexes to the book's property
var err error
book.Indexes, err = idxs.Build()
if err != nil {
	return err
}

// save book
```

### Search (example for Cloud Datastore)

```go
q := datastore.NewQuery("Book")

filters := xian.NewFilters(bookIndexesConfig).
    AddSomething(BookQueryLabelIsHolly, true).
    Add(BookQueryLabelStatusIN, statusInBuilder.Filter(BookStatusUnpublished, BookStatusPublished)).
    Add(BookQueryLabelPriceRange, "5000<=p<10000").
    AddBigrams(BookQueryLabelTitlePartial, title).
    AddBiunigrams(BookQueryLabelTitlePartial, title).
    AddSuffix(BookQueryLabelTitleSuffix, title)


built, err := filters.Build()
if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
}

for _, f := range built {
    q = q.Filter("Indexes =", f)
}

// query books
```
