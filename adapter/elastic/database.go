package elastic

type (
	// Map
	M map[string]interface{}
	// Query
	Doc      M
	Query    M
	QueryES2 M
	QueryES7 M
	// Document
	Document interface {
		GetID() string
	}
)

func (doc Doc) GetID() string {
	if value, ok := doc["id"]; ok {
		return value.(string)
	}
	return ""
}

func (q Query) Source() (interface{}, error) {
	return q, nil
}
