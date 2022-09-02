package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Database interface {
	CreateIndex(database, collection string, index *bson.M, unique bool) error
	Get(database, collection, id string, result interface{}) error
	Count(database, collection string, query *bson.M) (int64, error)
	FindOne(database, collection string, query *bson.M, sorts []string, offset int64, result interface{}) error
	FindMany(database, collection string, query *bson.M, sorts []string, size, offset int64, results interface{}) (int64, error)
	InsertOne(database, collection string, doc Document) error
	InsertMany(database, collection string, docs []Document, ordered bool) error
	UpdateByID(database, collection string, id interface{}, update interface{}) error
	UpdateOne(database, collection string, query *bson.M, update interface{}, upsert bool) error
	UpdateMany(database, collection string, query *bson.M, update interface{}, upsert bool) error
	DeleteByID(database, collection string, id interface{}) error
	DeleteOne(database, collection string, query *bson.M) error
	DeleteMany(database, collection string, query *bson.M) error
	Aggregate(database, collection string, pipeline []*bson.M, results interface{}) error
}
