package phoenix

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

// FindOne retrieves a single document from collectionName in db that matches filterQuery.
// A context deadline is applied using timeout. It returns the decoded document, a boolean
// indicating whether a matching document was found, and any error. If no document matches,
// it returns (zero value, false, nil) without an error.
func FindOne[D any](ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.FindOneOptions) (D, bool, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()

	var result D
	err := db.Collection(collectionName).FindOne(ctx, filterQuery, opts...).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, false, nil
		}
		return result, false, err
	}

	return result, true, nil
}

// FindAll retrieves all documents from collectionName in db that match filterQuery.
// A context deadline is applied using timeout. It returns a slice of decoded documents
// and any error. When no documents are found, a nil slice is returned with a nil error.
func FindAll[D any](ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.FindOptions) ([]D, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()

	cur, err := db.Collection(collectionName).Find(ctx, filterQuery, opts...)
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close is best-effort; the error is not returned to the caller.
		if err := cur.Close(ctx); err != nil {
			slog.Error("cursor close error", slog.String("error", err.Error()))
		}
	}()

	var documents []D
	for cur.Next(ctx) {
		var d D
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		documents = append(documents, d)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return documents, nil
}
