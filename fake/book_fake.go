package fake

import (
	"time"

	"github.com/cahyowhy/go-basit-restapi-revisit/model"
	"github.com/jaswdr/faker"
)

func GetBooks(total int) []model.Book {
	var books []model.Book
	faker := faker.New()
	fakeLorem := faker.Lorem()
	fakePerson := faker.Person()

	for i := 1; i <= total; i++ {
		book := model.Book{
			Title:        fakeLorem.Sentence(2),
			Sheet:        uint(faker.IntBetween(100, 250)),
			Introduction: fakeLorem.Paragraph(2),
			DateOffIssue: time.Date(2004, time.November, 12, 0, 0, 0, 0, time.UTC),
			Author:       []string{fakePerson.FirstName(), fakePerson.FirstName()},
		}

		books = append(books, book)
	}

	return books
}
