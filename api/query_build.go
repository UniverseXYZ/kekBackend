package api

import (
	"fmt"

	"github.com/lib/pq"
)

type Filters []QueryFilter

func (f *Filters) Add(key string, value interface{}) *Filters {
	*f = append(*f, QueryFilter{key, value})

	return f
}

type QueryFilter struct {
	Key   string
	Value interface{}
}

func buildQueryWithFilter(query string, filters *Filters, limit *int64, offset *int64) (string, []interface{}) {
	var where string
	var offsetFilter, limitFilter string
	var params []interface{}
	if len(*filters) > 0 {

		for _, filter := range *filters {
			if where != "" {
				where += " and "
			} else {
				where += "where "
			}

			switch filter.Value.(type) {
			case []string:
				params = append(params, pq.Array(filter.Value))

				where += fmt.Sprintf("%s= ANY($%d)", filter.Key, len(params))
			default:
				params = append(params, filter.Value)

				where += fmt.Sprintf("%s= $%d", filter.Key, len(params))
			}
		}
	}

	if offset != nil {
		params = append(params, offset)

		offsetFilter = fmt.Sprintf("offset $%d", len(params))
	}

	if limit != nil {
		params = append(params, limit)

		limitFilter = fmt.Sprintf("limit $%d", len(params))
	}

	return fmt.Sprintf(query, where, offsetFilter, limitFilter), params
}
