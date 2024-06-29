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

// FindOne returns a single document from the specified collection in the MongoDB database.
func FindOne[D any](ctx context.Context, timeOut time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.FindOneOptions) (D, bool, error) {

	ctx, ctxClose := context.WithTimeout(ctx, timeOut)
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

// FindAll retrieves all documents from the specified collection in the MongoDB database that match the given filter query.
func FindAll[D any](ctx context.Context, timeOut time.Duration, db *mongo.Database, collectionName string, filterQuery bson.M, opts ...*options.FindOptions) ([]D, error) {

	ctx, ctxClose := context.WithTimeout(ctx, timeOut)
	defer ctxClose()

	cur, err := db.Collection(collectionName).Find(ctx, filterQuery, opts...)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []D{}, nil
		}
		return nil, err
	}
	defer func() {
		if err := cur.Close(ctx); err != nil {
			slog.Error("Find  close err:", slog.String("error", err.Error()))
		}

	}()
	var documents []D
	for cur.Next(ctx) {
		var d D
		err := cur.Decode(&d)
		if err != nil {
			slog.Error("Error during decode", slog.String("error", err.Error()))
			return nil, err
		}
		documents = append(documents, d)
	}
	return documents, nil
}
