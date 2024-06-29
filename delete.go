package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

// DeleteAll deletes all documents that match the given filter query in the specified
// collection of the provided MongoDB database. It uses the provided context and
// timeout for the operation. Additional delete options can be passed as variadic
// arguments. It returns the number of documents deleted and an error, if any.
func DeleteAll(ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.DeleteOptions) (int64, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()
	deletedResp, err := db.Collection(collectionName).DeleteMany(ctx, filterQuery, opts...)
	if err != nil {
		slog.Error("Error deleting data",
			slog.String("collection", collectionName),
			slog.String("error", err.Error()))
		return 0, err
	}
	if deletedResp == nil {
		return 0, err
	}

	return deletedResp.DeletedCount, nil

}
