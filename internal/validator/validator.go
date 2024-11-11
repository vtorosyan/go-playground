package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}")

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key string, msg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = msg
	}
}

func (v *Validator) CheckField(ok bool, key string, msg string) {
	if !ok {
		v.AddFieldError(key, msg)
	}
}

func NotBlank(val string) bool {
	return strings.TrimSpace(val) != ""
}

func MaxChars(val string, n int) bool {
	return utf8.RuneCountInString(val) <= n
}

func PermittedValue[T comparable](val T, permitted ...T) bool {
	return slices.Contains(permitted, val)
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
