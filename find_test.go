package phoenix

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// testDoc is a shared document type used across root-level test files.
type testDoc struct {
	Name string `bson:"name"`
	Age  int32  `bson:"age"`
}

func TestFindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, bson.D{
			{Key: "name", Value: "Alice"},
			{Key: "age", Value: int32(30)},
		}))

		doc, found, err := FindOne[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !found {
			t.Fatal("expected document to be found")
		}
		if doc.Name != "Alice" || doc.Age != 30 {
			t.Errorf("unexpected document: %+v", doc)
		}
	})

	mt.Run("not found", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))

		_, found, err := FindOne[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found {
			t.Fatal("expected document not to be found")
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, found, err := FindOne[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if found {
			t.Fatal("expected not found on error")
		}
	})
}

func TestFindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("multiple documents", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch,
			bson.D{{Key: "name", Value: "Alice"}, {Key: "age", Value: int32(30)}},
			bson.D{{Key: "name", Value: "Bob"}, {Key: "age", Value: int32(25)}},
		))

		docs, err := FindAll[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(docs) != 2 {
			t.Fatalf("expected 2 documents, got %d", len(docs))
		}
		if docs[0].Name != "Alice" || docs[1].Name != "Bob" {
			t.Errorf("unexpected documents: %+v", docs)
		}
	})

	mt.Run("single document", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch,
			bson.D{{Key: "name", Value: "Charlie"}, {Key: "age", Value: int32(35)}},
		))

		docs, err := FindAll[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(docs))
		}
		if docs[0].Name != "Charlie" {
			t.Errorf("expected Charlie, got %s", docs[0].Name)
		}
	})

	mt.Run("empty result", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))

		docs, err := FindAll[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(docs) != 0 {
			t.Fatalf("expected 0 documents, got %d", len(docs))
		}
	})

	mt.Run("command error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "command failed",
			Name:    "InternalError",
		}))

		_, err := FindAll[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	// cursor iteration error: first batch succeeds but getMore fails.
	// cur.Err() must be propagated after the loop — this is the fix for the missing check.
	mt.Run("cursor iteration error", func(mt *mtest.T) {
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		// cursorID=1 tells the driver there is a second batch to fetch.
		mt.AddMockResponses(mtest.CreateCursorResponse(1, ns, mtest.FirstBatch,
			bson.D{{Key: "name", Value: "Alice"}, {Key: "age", Value: int32(30)}},
		))
		// The getMore request returns a command error.
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "network error",
			Name:    "NetworkError",
		}))

		docs, err := FindAll[testDoc](context.Background(), 5*time.Second, mt.DB, "testcol", bson.M{})
		if err == nil {
			t.Fatalf("expected cursor error to be propagated, got nil (docs: %v)", docs)
		}
		if docs != nil {
			t.Errorf("expected nil docs on error, got %v", docs)
		}
	})
}
