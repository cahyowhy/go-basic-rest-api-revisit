package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/test"
)

var a int
var user model.User
var token string
var cookieRefreshToken *http.Cookie

func init() {
	a = 2
	user = fake.GetUser("1234678")
}

func getUToken() []string {
	return []string{"Authorization", fmt.Sprintf("Bearer %s", token)}
}

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body, err := json.Marshal(user)
		if err != nil {
			t.Error(err.Error())
		}

		resp := test.ExecuteBaseRequest(t, "POST", "/api/users", bytes.NewReader(body), http.StatusOK)
		data := make(test.ResponseDataMap)

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			res := test.CheckVisibleDataMap(t, data, "first_name", "last_name", "email", "phone_number", "username", "ID")
			if val, ok := res["ID"]; ok {
				user.ID = uint(val.(float64))

				return
			}

			t.Error("user ID are empty !!")
		}
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body, err := json.Marshal(map[string]string{"username": user.Username, "password": user.Password})
		if err != nil {
			t.Errorf("error parse to json : %w", err)
		}

		resp := test.ExecuteBaseRequest(t, "POST", "/api/users/auth/login", bytes.NewReader(body), http.StatusOK)
		data := make(test.ResponseDataMap)

		cookies := resp.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "refresh_token" {
				cookieRefreshToken = cookie

				break
			}
		}

		if cookieRefreshToken == nil {
			t.Error("refresh_token not found")
		}

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			res := test.CheckVisibleDataMap(t, data, "token", "first_name", "last_name", "email", "phone_number")
			if val, ok := res["token"]; ok {
				token = val.(string)

				return
			}

			t.Error("token are empty !!")
		}
	})
}

func TestGetAllUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/users", nil, http.StatusOK, getUToken())
		data := make(test.ResponseDataArray)
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataArray(t, data, "username", "email", "first_name", "last_name", "phone_number", "ID")
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", fmt.Sprintf("/api/users/%d", user.ID), nil, http.StatusOK, getUToken())
		data := make(test.ResponseDataMap)
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataMap(t, data, "username", "email", "first_name", "last_name", "phone_number", "ID")
		}
	})
}

func TestCountUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/users/paging/count", nil, http.StatusOK, getUToken())
		data := make(test.ResponseDataTotal)
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataTotal(t, data)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body, err := json.Marshal(map[string]interface{}{"first_name": "Samantha", "last_name": "Kholil"})
		if err != nil {
			t.Error(err.Error())
		}

		resp := test.ExecuteBaseRequest(t, "PUT", fmt.Sprintf("/api/users/%d", user.ID), bytes.NewReader(body), http.StatusOK, getUToken())
		data := make(test.ResponseDataMap)

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataMap(t, data, "first_name", "last_name")
		}
	})
}

func TestSessionUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if cookieRefreshToken == nil {
			t.Error("refresh token cookie not found")
			return
		}

		data := make(test.ResponseDataMap)
		req, err := http.NewRequest("GET", "/api/users/auth/session", nil)
		if err != nil {
			t.Errorf("failed init req : %w", err)
		}

		req.AddCookie(cookieRefreshToken)

		resp, err := test.GetResp(req)
		if err != nil {
			t.Errorf("failed init resp : %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("http status should %d instead got %d", http.StatusOK, resp.StatusCode)
		}

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataMap(t, data, "token", "username", "email", "first_name", "last_name", "phone_number", "ID")
		}
	})
}

func TestLogoutUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		if cookieRefreshToken == nil {
			t.Error("refresh token cookie not found")
			return
		}

		data := make(map[string]interface{})
		req, err := http.NewRequest("GET", "/api/users/auth/logout", nil)
		if err != nil {
			t.Errorf("failed init req : %w", err)
		}

		authHead := getUToken()
		req.Header.Set(authHead[0], authHead[1])
		req.AddCookie(cookieRefreshToken)

		resp, err := test.GetResp(req)
		if err != nil {
			t.Errorf("failed init resp : %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("http status should %d instead got %d", http.StatusOK, resp.StatusCode)
		}

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			msg, ok := data["message"]

			if !ok {
				t.Error("data.message not showing")
			}

			if ok && msg != "logged out succeed" {
				t.Error("data.message should be logged out succeed")
			}
		}
	})
}
