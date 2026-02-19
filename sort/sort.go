// Package sort provides builder functions for constructing MongoDB sort documents
// using the functional options pattern. Use WithSort or WithSorts to declare sort
// fields, then pass the resulting Options to Generate.
package sort

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Option is a function that appends sort entries to a bson.D document.
type Option func(d bson.D) bson.D

// WithSort returns an Option that sorts by the single field described by e.
// Use 1 for ascending and -1 for descending order.
func WithSort(e bson.E) Option {
	return WithSorts(e)
}

// WithSorts returns an Option that appends one or more sort fields to the document.
// Fields are applied in the order provided; earlier entries take priority.
func WithSorts(sort ...bson.E) Option {
	return func(d bson.D) bson.D {
		return append(d, sort...)
	}
}

// Generate composes multiple Options into a single bson.D sort document.
func Generate(opt ...Option) bson.D {
	d := bson.D{}
	for i := range opt {
		d = opt[i](d)
	}
	return d
}
