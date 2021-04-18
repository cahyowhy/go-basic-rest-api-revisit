package handler

import (
	"net/http"

	"github.com/cahyowhy/go-basit-restapi-revisit/util"
)

type Adapter func(http.Handler) http.Handler

type QueryParam struct {
	Offset int
	Limit  int
	Filter map[string]interface{}
}

func GetQueryParam(r *http.Request) QueryParam {
	offset, okOffset := r.Context().Value(util.KeyOffset).(int64)
	limit, okLimit := r.Context().Value(util.KeyLimit).(int64)

	if !okOffset {
		offset = 0
	}

	if !okLimit {
		limit = 20
	}

	queryParam := QueryParam{Offset: int(offset), Limit: int(limit)}

	var filter = r.Context().Value(util.KeyFilter)
	filterFinal, ok := filter.(map[string]interface{})

	if ok {
		queryParam.Filter = filterFinal
	}

	return queryParam
}

func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	// The loop is reversed so the adapters/middleware gets executed in the same
	// order as provided in the array.
	for i := len(adapters); i > 0; i-- {
		handler = adapters[i-1](handler)
	}
	return handler
}
