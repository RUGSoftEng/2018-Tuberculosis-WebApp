package main

import (
	"github.com/pkg/errors"
	http "net/http"
)

// APIResponse : Type used by the Response Channel
// in the handlerWrapper (does not need json tags)
type APIResponse struct {
	Data       interface{}
	StatusCode int
	Error      error
}

func (a *APIResponse) init() {
	a.Data = nil
	a.StatusCode = http.StatusOK
	a.Error = nil
}

// Sets the error and the standard error message
func (a *APIResponse) setError(err error, errMessages ...string) {
	a.setErrorAndStatus(http.StatusInternalServerError, err, errMessages...)
}

// Sets the error and the given status. If given extra error messages,
// it wraps them with the original error.
func (a *APIResponse) setErrorAndStatus(status int, err error, errMessages ...string) {
	for _, errMessage := range errMessages {
		err = errors.Wrap(err, errMessage)
	}
	a.StatusCode = status
	a.Error = err
}

func (a *APIResponse) setResponse(data interface{}) {
	a.setResponseAndStatus(http.StatusOK, data)
}

func (a *APIResponse) setResponseAndStatus(status int, data interface{}) {
	a.StatusCode = status
	a.Data = data
}

func (a *APIResponse) setStatus(status int) {
	a.StatusCode = status
}
