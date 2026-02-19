package phoenix

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestDeleteAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(3)},
		})

		count, err := DeleteAll(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 3 {
			t.Errorf("expected 3 deleted, got %d", count)
		}
	})

	mt.Run("zero deleted", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(0)},
		})

		count, err := DeleteAll(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 deleted, got %d", count)
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, err := DeleteAll(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
