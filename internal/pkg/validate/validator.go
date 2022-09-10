package validate

import (
	"auth/ent/schema"
	"auth/internal/pkg/slices"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	customValidations = map[string]func(fl validator.FieldLevel) bool{
		`user_type`: func(fl validator.FieldLevel) bool {
			return slices.Includes(fl.Field().String(), schema.UserTypes)
		},
	}
)

func Init() error {
	validate = validator.New()

	for tag, fn := range customValidations {
		err := validate.RegisterValidation(tag, fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func Struct(s any) error {
	if validate == nil {
		if err := Init(); err != nil {
			return err
		}
	}
	return validate.Struct(s)
}

func IsInvalidValidationError(err error) bool {
	if err == nil {
		return false
	}
	_, is := err.(*validator.InvalidValidationError)
	return is
}

func AsValidationErrors(err error) validator.ValidationErrors {
	return err.(validator.ValidationErrors)
}

func Default(s any) error {
	err := Struct(s)
	if err != nil {
		if IsInvalidValidationError(err) {
			return err
		}
		for _, err := range AsValidationErrors(err) {
			return err // TODO Make better with err.Field() and others
		}
	}
	return nil
}
