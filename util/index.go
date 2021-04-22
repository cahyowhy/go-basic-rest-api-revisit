package util

import (
	"time"

	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	KeyOffset string = "offset"
	KeyLimit  string = "limit"
	KeyFilter string = "filter"
	KeyUser   string = "user"
)

func CountTotalFine(date time.Time) int {
	returnDate := time.Now()
	diffDay := returnDate.Sub(date).Hours() / 24

	if diffDay > 7 {
		lateDay := diffDay - 7

		return int(lateDay+0.5) * 1000
	}

	return 0
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
