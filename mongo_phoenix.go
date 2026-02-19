// Package phoenix provides a generic, type-safe wrapper around the MongoDB Go driver
// that simplifies common CRUD operations. It reduces boilerplate by managing context
// timeouts and collection references automatically.
package phoenix

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Phoenix is a generic MongoDB collection wrapper. The type parameter D represents
// the document model. All operations are bound to the collection and timeout
// configured at construction time.
type Phoenix[D any] struct {
	db             *mongo.Database
	timeOut        time.Duration
	collectionName string
}

// New creates a new Phoenix instance bound to the given database, collection name,
// and operation timeout. The type parameter D defines the document model used for
// deserialization.
func New[D any](db *mongo.Database, collectionName string, timeOut time.Duration) *Phoenix[D] {
	return &Phoenix[D]{
		db:             db,
		timeOut:        timeOut,
		collectionName: collectionName,
	}
}

// FindOne retrieves a single document matching filterQuery from the collection.
// It returns the document, a boolean indicating whether the document was found,
// and any error encountered. If no document matches, it returns (zero value, false, nil).
func (p *Phoenix[D]) FindOne(ctx context.Context, filterQuery bson.M, opts ...*options.FindOneOptions) (D, bool, error) {
	return FindOne[D](ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

// FindAll retrieves all documents matching filterQuery from the collection.
// It returns a slice of documents and any error encountered. An empty result
// set is returned as a nil slice with a nil error.
func (p *Phoenix[D]) FindAll(ctx context.Context, filterQuery bson.M, opts ...*options.FindOptions) ([]D, error) {
	return FindAll[D](ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

// DeleteAll removes all documents matching filterQuery from the collection.
// It returns the number of deleted documents and any error encountered.
func (p *Phoenix[D]) DeleteAll(ctx context.Context, filterQuery bson.M, opts ...*options.DeleteOptions) (int64, error) {
	return DeleteAll(ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

// Count returns the number of documents in the collection that match filterQuery.
func (p *Phoenix[D]) Count(ctx context.Context, filterQuery bson.M, opts ...*options.CountOptions) (int64, error) {
	return Count(ctx, p.timeOut, p.db, p.collectionName, filterQuery, opts...)
}

// UpdateMany updates all documents matching filter using the provided update document.
// It returns the number of matched documents, the number of modified documents,
// and any error encountered.
func (p *Phoenix[D]) UpdateMany(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	return UpdateMany(ctx, p.timeOut, p.db, p.collectionName, filter, update, opts...)
}

// UpdateOne updates the first document matching filter using the provided update document.
// It returns the number of matched documents, the number of modified documents,
// and any error encountered.
func (p *Phoenix[D]) UpdateOne(ctx context.Context, filter bson.M, update bson.M, opts ...*options.UpdateOptions) (int64, int64, error) {
	return UpdateOne(ctx, p.timeOut, p.db, p.collectionName, filter, update, opts...)
}
