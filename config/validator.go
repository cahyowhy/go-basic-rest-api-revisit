package config

import (
	"sync"

	"github.com/go-playground/validator"
)

var structValidator *validator.Validate
var onceApp sync.Once

func GetValidator() *validator.Validate {
	onceApp.Do(func() {
		structValidator = validator.New()
	})

	return structValidator
}
