package sort

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestWithSorts(t *testing.T) {
	tests := []struct {
		name string
		sort []bson.E
		want bson.D
	}{
		{
			name: "EmptySort",
			sort: []bson.E{},
			want: bson.D{},
		},
		{
			name: "SingleSort",
			sort: []bson.E{bson.E{Key: "key1", Value: 1}},
			want: bson.D{bson.E{Key: "key1", Value: 1}},
		},
		{
			name: "MultipleSort",
			sort: []bson.E{bson.E{Key: "key1", Value: 1}, bson.E{Key: "key2", Value: -1}},
			want: bson.D{bson.E{Key: "key1", Value: 1}, bson.E{Key: "key2", Value: -1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithSorts(tt.sort...)(bson.D{}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithSorts() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		opt  []Option
		want bson.D
	}{
		{
			name: "one option",
			opt:  []Option{WithSorts(bson.E{Key: "name", Value: -1})},
			want: bson.D{{"name", -1}},
		},
		{
			name: "multiple options",
			opt:  []Option{WithSorts(bson.E{Key: "name", Value: -1}), WithSorts(bson.E{Key: "age", Value: 1})},
			want: bson.D{{"name", -1}, {"age", 1}},
		},
		{
			name: "no options",
			opt:  []Option{},
			want: bson.D{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.opt...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
