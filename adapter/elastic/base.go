package elastic

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	es "github.com/olivere/elastic/v7"
)

type (
	ES struct {
		model *es.Client
	}

	BulkService struct {
		helper *es.BulkService
	}

	ScrollService struct {
		helper *es.ScrollService
	}
)

func newElastic(addr string) (*ES, error) {
	client, err := es.NewClient(es.SetURL(addr), es.SetSniff(false))
	if err != nil {
		return nil, err
	}
	// Success
	return &ES{
		model: client,
	}, nil
}

func (con *ES) Bulk() *BulkService {
	// Success
	return &BulkService{helper: es.NewBulkService(con.model)}
}

func (bs *BulkService) Index(index, id string, doc interface{}) {
	req := es.NewBulkIndexRequest().
		Index(index).
		Id(id).
		Doc(doc)
	bs.helper = bs.helper.Add(req)
	// Success
	return
}

func (bs *BulkService) Update(index, id string, update interface{}, upsert bool) {
	req := es.NewBulkUpdateRequest().
		Index(index).
		Id(id).
		Doc(update).
		DocAsUpsert(upsert)
	bs.helper = bs.helper.Add(req)
	// Success
	return
}

func (bs *BulkService) Delete(index, id string) {
	req := es.NewBulkDeleteRequest().
		Index(index).
		Id(id)
	bs.helper = bs.helper.Add(req)
	// Success
	return
}

func (bs *BulkService) Do() error {
	defer bs.helper.Reset()
	_, err := bs.helper.Refresh("true").Do(context.Background())
	if err != nil {
		return err
	}
	// Success
	return nil
}

func (con *ES) Scroll() *ScrollService {
	// Success
	return &ScrollService{helper: es.NewScrollService(con.model)}
}

func (ss *ScrollService) Search(index string, query Query, sorts []string, size int, scrollID, keepAlive string) (*es.SearchResult, error) {
	service := ss.helper.Index(index).Query(query).Size(size).ScrollId(scrollID).Scroll(keepAlive)
	if sorts != nil && len(sorts) > 0 {
		for _, sort := range sorts {
			if strings.HasPrefix(sort, "-") {
				service = service.Sort(strings.TrimPrefix(sort, "-"), false)
			} else if strings.HasPrefix(sort, "+") {
				service = service.Sort(strings.TrimPrefix(sort, "+"), true)
			}
		}
	}
	result, err := service.Do(context.Background())
	if err != nil {
		if err == io.EOF {
			return nil, errors.New(NotFoundError)
		}
		if ex := err.(*es.Error); ex.Status == http.StatusBadRequest || ex.Status == http.StatusNotFound {
			return nil, errors.New(NotFoundError)
		}
		return nil, err
	}
	// Success
	return result, nil
}

func (con *ES) IndexExists(index string) (bool, error) {
	// Success
	return con.model.IndexExists(index).
		Do(context.Background())
}

func (con *ES) CreateIndex(index string, mapping M) (*es.IndicesCreateResult, error) {
	// Success
	return con.model.CreateIndex(index).
		BodyJson(mapping).
		Do(context.Background())
}

func (con *ES) DeleteIndex(index string) (*es.IndicesDeleteResponse, error) {
	// Success
	return con.model.DeleteIndex(index).
		Do(context.Background())
}

func (con *ES) Count(index string, query Query) (int64, error) {
	// Success
	return con.model.Count().
		Index(index).
		Query(query).
		Do(context.Background())
}

func (con *ES) SearchScroll(index string, query Query, sorts []string, site int, scrollID, keepAlive string) (*es.SearchResult, error) {
	// Success
	return con.Scroll().Search(index, query, sorts, site, scrollID, keepAlive)
}

func (con *ES) SearchPaging(index string, query Query, sorts []string, page, size int) (*es.SearchResult, error) {
	// Success
	return con.SearchOffset(index, query, sorts, page*size, size)
}

func (con *ES) SearchOffset(index string, query Query, sorts []string, offset, size int) (*es.SearchResult, error) {
	service := con.model.Search().Index(index)
	if sorts != nil && len(sorts) > 0 {
		for _, sort := range sorts {
			if strings.HasPrefix(sort, "-") {
				service = service.Sort(strings.TrimPrefix(sort, "-"), false)
			} else if strings.HasPrefix(sort, "+") {
				service = service.Sort(strings.TrimPrefix(sort, "+"), true)
			}
		}
	}
	service = service.Query(query)
	if size == 0 {
		size = 10
	}
	service = service.Size(size).From(offset)
	result, err := service.Do(context.Background())
	if err != nil {
		if ex := err.(*es.Error); ex.Status == http.StatusBadRequest || ex.Status == http.StatusNotFound {
			return nil, errors.New(NotFoundError)
		}
		return nil, err
	}
	// Success
	return result, nil
}

func (con *ES) Get(index, id string) (*es.GetResult, error) {
	// Success
	return con.model.Get().
		Index(index).
		Id(id).
		Refresh("true").
		Do(context.Background())
}

func (con *ES) Exists(index, id string) (bool, error) {
	// Success
	return con.model.Exists().
		Index(index).
		Id(id).
		Refresh("true").
		Do(context.Background())
}

func (con *ES) Index(index string, doc Document) (*es.IndexResponse, error) {
	// Success
	return con.model.Index().
		Index(index).
		Id(doc.GetID()).
		BodyJson(doc).
		Refresh("true").
		Do(context.Background())
}

func (con *ES) Update(index, id string, update interface{}, upsert bool) (*es.UpdateResponse, error) {
	// Success
	return con.model.Update().
		Index(index).
		Id(id).
		Doc(update).
		DocAsUpsert(upsert).
		Refresh("true").
		Do(context.Background())
}

func (con *ES) DeleteByID(index, id string) (*es.DeleteResponse, error) {
	// Success
	return con.model.Delete().
		Index(index).
		Id(id).
		Refresh("true").
		Do(context.Background())
}

func (con *ES) DeleteByQuery(index string, query Query) (*es.BulkIndexByScrollResponse, error) {
	// Success
	return con.model.DeleteByQuery().
		Index(index).
		Query(query).
		Refresh("true").
		Do(context.Background())
}
