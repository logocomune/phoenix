package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DeleteAll removes all documents from collectionName in db that match filterQuery.
// A context deadline is applied using timeout. It returns the number of deleted
// documents and any error encountered.
func DeleteAll(ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.DeleteOptions) (int64, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()

	deletedResp, err := db.Collection(collectionName).DeleteMany(ctx, filterQuery, opts...)
	if err != nil {
		return 0, err
	}
	if deletedResp == nil {
		return 0, err
	}

	return deletedResp.DeletedCount, nil
}
