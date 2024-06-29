package query

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestByKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value interface{}
		want  bson.M
	}{
		{
			name:  "KeyValueTest",
			m:     bson.M{},
			key:   "key1",
			value: "value1",
			want:  bson.M{"key1": "value1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ByKeyValue(tt.m, tt.key, tt.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Fatalf("expected: %v, got: %v", tt.want, tt.m)
			}
		})
	}
}

func TestInKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value interface{}
		want  bson.M
	}{
		{
			name:  "InKeyValueTest",
			m:     bson.M{},
			key:   "key2",
			value: []string{"value2.1", "value2.2"},
			want:  bson.M{"key2": bson.M{"$in": []string{"value2.1", "value2.2"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InKeyValue(tt.m, tt.key, tt.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Fatalf("expected: %v, got: %v", tt.want, tt.m)
			}
		})
	}
}

//Similar test functions for GTEKeyValue, GTKeyValue, LTEKeyValue, LTKeyValue

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string
		opt  []Option
		want bson.M
	}{
		{
			name: "GenerateTest",
			opt: []Option{
				func(m bson.M) {
					ByKeyValue(m, "key1", "value1")
				},
				func(m bson.M) {
					InKeyValue(m, "key2", []string{"value2.1", "value2.2"})
				},
				//Add more Options here
			},
			want: bson.M{
				"key1": "value1",
				"key2": bson.M{"$in": []string{"value2.1", "value2.2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.opt...); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGTEKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value interface{}
		wantM bson.M
	}{
		{
			name:  "Existing key",
			m:     bson.M{"age": bson.M{"$gt": 18}},
			key:   "age",
			value: 30,
			wantM: bson.M{"age": bson.M{"$gt": 18, "$gte": 30}},
		},
		{
			name:  "New key",
			m:     bson.M{"age": 18},
			key:   "height",
			value: 150,
			wantM: bson.M{"age": 18, "height": bson.M{"$gte": 150}},
		},
		{
			name:  "Empty key",
			m:     bson.M{"age": 18},
			key:   "",
			value: 150,
			wantM: bson.M{"age": 18, "": bson.M{"$gte": 150}},
		},
		{
			name:  "Empty map",
			m:     bson.M{},
			key:   "age",
			value: 18,
			wantM: bson.M{"age": bson.M{"$gte": 18}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GTEKeyValue(tt.m, tt.key, tt.value)
			if !reflect.DeepEqual(tt.wantM, tt.m) {
				t.Errorf("GTKeyValue() = %v, want %v", tt.wantM, tt.m)
			}
		})
	}
}
func TestGTKeyValue(t *testing.T) {
	type args struct {
		m     bson.M
		key   string
		value interface{}
	}

	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "EmptyMap",
			args: args{
				m:     bson.M{},
				key:   "test",
				value: 123,
			},
			want: bson.M{
				"test": bson.M{
					"$gt": 123,
				},
			},
		},
		{
			name: "ExistingKey",
			args: args{
				m: bson.M{
					"test": bson.M{
						"$lt": 0,
					},
				},
				key:   "test",
				value: 123,
			},
			want: bson.M{
				"test": bson.M{
					"$lt": 0,
					"$gt": 123,
				},
			},
		},
		{
			name: "NonMapValue",
			args: args{
				m:     bson.M{"test": 123},
				key:   "test",
				value: 456,
			},
			want: bson.M{
				"test": bson.M{
					"$gt": 456,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GTKeyValue(tt.args.m, tt.args.key, tt.args.value)

			if !reflect.DeepEqual(tt.args.m, tt.want) {
				t.Errorf("GTKeyValue() = %v, want %v", tt.args.m, tt.want)
			}
		})
	}
}
func TestLTEKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		m     bson.M
		key   string
		value interface{}
		want  bson.M
	}{
		{
			name:  "IntLTE",
			m:     bson.M{},
			key:   "val",
			value: 9,
			want:  bson.M{"val": bson.M{"$lte": 9}},
		},
		{
			name:  "StringLTE",
			m:     bson.M{},
			key:   "val",
			value: "abc",
			want:  bson.M{"val": bson.M{"$lte": "abc"}},
		},
		{
			name:  "ExistingKey",
			m:     bson.M{"val": bson.M{"$gt": 5}},
			key:   "val",
			value: 9,
			want:  bson.M{"val": bson.M{"$gt": 5, "$lte": 9}},
		},
		{
			name:  "FloatLTE",
			m:     bson.M{},
			key:   "val",
			value: 0.34,
			want:  bson.M{"val": bson.M{"$lte": 0.34}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LTEKeyValue(tt.m, tt.key, tt.value)
			if !reflect.DeepEqual(tt.m, tt.want) {
				t.Errorf("got %v, want %v", tt.m, tt.want)
			}
		})
	}
}

func TestLTKeyValue(t *testing.T) {
	tests := []struct {
		name  string
		input bson.M
		key   string
		value interface{}
		want  bson.M
	}{
		{
			name:  "EmptyMap",
			input: bson.M{},
			key:   "Key1",
			value: "Value1",
			want:  bson.M{"Key1": bson.M{"$lt": "Value1"}},
		},
		{
			name:  "ExistingKey",
			input: bson.M{"Key1": bson.M{"$gt": "Value0"}},
			key:   "Key1",
			value: "Value1",
			want:  bson.M{"Key1": bson.M{"$gt": "Value0", "$lt": "Value1"}},
		},
		{
			name:  "NonExistingKey",
			input: bson.M{"Key2": bson.M{"$gt": "Value0"}},
			key:   "Key1",
			value: "Value1",
			want:  bson.M{"Key1": bson.M{"$lt": "Value1"}, "Key2": bson.M{"$gt": "Value0"}},
		},
		{
			name:  "NilValue",
			input: bson.M{},
			key:   "Key1",
			value: nil,
			want:  bson.M{"Key1": bson.M{"$lt": nil}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LTKeyValue(tt.input, tt.key, tt.value)
			if !reflect.DeepEqual(tt.input, tt.want) {
				t.Errorf("LTKeyValue() got %v, want %v", tt.input, tt.want)
			}
		})
	}
}
