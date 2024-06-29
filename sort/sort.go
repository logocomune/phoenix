package sort

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Option func(d bson.D) bson.D

// WithSort is a function that creates an Option to sort a given document by the provided field.
func WithSort(e bson.E) Option {
	return WithSorts(e)
}

// WithSorts is a function that creates an Option to append sorting options to a given document.
func WithSorts(sort ...bson.E) Option {
	return func(d bson.D) bson.D {

		return append(d, sort...)

	}
}

// Generate is a function that creates a bson.D document by applying a series of Option functions.
func Generate(opt ...Option) bson.D {
	d := bson.D{}
	for i := range opt {
		d = opt[i](d)
	}
	return d
}
