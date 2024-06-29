package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func UpdateMany(ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()
	updateResponse, err := db.Collection(collectionName).UpdateMany(ctx, filter, update, opts...)
	if err != nil || updateResponse == nil {
		return 0, 0, err
	}

	return updateResponse.MatchedCount, updateResponse.ModifiedCount, err

}

func UpdateOne(ctx context.Context, timeout time.Duration, db *mongo.Database, collectionName string, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	ctx, ctxClose := context.WithTimeout(ctx, timeout)
	defer ctxClose()
	updateResponse, err := db.Collection(collectionName).UpdateOne(ctx, filter, update, opts...)
	if err != nil || updateResponse == nil {
		return 0, 0, err
	}

	return updateResponse.MatchedCount, updateResponse.ModifiedCount, err

}
