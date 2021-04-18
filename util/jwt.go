package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
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

	return token.SignedString([]byte(config.GetConfig().JWTSECRET))
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

		jwtSecret := config.GetConfig().JWTSECRET

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
