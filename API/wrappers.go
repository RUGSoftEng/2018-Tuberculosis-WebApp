package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	http "net/http"
	"reflect"
	"runtime"
)

// Wraps each response function to simplify logging, error handling, response writing and data encoding
func handlerWrapper(handler func(r *http.Request, ar *APIResponse)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create API Response with standard variables
		var ar APIResponse
		ar.init()

		// Log incomming request: Requested URL and Function called
		funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.Printf("New request:\n |url:  %s\n |func: %s\n", r.URL, funcName)

		// Call the function
		handler(r, &ar)

		// Handle possible errors
		if ar.Error != nil {
			log.Printf("Server error: %v", ar.Error)
			if ar.StatusCode == http.StatusInternalServerError {
				ar.Error = errors.New(http.StatusText(http.StatusInternalServerError))
			}
			// If error is different than InternalServerError,
			// Give the error to the user
			http.Error(w, ar.Error.Error(), ar.StatusCode)
			return
		}

		// If there is no data sent back, write status code and return
		if ar.Data == nil {
			w.WriteHeader(ar.StatusCode)
			return
		}

		// Marshal outgoing data to JSON
		jsonData, err := json.Marshal(ar.Data)
		if err != nil {
			err := errors.Wrap(err, "Error during JSON Encoding")
			log.Printf("Error marshalling response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(ar.StatusCode)
		_, err = w.Write(jsonData) //returns an integer, not sure what it's used for
		if err != nil {
			log.Printf("Error sending response to request: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		return
	})
}

func authWrapper(handler func(r *http.Request, ar *APIResponse)) func(*http.Request, *APIResponse) {
	return func(r *http.Request, ar *APIResponse) {

		authenticate(r, ar)
		if ar.Error != nil {
			return
		}
		handler(r, ar)
	}
}
