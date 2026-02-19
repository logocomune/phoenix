package phoenix

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUpdateMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(3)},
			{Key: "nModified", Value: int32(2)},
		})

		matched, modified, err := UpdateMany(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 3 {
			t.Errorf("expected matched=3, got %d", matched)
		}
		if modified != 2 {
			t.Errorf("expected modified=2, got %d", modified)
		}
	})

	mt.Run("no match", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(0)},
			{Key: "nModified", Value: int32(0)},
		})

		matched, modified, err := UpdateMany(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 0 || modified != 0 {
			t.Errorf("expected matched=0 modified=0, got matched=%d modified=%d", matched, modified)
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, _, err := UpdateMany(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}

func TestUpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(1)},
			{Key: "nModified", Value: int32(1)},
		})

		matched, modified, err := UpdateOne(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 1 {
			t.Errorf("expected matched=1, got %d", matched)
		}
		if modified != 1 {
			t.Errorf("expected modified=1, got %d", modified)
		}
	})

	mt.Run("no match", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(0)},
			{Key: "nModified", Value: int32(0)},
		})

		matched, modified, err := UpdateOne(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 0 || modified != 0 {
			t.Errorf("expected matched=0 modified=0, got matched=%d modified=%d", matched, modified)
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, _, err := UpdateOne(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{}, bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
