package util

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
