package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/test"
)

var tokenUserBook string
var userBooksTestUb []model.UserBook
var userTestUb []model.User

type responseDataReturnBook struct {
	Data struct {
		Fine float64 `json:"fine"`
	} `json:"data"`
	Message string `json:"message"`
}

type responseDataBorrowBook struct {
	Data []model.UserBook `json:"data"`
}

func init() {
	tokenUserBookRes, err := test.InitLoginUser(test.LOGIN_ADMIN)
	if err != nil {
		log.Fatal(err)

		return
	}

	tokenUserBook = tokenUserBookRes
	usersInit, userBooksInit := initDataUserBook(false)
	userTestUb = usersInit
	userBooksTestUb = userBooksInit
}

func initDataUserBook(nilReturnDate bool) ([]model.User, []model.UserBook) {
	usersInit := fake.GetUsers(1)
	booksInit := fake.GetBooks(1)

	euCh, ebCh := createDataFirst(&usersInit), createDataFirst(&booksInit)
	eu, eb := <-euCh, <-ebCh

	userBooksInit := fake.GetUserBooks(booksInit, usersInit)

	if eu != nil {
		log.Fatal(eu)
	}

	if eb != nil {
		log.Fatal(eb)
	}

	if nilReturnDate {
		for i := 0; i < len(userBooksInit); i++ {
			userBooksInit[i].ReturnDate = sql.NullTime{
				Valid: false,
			}
		}
	}

	eubCh := createDataFirst(&userBooksInit)
	eub := <-eubCh

	if eub != nil {
		log.Fatal(eub)
	}

	return usersInit, userBooksInit
}

func createDataFirst(target interface{}) <-chan error {
	r := make(chan error)

	go func() {
		err := test.GetApp().DB.Create(target).Error
		r <- err
	}()

	return r
}

func getUbToken() []string {
	return []string{"Authorization", fmt.Sprintf("Bearer %s", tokenUserBook)}
}

func TestUserBookGetAll(t *testing.T) {
	resp := test.ExecuteBaseRequest(t, "GET", "/api/user-books", nil, http.StatusOK, getUbToken())
	data := make(test.ResponseDataArray)
	err := test.ParseJson(resp, &data)

	if err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataArray(t, data, "borrow_date", "return_date", "book", "user", "ID")
	}
}

func TestUserBookCount(t *testing.T) {
	resp := test.ExecuteBaseRequest(t, "GET", "/api/user-books/paging/count", nil, http.StatusOK, getUbToken())
	data := make(test.ResponseDataTotal)
	err := test.ParseJson(resp, &data)

	if err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataTotal(t, data)
	}
}

func TestUserBookBorrow(t *testing.T) {
	resp, err := executeBorrowReturn(t, false, userTestUb[0].ID, &userBooksTestUb)
	data := responseDataBorrowBook{}

	if err := test.ParseJson(resp, &data); err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		if len(data.Data) <= 0 {
			t.Error("user books borrowed not showing")
		}
	}
}
func TestUserBookReturn(t *testing.T) {
	resp, err := executeBorrowReturn(t, true, userTestUb[0].ID, &userBooksTestUb)
	data := responseDataReturnBook{}

	if err := test.ParseJson(resp, &data); err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		if data.Message != "return book succeed" {
			t.Error("message not showing")
		}

		if data.Data.Fine > 0 {
			t.Error("fine should be 0")
		}
	}
}

func executeBorrowReturn(t *testing.T, isReturn bool, userId uint, userBooksParams *[]model.UserBook) (*http.Response, error) {
	url := fmt.Sprintf("/api/user-books/borrows/%d", userId)
	method := "POST"

	if isReturn {
		method = "PUT"
		url = fmt.Sprintf("/api/user-books/returns/%d", userId)
	}

	body := make(map[string][]map[string]int)
	for _, v := range *userBooksParams {
		item := map[string]int{"book_id": int(v.ID)}
		body["books"] = append(body["books"], item)
	}

	b, err := json.Marshal(body)
	if err != nil {
		t.Errorf("error parse to json : %w", err)
		return nil, err
	}

	return test.ExecuteBaseRequest(t, method, url, bytes.NewReader(b), http.StatusOK, getUbToken()), nil
}
