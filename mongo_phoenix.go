package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Phoenix[D any] struct {
	db             *mongo.Database
	timeOut        time.Duration
	collectionName string
}

func New[D any](db *mongo.Database, collectionName string, timeOut time.Duration) *Phoenix[D] {
	return &Phoenix[D]{
		db:             db,
		timeOut:        timeOut,
		collectionName: collectionName,
	}
}

func (p *Phoenix[D]) FindOne(ctx context.Context, filterQuery bson.M, opts ...*options.FindOneOptions) (D, bool, error) {

	return FindOne[D](ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

func (p *Phoenix[D]) FindAll(ctx context.Context, filterQuery bson.M, opts ...*options.FindOptions) ([]D, error) {

	return FindAll[D](ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

func (p *Phoenix[D]) DeleteAll(ctx context.Context, filterQuery bson.M, opts ...*options.DeleteOptions) (int64, error) {
	return DeleteAll(ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

func (p *Phoenix[D]) Count(ctx context.Context, filterQuery bson.M, opts ...*options.CountOptions) (int64, error) {
	return Count(ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

func (p *Phoenix[D]) UpdateMany(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	return UpdateMany(ctx, p.timeOut, p.db, p.collectionName, filter, update, opts...)
}

func (p *Phoenix[D]) UpdateOne(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	return UpdateOne(ctx, p.timeOut, p.db, p.collectionName, filter, update, opts...)
}
