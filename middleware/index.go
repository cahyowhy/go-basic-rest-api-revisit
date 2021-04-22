package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var AuthenticateJWT fiber.Handler = func(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	valid, claims := util.IsJwtValid(auth)

	if !valid {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "Unauthorize"))
	}

	c.Locals(util.KeyUser, claims)

	return c.Next()
}

var AuthenticateAdmin fiber.Handler = func(c *fiber.Ctx) error {
	claims, okClaim := c.Locals(util.KeyUser).(jwt.MapClaims)

	if !okClaim {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "Unauthorize"))
	} else if param, ok := claims["user_role"]; ok && param != "ADMIN" {
		return c.Status(http.StatusUnauthorized).JSON(util.ToMapKey("message", "Unauthorize"))
	}

	return c.Next()
}

var ParseQueryFilter fiber.Handler = func(c *fiber.Ctx) error {
	var filter = make(map[string]interface{})
	var filterString = c.Query("filter")

	if len(filterString) > 0 {
		if error := json.Unmarshal([]byte(filterString), &filter); error == nil {
			c.Locals(util.KeyFilter, filter)
		}
	}

	offset, _ := strconv.ParseInt(c.Query("offset"), 10, 8)
	limit, _ := strconv.ParseInt(c.Query("limit"), 10, 8)

	if limit == 0 {
		limit = 20
	}

	c.Locals(util.KeyOffset, offset)
	c.Locals(util.KeyLimit, limit)

	return c.Next()
}
