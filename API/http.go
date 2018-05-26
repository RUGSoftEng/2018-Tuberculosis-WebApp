package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	http "net/http"
	"reflect"
	"runtime"
)

var (
	db *sql.DB
)

// Custom made error codes
const (
	StatusFailedOperation             = 599
	StatusDatabaseConstraintViolation = 598
)

func main() {
	var err error
	var dbUser, dbUserPassword, dbName, listenLocation string

	_, err = fmt.Scanf("%s", &dbUser)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Scanning of Database User failed").Error())
	}
	_, err = fmt.Scanf("%s", &dbUserPassword)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Scanning of Database User's Password failed").Error())
	}
	_, err = fmt.Scanf("%s", &dbName)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Scanning of Database Name failed").Error())
	}
	_, err = fmt.Scanf("%s", &listenLocation)
	if err != nil {
		log.Fatal(errors.Wrap(err, "Scanning of Listen Location failed").Error())
	}

	db, err = sql.Open("mysql", dbUser+":"+dbUserPassword+"@/"+dbName)
	if err != nil {
		log.Printf("encountered error while connecting to database: %v", err)
	}

	log.Printf("Connected to database '%s', and listening on '%s'...", dbName, listenLocation)
	router := mux.NewRouter()

	// GET Requests for Retrieving
	getRouter := router.Methods("GET").Subrouter()
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(authWrapper(getDosages)))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(authWrapper(getNotes)))
	getRouter.Handle("/api/general/videos/topics/{topic}", handlerWrapper(getVideoByTopic))
	getRouter.Handle("/api/general/videos/topics", handlerWrapper(getTopics))
	getRouter.Handle("/api/general/faq", handlerWrapper(getFAQs))
	getRouter.Handle("/api/general/physicians/{id:[0-9]+}/retrieve", handlerWrapper(authWrapper(getPatients)))

	// POST Requests for Updating
	postRouter := router.Methods("POST").Subrouter()
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(authWrapper(updatePatient)))
	postRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(authWrapper(updatePhysician)))
	postRouter.Handle("/api/accounts/login", handlerWrapper(login))
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(authWrapper(updateScheduledDosage)))

	// PUT Requests for Creating
	putRouter := router.Methods("PUT").Subrouter()
	putRouter.Handle("/api/accounts/patients", handlerWrapper(createPatient))
	putRouter.Handle("/api/accounts/physicians", handlerWrapper(createPhysician))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(createDosage))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(createScheduledDosages))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(createNote))
	putRouter.Handle("/api/general/videos", handlerWrapper(createVideo))
	putRouter.Handle("/api/admin/faq", handlerWrapper(createFAQ))

	// DELETE Requests for Deleting
	deleteRouter := router.Methods("DELETE").Subrouter()
	deleteRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(deletePatient))
	deleteRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(deletePhysician))

	// Starting the router
	err = http.ListenAndServe(listenLocation, router)
	if err != nil {
		log.Fatal(err)
	}
}

// APIResponse : Type used by the Response Channel
// in the handlerWrapper (does not need json tags)
type APIResponse struct {
	Data       interface{}
	StatusCode int
	Error      error
}

func (a *APIResponse) setError(err error, errMessage string) {
	a.setErrorAndStatus(http.StatusInternalServerError, err, errMessage)
}

func (a *APIResponse) setErrorAndStatus(status int, err error, errMessage string) {
	a.StatusCode = status
	a.Error = errors.Wrap(err, errMessage)
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

func handlerWrapper(handler func(r *http.Request, ar *APIResponse)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ar := APIResponse{nil, 200, nil}
		funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.Printf("New request:\n |url:  %s\n |func: %s\n", r.URL, funcName)
		handler(r, &ar)

		if ar.Error != nil {
			log.Printf("Server error: %v", ar.Error)
			if ar.StatusCode == http.StatusInternalServerError {
				ar.Error = errors.New(http.StatusText(http.StatusInternalServerError))
			}
			http.Error(w, ar.Error.Error(), ar.StatusCode)
			return
		}

		if ar.Data == nil {
			w.WriteHeader(ar.StatusCode)
			return
		}

		jsonData, err := json.Marshal(ar.Data)
		if err != nil {
			err := errors.Wrap(err, "Error during JSON Decoding")
			log.Printf("Error marshalling response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
