package main

import (
	"github.com/pkg/errors"
	http "net/http"
)

const (
	// *** Status for successfull requests (2xx)  ***

	// StatusDefault : Default status for all requests
	StatusDefault = 200
	// StatusCreated : Status for requests where object is created successfully
	StatusCreated = 201
	// StatusUpdated : Status where requests succesfully updated the data
	StatusUpdated = 200
	// StatusDeleted : Status where request succesfully deleted the data
	StatusDeleted = 200

	// *** Error Status ***
	// * Client Side 4xx *

	// StatusClientError : Default status for incorrect requests
	StatusClientError = 400
	// StatusObjectNotFound : Status when an object requested could not be found
	//  e.g. a medicine or patient does not exist.
	StatusObjectNotFound = 400
	// StatusUnauthorized : Status when authorisation is needed but failed
	StatusUnauthorized = 401
	// StatusInvalidJSON : Status when errored during decoding of the json body
	StatusInvalidJSON = 499
	// StatusInvalidLanguage : Langauge specified is not supported / invalid
	StatusInvalidLanguage = 498

	// * Server Side (5xx) *

	// StatusServerError : Default status if something went wrong in the server
	StatusServerError = 500
	// StatusDatabaseError : All errors related to the database (58x)
	StatusDatabaseError = 580
	// StatusDatabaseConstraintViolation : Status for constraint violations in the database
	StatusDatabaseConstraintViolation = 581
	// StatusFailedOperation : ?
	StatusFailedOperation = 599
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
	a.StatusCode = StatusDefault
	a.Error = nil
}

// Sets the error and the standard error message
func (a *APIResponse) setError(err error, errMessages ...string) {
	a.setErrorAndStatus(StatusServerError, err, errMessages...)
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
	a.setResponseAndStatus(StatusDefault, data)
}

func (a *APIResponse) setResponseAndStatus(status int, data interface{}) {
	a.StatusCode = status
	a.Data = data
}

func (a *APIResponse) setStatus(status int) {
	a.StatusCode = status
}
