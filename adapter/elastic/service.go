package elastic

import (
	"encoding/json"
	"errors"
	"reflect"

	es "github.com/olivere/elastic/v7"
)

type ModelV7 struct {
	model *ES
}

func NewService(conf Config) (Database, error) {
	con, err := newElastic(conf.String())
	if err != nil {
		return nil, err
	}
	// Success
	return &ModelV7{model: con}, nil
}

func (con *ModelV7) Get(database, _, id string, result interface{}) error {
	res, err := con.model.Get(database, id)
	if err != nil {
		if es.IsNotFound(err) {
			return errors.New(NotFoundError)
		}
		return err
	}
	err = json.Unmarshal(res.Source, result)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) Exists(database, _, id string) (bool, error) {
	// Success
	return con.model.Exists(database, id)
}

func (con *ModelV7) Count(database, _ string, query Query) (int64, error) {
	// Success
	return con.model.Count(database, query)
}

func (con *ModelV7) FindOne(database, _ string, query Query, sorts []string, result interface{}) error {
	res, err := con.model.SearchOffset(database, query, sorts, 0, 1)
	if err != nil || res.Hits == nil {
		return err
	}
	if res.Hits.TotalHits.Value == 0 {
		return errors.New(NotFoundError)
	}
	err = json.Unmarshal(res.Hits.Hits[0].Source, result)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) FindPaging(database, _ string, query Query, sorts []string, page, size int, results interface{}) (int64, error) {
	res, err := con.model.SearchPaging(database, query, sorts, page, size)
	if err != nil || res.Hits == nil {
		return 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return count, nil
}

func (con *ModelV7) FindOffset(database, _ string, query Query, sorts []string, offset, size int, results interface{}) (int64, error) {
	res, err := con.model.SearchOffset(database, query, sorts, offset, size)
	if err != nil || res.Hits == nil {
		return 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return count, nil
}

func (con *ModelV7) FindScroll(database, _ string, query Query, sorts []string, size int, scrollID, keepAlive string, results interface{}) (string, int64, error) {
	res, err := con.model.SearchScroll(database, query, sorts, size, scrollID, keepAlive)
	if err != nil || res.Hits == nil {
		return "", 0, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return "", 0, errors.New(NotFoundError)
	}
	resultType := reflect.TypeOf(results)
	resultValue := reflect.ValueOf(results)
	resultElemType := resultType.Elem().Elem()
	if resultType.Kind() != reflect.Ptr {
		return "", 0, errors.New(ResultNotAPointer)
	}
	count := res.Hits.TotalHits.Value
	for _, hit := range res.Hits.Hits {
		itemValue := reflect.New(resultElemType)
		err = json.Unmarshal(hit.Source, itemValue.Interface())
		if err != nil {
			return "", count, err
		}
		resultValue.Elem().Set(reflect.Append(resultValue.Elem(), itemValue.Elem()))
	}
	// Success
	return res.ScrollId, count, nil
}

func (con *ModelV7) InsertOne(database, _ string, doc Document) error {
	_, err := con.model.Index(database, doc)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) InsertMany(database, _ string, docs []Document) error {
	bulk := con.model.Bulk()
	for idx := range docs {
		bulk.Index(database, docs[idx].GetID(), docs[idx])
	}
	err := bulk.Do()
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) UpdateByID(database, _, id string, update interface{}, upsert bool) error {
	_, err := con.model.Update(database, id, update, upsert)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) UpdateOne(database, _ string, query Query, update interface{}, upsert bool) error {
	// TODO
	return nil
}

func (con *ModelV7) UpdateMany(database, _ string, query Query, update interface{}, upsert bool) error {
	// TODO
	return nil
}

func (con *ModelV7) DeleteByID(database, _, id string) error {
	_, err := con.model.DeleteByID(database, id)
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ModelV7) DeleteMany(database, _ string, query Query) error {
	_, err := con.model.DeleteByQuery(database, query)
	if err != nil {
		return err
	}
	// Success
	return nil
}
