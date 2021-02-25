package api

import "fmt"

func buildQueryWithFilter(query string, filters map[string]interface{}, limit *int64, offset *int64) (string, []interface{}) {
	var where string
	var offsetFilter, limitFilter string
	var params []interface{}

	for k, v := range filters {
		if where != "" {
			where += " and "
		}

		params = append(params, v)

		where += fmt.Sprintf("%s = $%d", k, len(params))
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
