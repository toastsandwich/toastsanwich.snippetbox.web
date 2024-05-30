package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func NewValidator() Validator {
	return Validator{
		FieldErrors:    make(map[string]string),
		NonFieldErrors: make([]string, 0),
	}
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, val string) {
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = val
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, val string) {
	if !ok {
		v.AddFieldError(key, val)
	}
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func MinChar(value string, lim int) bool {
	return utf8.RuneCountInString(value) >= lim
}

func NotABlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxCharLimit(value string, lim int) bool {
	return utf8.RuneCountInString(value) <= lim
}

func PermittedInts(value int, permittedvalues ...int) bool {
	for _, val := range permittedvalues {
		if val == value {
			return true
		}
	}
	return false
}
