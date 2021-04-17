package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
)

var AuthenticateJWT handler.Adapter = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")

		if !util.IsJwtValid(authHeader) {
			body := make(map[string]string)
			body["message"] = "Unauthorize"

			res.WriteHeader(http.StatusUnauthorized)
			util.ResponseSendJson(res, body)

			return
		}

		next.ServeHTTP(res, req)
	})
}

var ParseQueryFilter handler.Adapter = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		var filter = make(map[string]interface{})
		var filterString = query.Get("filter")

		if len(filterString) > 0 {
			if error := json.Unmarshal([]byte(filterString), &filter); error == nil {
				req = req.WithContext(context.WithValue(req.Context(), "filter", filter))
			}
		}

		offset, _ := strconv.ParseInt(query.Get("offset"), 10, 8)
		limit, _ := strconv.ParseInt(query.Get("limit"), 10, 8)

		if limit == 0 {
			limit = 20
		}

		req = req.WithContext(context.WithValue(req.Context(), "offset", offset))
		req = req.WithContext(context.WithValue(req.Context(), "limit", limit))

		next.ServeHTTP(res, req)
	})
}
