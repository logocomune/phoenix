// Package query provides builder functions for constructing MongoDB filter documents
// using the functional options pattern. Functions like ByKeyValue, InKeyValue, GTKeyValue,
// etc. can be composed with Generate to produce a complete filter bson.M.
package query

import "go.mongodb.org/mongo-driver/bson"

// Option is a function that modifies a bson.M filter document.
type Option func(m bson.M)

// ByKeyValue sets an exact-match condition for key in m, requiring the field
// to equal value.
func ByKeyValue(m bson.M, key string, value any) {
	m[key] = value
}

// InKeyValue sets an $in condition for key in m, matching documents where the
// field value is contained in value (typically a slice).
func InKeyValue(m bson.M, key string, value any) {
	m[key] = bson.M{
		"$in": value,
	}
}

// setKeyValue merges a comparison operator op and value into the existing
// sub-document for key in m, creating it if it does not exist.
func setKeyValue(m bson.M, key string, op string, value any) {
	var q bson.M
	var ok bool

	if q, ok = m[key].(bson.M); !ok {
		q = bson.M{}
	}

	q[op] = value
	m[key] = q
}

// GTEKeyValue sets a $gte (greater-than-or-equal) condition for key in m.
func GTEKeyValue(m bson.M, key string, value any) {
	setKeyValue(m, key, "$gte", value)
}

// GTKeyValue sets a $gt (greater-than) condition for key in m.
func GTKeyValue(m bson.M, key string, value any) {
	setKeyValue(m, key, "$gt", value)
}

// LTEKeyValue sets a $lte (less-than-or-equal) condition for key in m.
func LTEKeyValue(m bson.M, key string, value any) {
	setKeyValue(m, key, "$lte", value)
}

// LTKeyValue sets a $lt (less-than) condition for key in m.
func LTKeyValue(m bson.M, key string, value any) {
	setKeyValue(m, key, "$lt", value)
}

// Generate composes multiple Options into a single bson.M filter document.
// Options are applied in order and may freely combine multiple conditions on
// the same or different fields.
func Generate(opt ...Option) bson.M {
	m := bson.M{}
	for i := range opt {
		opt[i](m)
	}
	return m
}
