package util

import (
	"fmt"
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

func NestedMapLookup(m map[string]interface{}, ks ...string) (rval interface{}, err error) {
	var ok bool

	if len(ks) == 0 { // degenerate input
		return nil, fmt.Errorf("nestedMapLookup needs at least one key")
	}
	if rval, ok = m[ks[0]]; !ok {
		return nil, fmt.Errorf("key not found; remaining keys: %v", ks)
	} else if len(ks) == 1 { // we've reached the final key
		return rval, nil
	} else if m, ok = rval.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return NestedMapLookup(m, ks[1:]...)
	}
}
