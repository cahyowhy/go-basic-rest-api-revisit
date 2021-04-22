package handler

import (
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/gofiber/fiber/v2"
)

type QueryParam struct {
	Offset int
	Limit  int
	Filter map[string]interface{}
}

func GetQueryParam(c *fiber.Ctx) QueryParam {
	offset, okOffset := c.Locals(util.KeyOffset).(int64)
	limit, okLimit := c.Locals(util.KeyLimit).(int64)

	if !okOffset {
		offset = 0
	}

	if !okLimit {
		limit = 20
	}

	queryParam := QueryParam{Offset: int(offset), Limit: int(limit)}

	filterFinal, ok := c.Locals(util.KeyFilter).(map[string]interface{})

	if ok {
		queryParam.Filter = filterFinal
	}

	return queryParam
}
