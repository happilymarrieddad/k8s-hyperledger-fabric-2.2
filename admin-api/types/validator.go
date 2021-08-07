package types

import (
	"github.com/go-playground/validator/v10"
)

var (
	validCop *validator.Validate = validator.New()
)

type ValidatingEntity interface {
	Validate() error
}

func Validate(entity interface{}) (err error) {
	ve, ok := entity.(ValidatingEntity)
	if ok {
		err = ve.Validate()
	} else {
		err = validCop.Struct(entity)
	}

	return
}
