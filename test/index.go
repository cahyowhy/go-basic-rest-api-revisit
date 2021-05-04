package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/app"
	"github.com/cahyowhy/go-basit-restapi-revisit/config"
	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/joho/godotenv"
)

var appTest *app.App

type ResponseDataArray map[string][]map[string]interface{}
type ResponseDataMap map[string]map[string]interface{}
type ResponseDataTotal map[string]int

func GetConfigTest() *config.Config {
	if err := godotenv.Load(".test.env"); err != nil {
		if err := godotenv.Load("../.test.env"); err != nil {
			log.Fatalf(".test.env : %s", err.Error())

			return nil
		}
	}

	return &config.Config{
		AppEnv: os.Getenv("APP_ENV"),
		DbConfig: config.DbConfig{
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		JWTSECRET: os.Getenv("JWT_SECRET"),
		PORT:      os.Getenv("PORT"),
	}
}

func GetApp() *app.App {
	if appTest == nil {
		var configApp = GetConfigTest()
		appTest = app.GetApp(configApp)
	}

	return appTest
}

func GetResp(req *http.Request) (*http.Response, error) {
	return GetApp().FiberApp.Test(req)
}

func ParseJson(resp *http.Response, dest interface{}) error {
	if resp == nil {
		return errors.New("resp is nil")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dest)
}

func ExecuteBaseRequest(t *testing.T, method, url string, body io.Reader, expStatus int, headers ...[]string) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Errorf("failed init req : %w", err)
	}

	if err == nil && len(headers) > 0 {
		for _, h := range headers {
			if len(h) == 2 {
				req.Header.Set(h[0], h[1])
			}
		}
	}

	resp, err := GetResp(req)
	if err != nil {
		t.Errorf("failed init resp : %w", err)
	}

	if expStatus != resp.StatusCode {
		t.Errorf("http status should %d instead got %d", expStatus, resp.StatusCode)

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Log(string(b))
		}
	}

	return resp
}

func CheckVisibleDataArray(t *testing.T, data ResponseDataArray, expKeys ...string) {
	vals, ok := data["data"]
	if !ok {
		t.Error("data not found")
	}

	if !(ok && len(vals) > 0) {
		t.Error("data are empty")
	}

	if len(vals) > 0 && len(expKeys) > 0 {
		for _, expKey := range expKeys {
			prop, ok := vals[0][expKey]
			if ok {
				switch v := prop.(type) {
				case string:
					if len(v) <= 0 {
						t.Errorf("%s are empty", expKey)
					}
				}
			}

			if !ok {
				t.Errorf("%s are not present", expKey)
			}
		}
	}
}

func CheckVisibleDataTotal(t *testing.T, data ResponseDataTotal) {
	_, ok := data["data"]
	if !ok {
		t.Error("data not found")
	}
}

func CheckVisibleDataMap(t *testing.T, data ResponseDataMap, expKeys ...string) map[string]interface{} {
	expKeyRes := make(map[string]interface{})

	vals, ok := data["data"]
	if !ok {
		t.Error("data not found")
	}

	if !(ok && len(vals) > 0) {
		t.Error("data are empty")
	}

	if len(vals) > 0 && len(expKeys) > 0 {
		for _, expKey := range expKeys {
			prop, ok := vals[expKey]

			if !ok {
				t.Errorf("%s are not present", expKey)
				continue
			}

			switch v := prop.(type) {
			case string:
				if len(v) <= 0 {
					t.Errorf("%s are empty", expKey)
					continue
				}

				expKeyRes[expKey] = v
			default:
				expKeyRes[expKey] = v
			}
		}
	}

	return expKeyRes
}

type TypeLogin int

const (
	LOGIN_USER  TypeLogin = iota
	LOGIN_ADMIN TypeLogin = iota
)

// create user from db => login => token
func InitLoginUser(loginType TypeLogin) (string, model.User, error) {
	user := fake.GetUsers(1)[0]

	if loginType == LOGIN_ADMIN {
		user.UserRole = model.ADMIN
	}

	if err := GetApp().DB.Create(&user).Error; err != nil {
		return "", user, err
	}

	body, err := json.Marshal(map[string]string{"username": user.Username, "password": "12345678"})
	if err != nil {
		return "", user, err
	}

	req, err := http.NewRequest("POST", "/api/users/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", user, err
	}

	resp, err := GetResp(req)
	if err != nil {
		return "", user, err
	}

	data := make(map[string]interface{})
	if err := ParseJson(resp, &data); err != nil {
		return "", user, err
	}

	val, err := util.NestedMapLookup(data, "data", "token")
	if err != nil {
		return "", user, err
	}

	valStr, ok := val.(string)
	if !ok {
		return "", user, errors.New("data.token are not string")
	}

	return valStr, user, nil
}
