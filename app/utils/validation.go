package utils

import (
	"employee-management/app/models"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

func ValidateRegisterUser(user *models.User) map[string]string {
	validate := *validator.New()
	err := validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		errorList := make(map[string]string)

		for _, e := range errors {
			var errMsg string
			field, _ := reflect.TypeOf(*user).FieldByName(e.StructField())
			fieldName := field.Tag.Get("json")

			switch e.Tag() {
			case "required":
				errMsg = fmt.Sprintf("%s tidak boleh kosong", fieldName)
			case "email":
				errMsg = fmt.Sprintf("ini bukan %s yang valid", fieldName)
			case "oneof":
				errMsg = fmt.Sprintf("%s harus berupa %s", fieldName, e.Param())
			case "min":
				errMsg = fmt.Sprintf("%s minimal %s karakter", fieldName, e.Param())
			case "max":
				errMsg = fmt.Sprintf("%s minimal %s karakter", fieldName, e.Param())
			}

			errorList[fieldName] = errMsg
		}

		return errorList
	}

	return nil
}
