package phoenix

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNew(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("returns non-nil Phoenix", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		if p == nil {
			t.Fatal("expected non-nil Phoenix")
		}
	})
}

func TestPhoenix_FindOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch, bson.D{
			{Key: "name", Value: "Alice"},
			{Key: "age", Value: int32(30)},
		}))

		doc, found, err := p.FindOne(context.Background(), bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !found {
			t.Fatal("expected document to be found")
		}
		if doc.Name != "Alice" {
			t.Errorf("expected Alice, got %s", doc.Name)
		}
	})

	mt.Run("not found", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch))

		_, found, err := p.FindOne(context.Background(), bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if found {
			t.Fatal("expected document not to be found")
		}
	})
}

func TestPhoenix_FindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch,
			bson.D{{Key: "name", Value: "Alice"}, {Key: "age", Value: int32(30)}},
			bson.D{{Key: "name", Value: "Bob"}, {Key: "age", Value: int32(25)}},
		))

		docs, err := p.FindAll(context.Background(), bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(docs) != 2 {
			t.Fatalf("expected 2 documents, got %d", len(docs))
		}
	})
}

func TestPhoenix_Count(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		ns := fmt.Sprintf("%s.testcol", mt.DB.Name())
		mt.AddMockResponses(mtest.CreateCursorResponse(0, ns, mtest.FirstBatch,
			bson.D{{Key: "n", Value: int32(7)}},
		))

		count, err := p.Count(context.Background(), bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 7 {
			t.Errorf("expected count 7, got %d", count)
		}
	})
}

func TestPhoenix_DeleteAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(2)},
		})

		count, err := p.DeleteAll(context.Background(), bson.M{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 deleted, got %d", count)
		}
	})
}

func TestPhoenix_UpdateMany(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(4)},
			{Key: "nModified", Value: int32(3)},
		})

		matched, modified, err := p.UpdateMany(context.Background(), bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 4 || modified != 3 {
			t.Errorf("expected matched=4 modified=3, got matched=%d modified=%d", matched, modified)
		}
	})
}

func TestPhoenix_UpdateOne(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		p := New[testDoc](mt.DB, "testcol", 5*time.Second)
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: int32(1)},
			{Key: "nModified", Value: int32(1)},
		})

		matched, modified, err := p.UpdateOne(context.Background(), bson.M{}, bson.M{"$set": bson.M{"status": "active"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if matched != 1 || modified != 1 {
			t.Errorf("expected matched=1 modified=1, got matched=%d modified=%d", matched, modified)
		}
	})
}
