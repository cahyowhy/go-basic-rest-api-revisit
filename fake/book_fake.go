package fake

import (
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"syreclabs.com/go/faker"
)

func GetBooks(total int) []model.Book {
	var books []model.Book

	for i := 1; i <= total; i++ {
		book := model.Book{
			Title:        faker.Lorem().Sentence(2),
			Sheet:        uint(faker.RandomInt(100, 300)),
			Introduction: faker.Lorem().Paragraph(4),
			DateOffIssue: time.Date(2004, time.November, 12, 0, 0, 0, 0, time.UTC),
			Author:       []string{faker.Name().FirstName(), faker.Name().FirstName()},
		}

		books = append(books, book)
	}

	return books
}
