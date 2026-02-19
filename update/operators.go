// Package update provides builder functions for constructing MongoDB update documents
// using the functional options pattern. Each operator function returns an Option that
// can be composed with Generate to produce a complete update bson.M.
package update

import (
	"go.mongodb.org/mongo-driver/bson"
)

const (
	pushOperator = "$push"
	pullOperator = "$pull"
	setOperator  = "$set"
	incOperator  = "$inc"
	addToSet     = "$addToSet"
	setOnInsert  = "$setOnInsert"
)

// Option is a function that modifies a bson.M update document.
type Option func(m bson.M)

// initializeMap returns the existing operator sub-document from m, or a new
// empty bson.M if the operator has not been set yet.
func initializeMap(m bson.M, operator string) bson.M {
	if valMap, ok := m[operator]; ok {
		return valMap.(bson.M)
	}
	return bson.M{}
}

// Pull returns an Option that appends a $pull operator entry to the update document,
// removing from the array field identified by key all elements equal to value.
func Pull(key string, value any) Option {
	return func(m bson.M) {
		update(m, pullOperator, key, value)
	}
}

// Push returns an Option that appends a $push operator entry to the update document,
// adding value to the array field identified by key.
func Push(key string, value any) Option {
	return func(m bson.M) {
		update(m, pushOperator, key, value)
	}
}

// Set returns an Option that appends a $set operator entry to the update document,
// setting the field identified by key to value.
func Set(key string, value any) Option {
	return func(m bson.M) {
		update(m, setOperator, key, value)
	}
}

// Inc returns an Option that appends an $inc operator entry to the update document,
// incrementing the field identified by key by value.
func Inc(key string, value any) Option {
	return func(m bson.M) {
		update(m, incOperator, key, value)
	}
}

// AddToSet returns an Option that appends an $addToSet operator entry to the update
// document, adding value to the set field identified by key only if it is not already
// present.
func AddToSet(key string, value any) Option {
	return func(m bson.M) {
		update(m, addToSet, key, value)
	}
}

// SetOnInsert returns an Option that appends a $setOnInsert operator entry to the
// update document. The field identified by key is set to value only when an upsert
// results in a new document being inserted.
func SetOnInsert(key string, value any) Option {
	return func(m bson.M) {
		update(m, setOnInsert, key, value)
	}
}

// update writes key=value into the sub-document for operator inside m.
func update(m bson.M, operator string, key string, value any) {
	v := initializeMap(m, operator)
	v[key] = value
	m[operator] = v
}

// Generate composes multiple Options into a single bson.M update document.
// Options are applied in order; later Options for the same operator and key
// overwrite earlier ones.
func Generate(opt ...Option) bson.M {
	m := bson.M{}
	for i := range opt {
		opt[i](m)
	}
	return m
}
