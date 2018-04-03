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

func exampleHandler(r *http.Request, responseChan chan []byte, errorChan chan error) {
	ID := 0
	apiToken := r.Header.Get("api_token")

	// This is a join example for a patient call, change to physician it is a call only a physician can make
	// remove join part if it is a call able for both
	err := db.QueryRow(`SELECT id
			   FROM Patients AS pa 
			   INNER JOIN Accounts AS acc 
			   ON pa.id = acc.id  
			   WHERE acc.api_token = ?`,
		apiToken).Scan(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorChan <- errors.Wrap(err, "no valid login credentials")
			return
		}
		errorChan <- errors.Wrap(err, "encountered error during query")
		return
	}

	// if you are going to insert multiple things in the database do this using a transaction.
	// see insertPatient

	// do your own querries,
	// if you encounter a "err != nil" send it to the errorChan in the above matter
	// if all goed well, marshal your results and sen them to responseChan

	// End for a get function
	// responseChan <- "your marshalled data"

	// End for a succesfull push or put function
	// errorChan <- nil
}
