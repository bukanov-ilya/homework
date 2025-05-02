package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

func (e *MultiError) Error() string {
	if len(e.Errors) == 0 {
		return ""
	}

	finalErr := fmt.Sprintf("%d errors occured:\n", len(e.Errors))

	for _, err := range e.Errors {
		finalErr += fmt.Sprintf("\t* %s", err)
	}

	finalErr += "\n"

	return finalErr
}

func Append(err error, errs ...error) *MultiError {
	var multiError *MultiError

	if errors.As(err, &multiError) {
		multiError.Errors = append(multiError.Errors, errs...)

		return multiError
	}

	allErrors := make([]error, 0, len(errs)+1)
	if err != nil {
		allErrors = append(allErrors, err)
	}

	allErrors = append(allErrors, errs...)

	return &MultiError{
		Errors: allErrors,
	}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
