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
var userTestUf []model.User

func init() {
	tokenUserFineRes, err := test.InitLoginUser(test.LOGIN_ADMIN)
	if err != nil {
		log.Fatal(err)

		return
	}

	tokenUserFine = tokenUserFineRes
	usersInit, userBooksInit := initDataUserBook(true)
	userTestUf = usersInit
	userBooksTestUf = userBooksInit
}

func getUfToken() []string {
	return []string{"Authorization", fmt.Sprintf("Bearer %s", tokenUserFine)}
}

func TestUserReturnBookInit(t *testing.T) {
	resp, err := executeBorrowReturn(t, true, userTestUf[0].ID, &userBooksTestUf)
	if err != nil {
		data := make(map[string]interface{})
		if err := test.ParseJson(resp, &data); err == nil {
			fmt.Println(data)
		}
	}
}

func TestUserFinePay(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		user_id := userTestUf[0].ID
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

func TestUserFineGetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/user-books", nil, http.StatusOK, getUfToken())
		data := make(test.ResponseDataArray)
		err := test.ParseJson(resp, &data)

		if err != nil {
			t.Errorf("error parsing : %w", err)
		}

		if err == nil {
			test.CheckVisibleDataArray(t, data, "borrow_date", "return_date", "book", "user", "ID")
		}
	})
}

func TestUserFineCount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := test.ExecuteBaseRequest(t, "GET", "/api/user-books/paging/count", nil, http.StatusOK, getUfToken())
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
