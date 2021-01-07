package api

import (
	"strconv"
)

func calculateOffset(limit string, page string) (string, error) {
	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return "", err
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return "", err
	}

	offset := (p - 1) * l

	return strconv.FormatInt(offset, 10), nil
}
