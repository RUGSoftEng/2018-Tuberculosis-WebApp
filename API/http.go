package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	http "net/http"
	"time"
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
	router.Handle("/api/your extension", handlerWrapper(exampleHandler))

	// GET Requests for Retrieving
	getRouter := router.Methods("GET").Subrouter()
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(getDosages))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(getNotes))
	getRouter.Handle("/api/general/videos/topics/{topic}", handlerWrapper(getVideoByTopic))
	getRouter.Handle("/api/general/videos/topics", handlerWrapper(getTopics))

	// POST Requests for Updating
	postRouter := router.Methods("POST").Subrouter()
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(modifyPatient))
	postRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(modifyPhysician))
	postRouter.Handle("/api/accounts/login", handlerWrapper(login))
	postRouter.Handle("/api/accounts/authenticate", handlerWrapper(authenticate))

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
	http.ListenAndServe(listenLocation, router)
}

func handlerWrapper(handler func(r *http.Request, responseChan chan []byte, errorChan chan error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseChan := make(chan []byte)
		errorChan := make(chan error)

		go handler(r, responseChan, errorChan)

		time.After(2 * time.Second)

		select {
		case body := <-responseChan:
			w.Header().Set("Content-Type", "application/json")
			w.Write(body)
		case err := <-errorChan:
			if err != nil {
				log.Printf("Server error: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		case <-time.After(5 * time.Second):
			log.Printf("Response timeout")
		}
		return
	})
}
