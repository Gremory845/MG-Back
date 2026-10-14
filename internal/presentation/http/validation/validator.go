package validation

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(s interface{}) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs []ValidationError
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, ValidationError{
			Field:   e.Field(),
			Message: humanMessage(e),
		})
	}
	return errs
}

func humanMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "es obligatorio"
	case "email":
		return "debe ser un correo válido"
	case "min":
		return "es demasiado corto"
	case "max":
		return "es demasiado largo"
	default:
		return "es inválido"
	}
}
