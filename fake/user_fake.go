package fake

import (
	"fmt"
	"strings"
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/cahyowhy/go-basit-restapi-revisit/util"
	"github.com/jaswdr/faker"
)

func GetUser(password string) model.User {
	faker := faker.New()
	fakePerson := faker.Person()
	firstName := strings.ToLower(fakePerson.FirstName())
	lastName := strings.ToLower(fakePerson.LastName())
	email := fmt.Sprintf("%s_%s@mail.com", firstName, lastName)
	username := fmt.Sprintf("%s_%s", firstName, lastName)
	userRole := model.USER

	return model.User{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: fakePerson.Contact().Phone,
		Username:    username,
		Password:    password,
		BirthDate:   time.Date(1996, time.November, 12, 0, 0, 0, 0, time.UTC),
		UserRole:    userRole,
	}
}

func GetUsers(total int) []model.User {
	var users []model.User

	for i := 1; i <= total; i++ {
		password, err := util.GeneratePassword("12345678")

		if err != nil {
			continue
		}

		user := GetUser(password)
		if i <= 1 {
			user.FirstName = "admin"
			user.LastName = "admin"
			user.Username = "admin"
			user.Email = "admin@mail.com"
			user.UserRole = model.ADMIN
		}

		users = append(users, user)
	}

	return users
}
