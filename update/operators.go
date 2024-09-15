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

type Option func(m bson.M)

func initializeMap(m bson.M, operator string) bson.M {
	if valMap, ok := m[operator]; ok {
		return valMap.(bson.M)
	}
	return bson.M{}
}

func Pull(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, pullOperator, key, value)
	}
}

func Push(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, pushOperator, key, value)
	}
}

func Set(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, setOperator, key, value)
	}
}

func Inc(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, incOperator, key, value)
	}
}

func AddToSet(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, addToSet, key, value)
	}
}

func SetOnInsert(key string, value interface{}) Option {
	return func(m bson.M) {
		update(m, setOnInsert, key, value)
	}
}

func update(m bson.M, operator string, key string, value interface{}) {
	v := initializeMap(m, operator)
	v[key] = value
	m[operator] = v
}

func Generate(opt ...Option) bson.M {
	m := bson.M{}
	for i := range opt {
		opt[i](m)
	}
	return m
}
