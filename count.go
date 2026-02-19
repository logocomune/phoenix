package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Count returns the number of documents in collectionName in db that match filterQuery.
// A context deadline is applied using timeout.
func Count(ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.CountOptions) (int64, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()

	countDoc, err := db.Collection(collectionName).CountDocuments(ctx, filterQuery, opts...)
	if err != nil {
		return 0, err
	}
	return countDoc, nil
}
