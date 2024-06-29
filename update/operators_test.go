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
		value interface{}
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
