package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/test"
)

var tokenUserFine string
var userBooksTestUf []model.UserBook
var userLoginUserFine model.User

func init() {
	tokenUserFineRes, user, err := test.InitLoginUser(test.LOGIN_ADMIN)
	if err != nil {
		log.Fatal(err)

		return
	}

	tokenUserFine = tokenUserFineRes
	userLoginUserFine = user
	userBooksInit := initDataUserBook(true, userLoginUserFine)
	userBooksTestUf = userBooksInit
}

func getUfToken() []string {
	return []string{"Authorization", fmt.Sprintf("Bearer %s", tokenUserFine)}
}

func TestUserReturnBookInit(t *testing.T) {
	resp, err := executeBorrowReturn(t, true, userLoginUserFine.ID, &userBooksTestUf)
	if err == nil {
		data := make(map[string]interface{})
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Error("err parse json %w", err)
			return
		}

		t.Log(data)
	}
}

func TestUserFineGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/user-fine-histories", nil, http.StatusOK, getUfToken())
		data := make(test.ResponseDataArray)
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataArray(t, data, "user", "fine", "has_paid", "user_book_ids", "ID")
		}
	})
}

func TestUserFineCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/user-fine-histories/paging/count", nil, http.StatusOK, getUfToken())
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

func TestUserFinePay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user_id := userLoginUserFine.ID
		userFines := []model.UserFineHistory{}
		paramFind := map[string]interface{}{"user_id": user_id, "has_paid": false}

		if err := test.GetApp().DB.Select("id", "user_id").Where(paramFind).Find(&userFines).Error; err != nil {
			t.Errorf("find user fine history err %w", err)
			return
		}

		if len(userFines) == 0 {
			t.Error("user fines empty")
			return
		}

		body := make(map[string][]map[string]int)
		for _, val := range userFines {
			item := map[string]int{
				"id": int(val.ID),
			}
			body["fines"] = append(body["fines"], item)
		}

		b, err := json.Marshal(body)
		if err != nil {
			t.Errorf("error parse to json : %w", err)
			return
		}

		data := make(map[string]interface{})
		url := fmt.Sprintf("/api/user-fine-histories/pay/%d", user_id)
		resp := test.ExecuteBaseRequest(t, "PUT", url, bytes.NewReader(b), http.StatusOK, getUfToken())

		if err := test.ParseJson(resp, &data); err != nil {
			t.Errorf("error parsing : %w", err)
		}

		t.Log(data)
		if err == nil {
			msg, ok := data["message"]

			if !ok {
				t.Error("data.message not showing")
			}

			if ok && msg != "Success paid fine" {
				t.Error("data.message should be Success paid fine")
			}
		}
	})
}
