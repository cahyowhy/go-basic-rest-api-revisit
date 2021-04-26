package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/test"
)

// var a int
var tokenBook string
var book model.Book

func init() {
	tokenBookRes, err := test.InitLoginUser(test.LOGIN_ADMIN)
	if err != nil {
		log.Fatal(err)

		return
	}

	tokenBook = tokenBookRes
	book = fake.GetBooks(1)[0]
}

func TestBookCreate(t *testing.T) {
	body, err := json.Marshal(book)
	if err != nil {
		t.Error(err.Error())
	}

	resp := test.ExecuteBaseRequest(t, "POST", "/api/books", bytes.NewReader(body), http.StatusOK, []string{"Authorization", fmt.Sprintf("Bearer %s", tokenBook)})
	data := make(test.ResponseDataMap)

	if err := test.ParseJson(resp, &data); err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		res := test.CheckVisibleDataMap(t, data, "title", "sheet", "date_off_issue", "introduction", "author", "ID")
		if val, ok := res["ID"]; ok {
			book.ID = uint(val.(float64))

			return
		}

		t.Error("book ID are empty !!")
	}
}

func TestBookGetAll(t *testing.T) {
	resp := test.ExecuteBaseRequest(t, "GET", "/api/books", nil, http.StatusOK, []string{"Authorization", fmt.Sprintf("Bearer %s", tokenBook)})
	data := make(test.ResponseDataArray)
	err := test.ParseJson(resp, &data)

	if err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataArray(t, data, "title", "sheet", "date_off_issue", "introduction", "author", "ID")
	}
}

func TestBookGet(t *testing.T) {
	resp := test.ExecuteBaseRequest(t, "GET", fmt.Sprintf("/api/books/%d", book.ID), nil, http.StatusOK, []string{"Authorization", fmt.Sprintf("Bearer %s", tokenBook)})
	data := make(test.ResponseDataMap)
	err := test.ParseJson(resp, &data)

	if err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataMap(t, data, "title", "sheet", "date_off_issue", "introduction", "author", "ID")
	}
}

func TestBookCount(t *testing.T) {
	resp := test.ExecuteBaseRequest(t, "GET", "/api/books/paging/count", nil, http.StatusOK, []string{"Authorization", fmt.Sprintf("Bearer %s", tokenBook)})
	data := make(test.ResponseDataTotal)
	err := test.ParseJson(resp, &data)

	if err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataTotal(t, data)
	}
}

func TestBookUpdate(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{"title": "How to 101 kill your self", "sheet": 210})
	if err != nil {
		t.Error(err.Error())
	}

	resp := test.ExecuteBaseRequest(t, "PUT", fmt.Sprintf("/api/books/%d", book.ID), bytes.NewReader(body), http.StatusOK, []string{"Authorization", fmt.Sprintf("Bearer %s", tokenBook)})
	data := make(test.ResponseDataMap)

	if err := test.ParseJson(resp, &data); err != nil {
		t.Errorf("error parsing : %w", err)
	}

	if err == nil {
		test.CheckVisibleDataMap(t, data, "title", "sheet")
	}
}
