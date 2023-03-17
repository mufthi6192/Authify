package validation

import (
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

func MinimumChar(field validator.FieldLevel) bool {

	value := field.Field().String()
	param := field.Param()

	minChar, err := strconv.Atoi(param)

	if err != nil {
		panic("Failed to convert")
	}

	countSpace := strings.Count(value, " ")
	countChar := len(value)
	totalChar := countChar - countSpace

	if totalChar < minChar {
		return false
	}

	return true
}
