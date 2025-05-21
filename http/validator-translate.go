package http

import (
	"github.com/go-playground/validator/v10"
	"strings"
)

func FormatValidation(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {

		errors = append(errors, e.Translate(Trans))
	}
	return errors
}
func Format(err error) (interface{}, bool) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil, false
	}
	// validator.ValidationErrors类型错误则进行翻译
	mErr := errs.Translate(Trans)
	return RemoveTopStruct(mErr), true
}
func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
