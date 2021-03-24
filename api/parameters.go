package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func getQueryLimit(c *gin.Context) (int64, error) {
	limit := c.DefaultQuery("limit", "10")

	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "invalid 'limit' parameter")
	}

	return l, nil
}

func getQueryPage(c *gin.Context) (int64, error) {
	page := c.DefaultQuery("page", "1")

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid 'page' parameter")
	}

	return p, nil
}

func getQueryAddress(c *gin.Context, paramName string) (string, error) {
	addr := c.Param(paramName)

	return utils.ValidateAccount(addr)
}

func getQueryTimestamp(c *gin.Context) (int64, error) {
	timestamp := c.DefaultQuery("timestamp", "-1")

	t, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "invalid 'timestamp' parameter")
	}

	return t, nil
}
