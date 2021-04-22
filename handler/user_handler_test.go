package handler

import (
	"testing"

	"github.com/cahyowhy/go-basit-restapi-revisit/fake"
	"github.com/cahyowhy/go-basit-restapi-revisit/model"
)

var a int
var user model.User

func init() {
	a = 2
	user = fake.GetUser("1234678")
}

func TestLogin(t *testing.T) {
	// body, err := json.Marshal(user)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// req, err := http.NewRequest("POST", "/users", bytes.NewReader(body))
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// resp := httptest.NewRecorder()
	// test.GetApp(config.GetConfig(".test.env")).Subrouter.ServeHTTP(resp, req)

	// if http.StatusOK != resp.Code {
	// 	t.Errorf("http status should %d instead got", http.StatusOK, resp.Code)
	// }
}
