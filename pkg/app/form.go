package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	val "github.com/go-playground/validator/v10"
	"go-gin/pkg/translations"
	"strings"
)

type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		var valErr val.ValidationErrors
		ok := errors.As(err, &valErr)
		if !ok {
			return false, errs
		}

		for key, value := range valErr.Translate(translations.Trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
