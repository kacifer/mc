package mgin

import (
	"fmt"
	"strings"
)

type ErrorDetails map[string]any

type E struct {
	// Code is the error code, if not specified, it will be the same as the http status code.
	Code int `json:"code"`

	// Message is the error message.
	Message string `json:"message"`

	// Details can be used to provide more details about the error, such as the invalid form fields.
	Details ErrorDetails `json:"details,omitempty"`
}

func (e *E) Error() string {
	var details []string
	if e.Details != nil {
		for k, v := range e.Details {
			details = append(details, fmt.Sprintf("%s: %v", k, v))
		}
	}
	return fmt.Sprintf("[%d] %s%s", e.Code, e.Message, strings.Join(details, ", "))
}
