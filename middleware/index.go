package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cahyowhy/go-basit-restapi-revisit/handler"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/dgrijalva/jwt-go"
)

var AuthenticateJWT handler.Adapter = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		auth := req.Header.Get("Authorization")
		valid, claims := util.IsJwtValid(auth)

		if !valid {
			util.ResponseSendJson(res, util.ToMapKey("message", "Unauthorize"), http.StatusUnauthorized)

			return
		}

		req = req.WithContext(context.WithValue(req.Context(), util.KeyUser, claims))

		next.ServeHTTP(res, req)
	})
}

var AuthenticateAdmin handler.Adapter = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		claims, okClaim := req.Context().Value(util.KeyUser).(jwt.MapClaims)

		if !okClaim {
			util.ResponseSendJson(res, util.ToMapKey("message", "Unauthorize"), http.StatusUnauthorized)

			return
		} else if param, ok := claims["user_role"]; ok && param != "ADMIN" {
			util.ResponseSendJson(res, util.ToMapKey("message", "Unauthorize"), http.StatusUnauthorized)

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
				req = req.WithContext(context.WithValue(req.Context(), util.KeyFilter, filter))
			}
		}

		offset, _ := strconv.ParseInt(query.Get("offset"), 10, 8)
		limit, _ := strconv.ParseInt(query.Get("limit"), 10, 8)

		if limit == 0 {
			limit = 20
		}

		req = req.WithContext(context.WithValue(req.Context(), util.KeyOffset, offset))
		req = req.WithContext(context.WithValue(req.Context(), util.KeyLimit, limit))

		next.ServeHTTP(res, req)
	})
}
