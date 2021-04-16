package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJwt(payload jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	return token.SignedString([]byte(config.GetConfig().JWTSECRET))
}

func IsJwtValid(paramToken string) bool {
	if len(paramToken) == 0 {
		return false
	}

	token, _ := jwt.Parse(paramToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Failed Parse Token")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		expired, okExp := claims["expired"].(string)
		expiredTime, err := strconv.ParseInt(expired, 10, 64)

		return expiredTime > time.Now().Unix() && err == nil && okExp
	}

	return false
}

func ResponseSendJson(w http.ResponseWriter, response interface{}, httpStatus ...int) {
	w.Header().Set("Content-Type", "application/json")

	if len(httpStatus) > 0 {
		w.WriteHeader(httpStatus[0])
	}

	json.NewEncoder(w).Encode(response)
}

func GetReponseMessage(message string) map[string]interface{} {
	body := make(map[string]interface{})
	body["message"] = message

	return body
}

func GeneratePassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashed), err
}

func CompareHashPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
