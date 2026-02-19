package update

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestPull(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "TestPullWithValuesInMap",
			m:     bson.M{"test": "test"},
			key:   "item",
			value: 1,
			want:  bson.M{"test": "test", "$pull": bson.M{"item": 1}},
		},
		{
			name:  "TestPullWithEmptyBSONMap",
			m:     bson.M{},
			key:   "item",
			value: 1,
			want:  bson.M{"$pull": bson.M{"item": 1}},
		},
		{
			name:  "TestPullWithExistingBSONMap",
			m:     bson.M{"$pull": bson.M{"existing": 2}},
			key:   "new",
			value: 3,
			want:  bson.M{"$pull": bson.M{"existing": 2, "new": 3}},
		},
		{
			name:  "TestPullWithSameKeyBSONMap",
			m:     bson.M{"$pull": bson.M{"existing": 4}},
			key:   "existing",
			value: 5,
			want:  bson.M{"$pull": bson.M{"existing": 5}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update(tt.m, pullOperator, tt.key, tt.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Pull() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestPullExported(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "item",
			value: 1,
			want:  bson.M{"$pull": bson.M{"item": 1}},
		},
		{
			name:  "ExistingOperator",
			m:     bson.M{"$pull": bson.M{"existing": 2}},
			key:   "new",
			value: 3,
			want:  bson.M{"$pull": bson.M{"existing": 2, "new": 3}},
		},
		{
			name:  "OverwriteKey",
			m:     bson.M{"$pull": bson.M{"item": 4}},
			key:   "item",
			value: 5,
			want:  bson.M{"$pull": bson.M{"item": 5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Pull(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Pull() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestPush(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "items",
			value: "newItem",
			want:  bson.M{"$push": bson.M{"items": "newItem"}},
		},
		{
			name:  "ExistingOperator",
			m:     bson.M{"$push": bson.M{"existing": "value"}},
			key:   "new",
			value: "anotherValue",
			want:  bson.M{"$push": bson.M{"existing": "value", "new": "anotherValue"}},
		},
		{
			name:  "OverwriteKey",
			m:     bson.M{"$push": bson.M{"items": "old"}},
			key:   "items",
			value: "new",
			want:  bson.M{"$push": bson.M{"items": "new"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Push(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Push() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "field",
			value: "value",
			want:  bson.M{"$set": bson.M{"field": "value"}},
		},
		{
			name:  "WithExistingSet",
			m:     bson.M{"$set": bson.M{"other": "val"}},
			key:   "field",
			value: 42,
			want:  bson.M{"$set": bson.M{"other": "val", "field": 42}},
		},
		{
			name:  "OverwriteKey",
			m:     bson.M{"$set": bson.M{"field": "old"}},
			key:   "field",
			value: "new",
			want:  bson.M{"$set": bson.M{"field": "new"}},
		},
		{
			name:  "IntValue",
			m:     bson.M{},
			key:   "count",
			value: 99,
			want:  bson.M{"$set": bson.M{"count": 99}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Set() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestInc(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "counter",
			value: 1,
			want:  bson.M{"$inc": bson.M{"counter": 1}},
		},
		{
			name:  "WithExistingInc",
			m:     bson.M{"$inc": bson.M{"views": 5}},
			key:   "clicks",
			value: 3,
			want:  bson.M{"$inc": bson.M{"views": 5, "clicks": 3}},
		},
		{
			name:  "NegativeIncrement",
			m:     bson.M{},
			key:   "balance",
			value: -10,
			want:  bson.M{"$inc": bson.M{"balance": -10}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Inc(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("Inc() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestAddToSet(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "tags",
			value: "go",
			want:  bson.M{"$addToSet": bson.M{"tags": "go"}},
		},
		{
			name:  "WithExistingAddToSet",
			m:     bson.M{"$addToSet": bson.M{"colors": "red"}},
			key:   "sizes",
			value: "M",
			want:  bson.M{"$addToSet": bson.M{"colors": "red", "sizes": "M"}},
		},
		{
			name:  "OverwriteKey",
			m:     bson.M{"$addToSet": bson.M{"tags": "old"}},
			key:   "tags",
			value: "new",
			want:  bson.M{"$addToSet": bson.M{"tags": "new"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddToSet(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("AddToSet() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestSetOnInsert(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value any
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			m:     bson.M{},
			key:   "createdAt",
			value: "2024-01-01",
			want:  bson.M{"$setOnInsert": bson.M{"createdAt": "2024-01-01"}},
		},
		{
			name:  "WithExistingSetOnInsert",
			m:     bson.M{"$setOnInsert": bson.M{"field1": "val1"}},
			key:   "field2",
			value: "val2",
			want:  bson.M{"$setOnInsert": bson.M{"field1": "val1", "field2": "val2"}},
		},
		{
			name:  "IntValue",
			m:     bson.M{},
			key:   "version",
			value: 1,
			want:  bson.M{"$setOnInsert": bson.M{"version": 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetOnInsert(tt.key, tt.value)(tt.m)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("SetOnInsert() = %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want bson.M
	}{
		{
			name: "NoOptions",
			opts: []Option{},
			want: bson.M{},
		},
		{
			name: "SingleSet",
			opts: []Option{Set("field", "value")},
			want: bson.M{"$set": bson.M{"field": "value"}},
		},
		{
			name: "MultipleOperators",
			opts: []Option{
				Set("name", "Alice"),
				Inc("counter", 1),
				Push("items", "newItem"),
			},
			want: bson.M{
				"$set":  bson.M{"name": "Alice"},
				"$inc":  bson.M{"counter": 1},
				"$push": bson.M{"items": "newItem"},
			},
		},
		{
			name: "AllOperators",
			opts: []Option{
				Set("a", 1),
				Push("b", 2),
				Pull("c", 3),
				Inc("d", 4),
				AddToSet("e", 5),
				SetOnInsert("f", 6),
			},
			want: bson.M{
				"$set":         bson.M{"a": 1},
				"$push":        bson.M{"b": 2},
				"$pull":        bson.M{"c": 3},
				"$inc":         bson.M{"d": 4},
				"$addToSet":    bson.M{"e": 5},
				"$setOnInsert": bson.M{"f": 6},
			},
		},
		{
			name: "SameOperatorMultipleTimes",
			opts: []Option{
				Set("a", 1),
				Set("b", 2),
			},
			want: bson.M{
				"$set": bson.M{"a": 1, "b": 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Generate(tt.opts...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
