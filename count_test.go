package phoenix

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("with documents", func(mt *mtest.T) {
		// CountDocuments uses an aggregate pipeline internally;
		// the result cursor contains a document with field "n".
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch,
			bson.D{{Key: "n", Value: int32(5)}},
		))

		count, err := Count(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 5 {
			t.Errorf("expected count 5, got %d", count)
		}
	})

	mt.Run("zero count", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))

		count, err := Count(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 0 {
			t.Errorf("expected count 0, got %d", count)
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, err := Count(context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
