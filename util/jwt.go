package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/dgrijalva/jwt-go"
)

func GenerateJwt(user model.User) (string, error) {
	expired := strconv.FormatInt(time.Now().Add(time.Hour*12).Unix(), 10)
	claims := jwt.MapClaims{
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"expired":    expired,
		"user_role":  user.UserRole,
		"ID":         user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ToInt(param interface{}) (int, bool) {
	var id int
	ok := true

	switch v := param.(type) {
	case uint:
	case int8:
	case int16:
	case int32:
	case int64:
	case float32:
	case float64:
		id = int(v)
	case int:
		id = v
	case string:
		idParse, errParse := strconv.ParseInt(v, 10, 8)

		if errParse != nil {
			ok = false
		} else {
			id = int(idParse)
		}
	default:
		ok = false
	}

	return id, ok
}

func IsJwtValid(paramToken string) (bool, jwt.MapClaims) {
	if len(paramToken) == 0 {
		return false, nil
	}

	paramFinalToken := paramToken
	splits := strings.Split(paramFinalToken, "Bearer ")
	paramFinalToken = splits[len(splits)-1]

	if len(paramFinalToken) == 0 {
		return false, nil
	}

	token, err := jwt.Parse(paramFinalToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed Parse Token")
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return false, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		expired, okExp := claims["expired"].(string)
		expiredTime, err := strconv.ParseInt(expired, 10, 64)
		valid := expiredTime > time.Now().Unix() && err == nil && okExp

		return valid, claims
	}

	return false, nil
}
