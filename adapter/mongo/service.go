package mongo

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Model struct {
		model *mongo.Client
	}
)

func NewService(conf Config) (Database, error) {
	if conf.Auth.Enable {
		con, err := mongo.NewClient(
			options.Client().ApplyURI(conf.String()),
			options.Client().SetAuth(options.Credential{
				AuthSource: conf.Auth.AuthDB,
				Username:   conf.Auth.Username,
				Password:   conf.Auth.Password,
			}),
		)
		if err != nil {
			return nil, err
		}
		if err = con.Connect(context.Background()); err != nil {
			return nil, err
		}
		// Success
		return &Model{model: con}, nil
	} else {
		con, err := mongo.NewClient(
			options.Client().ApplyURI(conf.String()),
		)
		if err != nil {
			return nil, err
		}
		if err = con.Connect(context.Background()); err != nil {
			return nil, err
		}
		// Success
		return &Model{model: con}, nil
	}
}

func (con *Model) CreateIndex(database, collection string, index *bson.M, unique bool) error {
	opts := options.Index()
	opts.SetBackground(true)
	opts.SetUnique(unique)
	_, err := con.model.Database(database).Collection(collection).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    index,
		Options: opts,
	})
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) Get(database, collection, id string, result interface{}) error {
	res := con.model.Database(database).Collection(collection).FindOne(context.Background(), &bson.M{"_id": id})
	if err := res.Err(); err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return errors.New(NotFoundError)
		}
		return err
	}
	if err := res.Decode(result); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) Count(database, collection string, query *bson.M) (int64, error) {
	res, err := con.model.Database(database).Collection(collection).CountDocuments(context.Background(), query)
	if err != nil {
		return 0, err
	}
	// Success
	return res, err
}

func (con *Model) FindOne(database, collection string, query *bson.M, sorts []string, offset int64, result interface{}) error {
	opts := options.FindOne()
	opts.SetSkip(offset)
	if sorts != nil && len(sorts) > 0 {
		s := bson.D{}
		for _, sort := range sorts {
			if strings.HasPrefix(sort, "-") {
				s = append(s, bson.E{Key: strings.TrimPrefix(sort, "-"), Value: -1})
			} else if strings.HasPrefix(sort, "+") {
				s = append(s, bson.E{Key: strings.TrimPrefix(sort, "+"), Value: 1})
			}
		}
		opts.SetSort(s)
	}
	res := con.model.Database(database).Collection(collection).FindOne(context.Background(), query, opts)
	if err := res.Err(); err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return errors.New(NotFoundError)
		}
		return err
	}
	if err := res.Decode(result); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) FindMany(database, collection string, query *bson.M, sorts []string, size, offset int64, results interface{}) (int64, error) {
	opts := options.Find()
	opts.SetBatchSize(DefaultBatchSize)
	if offset > 0 {
		opts.SetSkip(offset)
	}
	if size > 0 {
		opts.SetLimit(size)
	}
	if sorts != nil && len(sorts) > 0 {
		s := bson.D{}
		for _, sort := range sorts {
			if strings.HasPrefix(sort, "-") {
				s = append(s, bson.E{Key: strings.TrimPrefix(sort, "-"), Value: -1})
			} else if strings.HasPrefix(sort, "+") {
				s = append(s, bson.E{Key: strings.TrimPrefix(sort, "+"), Value: 1})
			}
		}
		opts.SetSort(s)
	}
	cur, err := con.model.Database(database).Collection(collection).Find(context.Background(), query, opts)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return 0, errors.New(NotFoundError)
		}
		return 0, err
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New(ResultNotAPointer)
	}
	var total int64 = 0
	for cur.Next(context.Background()) {
		itemValue := reflect.New(resultElemType)
		if err = cur.Decode(itemValue.Interface()); err != nil {
			return 0, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
		total += 1
	}
	// Success
	return total, nil
}

func (con *Model) InsertOne(database, collection string, doc Document) error {
	opts := options.InsertOne()
	opts.SetBypassDocumentValidation(true)
	if _, err := con.model.Database(database).Collection(collection).InsertOne(context.Background(), doc, opts); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) InsertMany(database, collection string, docs []Document, ordered bool) error {
	opts := options.InsertMany()
	opts.SetBypassDocumentValidation(true)
	opts.SetOrdered(ordered)
	if len(docs) == 0 {
		return nil
	}
	documents := make([]interface{}, 0)
	for _, doc := range docs {
		documents = append(documents, doc)
	}
	if _, err := con.model.Database(database).Collection(collection).InsertMany(context.Background(), documents, opts); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) UpdateByID(database, collection string, id interface{}, update interface{}) error {
	if _, err := con.model.Database(database).Collection(collection).UpdateByID(context.Background(), id, update); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) UpdateOne(database, collection string, query *bson.M, update interface{}, upsert bool) error {
	opts := options.Update()
	opts.SetUpsert(upsert)
	opts.SetBypassDocumentValidation(true)
	if _, err := con.model.Database(database).Collection(collection).UpdateOne(context.Background(), query, update); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) UpdateMany(database, collection string, query *bson.M, update interface{}, upsert bool) error {
	opts := options.Update()
	opts.SetUpsert(upsert)
	opts.SetBypassDocumentValidation(true)
	_, err := con.model.Database(database).Collection(collection).UpdateMany(context.Background(), query, update)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) DeleteByID(database, collection string, id interface{}) error {
	if _, err := con.model.Database(database).Collection(collection).DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) DeleteOne(database, collection string, query *bson.M) error {
	if _, err := con.model.Database(database).Collection(collection).DeleteOne(context.Background(), query); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) DeleteMany(database, collection string, query *bson.M) error {
	if _, err := con.model.Database(database).Collection(collection).DeleteMany(context.Background(), query); err != nil {
		return err
	}
	// Success
	return nil
}

func (con *Model) Aggregate(database, collection string, pipeline []*bson.M, results interface{}) error {
	opts := options.Aggregate()
	opts.SetBypassDocumentValidation(true)
	cur, err := con.model.Database(database).Collection(collection).Aggregate(context.Background(), pipeline, opts)
	if err != nil {
		return err
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return errors.New(ResultNotAPointer)
	}
	for cur.Next(context.Background()) {
		itemValue := reflect.New(resultElemType)
		if err = cur.Decode(itemValue.Interface()); err != nil {
			return err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return nil
}
