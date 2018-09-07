package api

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/jinzhu/copier"
)

// Query is used for find calls
type Query struct {
	// Offset of the query
	Offset int
	// Limit in elements of the query
	Limit int
	// Sort key, minus can be used for reverse
	Sort string
	// Filter based on Mongo query definition
	Filter string
}

// RequestToQuery helper function to build Mongo Query from a URL
func RequestToQuery(url *url.URL) (*Query, error) {
	query := &Query{
		Offset: 0,
		Limit:  10,
		Sort:   "",
		Filter: "",
	}

	var filter interface{}
	for key, value := range url.Query() {
		if key == "limit" {
			limit, err := strconv.Atoi(value[0])
			if err == nil {
				query.Limit = limit
			}
		} else if key == "offset" {
			offset, err := strconv.Atoi(value[0])
			if err == nil {
				query.Offset = offset
			}
		} else if key == "sort" {
			query.Sort = value[0]
		} else if key == "query" {
			var filterTemplate interface{}
			err := json.Unmarshal([]byte(value[0]), &filterTemplate)
			if err == nil {
				filter = filterTemplate.(map[string]interface{})
			}
		} else {
			if filter == nil {
				filter = make(map[string]string)
			}
			filter.(map[string]string)[key] = value[0]
		}
	}

	if filter != nil {
		if filterBytes, err := json.Marshal(filter); err != nil {
			return nil, err
		} else {
			query.Filter = string(filterBytes)
		}
	}

	return query, nil
}

// Result is used for find result to format response
type Result struct {
	// Total items founs
	Total int `json:"total"`
	// Offset of the query
	Offset int `json:"offset"`
	// Limit in elements of the query
	Limit int `json:"limit"`
	// Sort key, minus can be used for reverse
	Sort string `json:"sort"`
	// Items returned
	Items interface{} `json:"items"`
}

// ResultToBody format a raw result to a well formed API result
func ResultToBody(result interface{}) (*Result, error) {
	bodyResult := &Result{}
	if result != nil {
		if err := copier.Copy(bodyResult, result); err != nil {
			return nil, err
		}
	}

	// TODO find a way to return none null -> below is KO
	if bodyResult.Items == nil {
		bodyResult.Items = make([]interface{}, 0)
	}
	return bodyResult, nil
}
