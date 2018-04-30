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
)

var (
	db *sql.DB
)

func main() {
	var err error
	rootpasswd, dbname, listenLocation := "pass", "database", "localhost:8080" // just some values
	fmt.Scanf("%s", &rootpasswd)
	fmt.Scanf("%s", &dbname)
	fmt.Scanf("%s", &listenLocation)
	db, err = sql.Open("mysql", "root:"+rootpasswd+"@/"+dbname)

	if err != nil {
		log.Printf("encountered error while connecting to database: %v", err)
	}

	log.Printf("Connected to database '%s', and listening on '%s'...", dbname, listenLocation)
	router := mux.NewRouter()

	// GET Requests for Retrieving
	getRouter := router.Methods("GET").Subrouter()
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(authWrapper(getDosages)))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(authWrapper(getNotes)))
	getRouter.Handle("/api/general/videos/topics/{topic}", handlerWrapper(getVideoByTopic))
	getRouter.Handle("/api/general/videos/topics", handlerWrapper(getTopics))
	getRouter.Handle("/api/general/faq", handlerWrapper(getFAQs))

	// POST Requests for Updating
	postRouter := router.Methods("POST").Subrouter()
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(authWrapper(modifyPatient)))
	postRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(authWrapper(modifyPhysician)))
	postRouter.Handle("/api/accounts/login", handlerWrapper(login))

	// PUT Requests for Creating
	putRouter := router.Methods("PUT").Subrouter()
	putRouter.Handle("/api/accounts/patients", handlerWrapper(pushPatient))
	putRouter.Handle("/api/accounts/physicians", handlerWrapper(pushPhysician))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(pushDosage))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(addNote))
	putRouter.Handle("/api/general/videos", handlerWrapper(addVideo))

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

func handlerWrapper(handler func(r *http.Request, ar *APIResponse)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ar := APIResponse{nil, 200, nil}
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
