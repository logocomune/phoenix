package query

import "go.mongodb.org/mongo-driver/bson"

type Option func(m bson.M)

// ByKeyValue sets the value of a given key in a bson.M map to the specified value.
func ByKeyValue(m bson.M, key string, value interface{}) {
	m[key] = value
}

// InKeyValue sets the value of a given key in a bson.M map to a map with the key "$in" and the specified value.
func InKeyValue(m bson.M, key string, value interface{}) {
	m[key] = bson.M{
		"$in": value,
	}
}

// setKeyValue sets the value of a given key in a bson.M map to the specified value.
func setKeyValue(m bson.M, key string, op string, value interface{}) {
	var q bson.M
	var ok bool

	if q, ok = m[key].(bson.M); !ok {
		q = bson.M{}
	}

	q[op] = value
	m[key] = q
}

// GTEKeyValue sets the value of a given key in a bson.M map to the specified value using the "$gte" operation.
func GTEKeyValue(m bson.M, key string, value interface{}) {
	setKeyValue(m, key, "$gte", value)
}

// GTKeyValue sets the value of a given key in a bson.M map to the specified value
func GTKeyValue(m bson.M, key string, value interface{}) {
	setKeyValue(m, key, "$gt", value)
}

// LTEKeyValue sets the value of a given key in a bson.M map to the specified value using the "$lte" operation.
func LTEKeyValue(m bson.M, key string, value interface{}) {
	setKeyValue(m, key, "$lte", value)
}

// LTKeyValue sets the value of a given key in a bson.M map to the specified value with the
func LTKeyValue(m bson.M, key string, value interface{}) {
	setKeyValue(m, key, "$lt", value)
}

// Generate creates a bson.M map by applying a series of Option functions.
func Generate(opt ...Option) bson.M {
	m := bson.M{}
	for i := range opt {
		opt[i](m)
	}
	return m
}
