package util

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type keyCtk int

const (
	KeyOffset keyCtk = iota
	KeyLimit
	KeyFilter
	KeyUser
)

func ResponseSendJson(w http.ResponseWriter, response interface{}, httpStatus ...int) {
	w.Header().Set("Content-Type", "application/json")

	if len(httpStatus) > 0 {
		w.WriteHeader(httpStatus[0])
	}

	json.NewEncoder(w).Encode(response)
}

func CountTotalFine(date time.Time) int {
	returnDate := time.Now()
	diffDay := returnDate.Sub(date).Hours() / 24

	if diffDay > 7 {
		lateDay := diffDay - 7

		return int(lateDay+0.5) * 1000
	}

	return 0
}

func GetReponseData(data interface{}) map[string]interface{} {
	return ToMapKey("data", data)
}

func ToMapKey(key string, data interface{}) map[string]interface{} {
	body := make(map[string]interface{})
	body[key] = data

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

func ValidationErrToString(param validator.ValidationErrors) []string {
	params := []string{}

	for _, e := range param {
		params = append(params, e.StructNamespace()+" Are not valid")
	}

	return params
}

func GetUUID() (string, error) {
	u, err := uuid.NewV4()

	return u.String(), err
}
