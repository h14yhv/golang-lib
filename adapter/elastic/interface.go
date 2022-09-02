package elastic

type Database interface {
	Get(database, collection, id string, result interface{}) error
	Exists(database, collection, id string) (bool, error)
	Count(database, collection string, query Query) (int64, error)
	FindOne(database, collection string, query Query, sorts []string, result interface{}) error
	FindPaging(database, collection string, query Query, sorts []string, page, size int, results interface{}) (int64, error)
	FindOffset(database, collection string, query Query, sorts []string, offset, size int, results interface{}) (int64, error)
	FindScroll(database, collection string, query Query, sorts []string, size int, scrollID, keepAlive string, results interface{}) (string, int64, error)
	InsertOne(database, collection string, doc Document) error
	InsertMany(database, collection string, docs []Document) error
	UpdateByID(database, collection, id string, update interface{}, upsert bool) error
	UpdateOne(database, collection string, query Query, update interface{}, upsert bool) error
	UpdateMany(database, collection string, query Query, update interface{}, upsert bool) error
	DeleteByID(database, collection, id string) error
	DeleteMany(database, collection string, query Query) error
}
