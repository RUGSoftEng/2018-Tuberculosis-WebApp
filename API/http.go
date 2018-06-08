package main

import (
	"database/sql"
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
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(authWrapper(retrieveScheduledDosages)))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(authWrapper(retrieveDosages)))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(authWrapper(retrieveNotes)))
	getRouter.Handle("/api/accounts/patients/{id:[0-9]+}/retrieveByID", handlerWrapper(retrieveByID))
	getRouter.Handle("/api/accounts/patients/{username}/retrieveByUsername", handlerWrapper(retrieveByUsername))
	getRouter.Handle("/api/accounts/physicians/{id:[0-9]+}/retrievePyByID", handlerWrapper(retrievePyByID))
	getRouter.Handle("/api/accounts/physicians/{username}/retrievePyByUsername", handlerWrapper(retrievePyByUsername))
	getRouter.Handle("/api/general/videos/topics/{topic}", handlerWrapper(retrieveVideoByTopic))
	getRouter.Handle("/api/general/videos/topics", handlerWrapper(retrieveTopics))
	getRouter.Handle("/api/general/faqs", handlerWrapper(retrieveFAQs))
	getRouter.Handle("/api/general/physicians/{id:[0-9]+}/retrieve", handlerWrapper(authWrapper(retrievePatients)))

	// POST Requests for Updating
	postRouter := router.Methods("POST").Subrouter()
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(authWrapper(updatePatient)))
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(authWrapper(updateScheduledDosage)))
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes/{note_id:[0-9]+}", handlerWrapper(authWrapper(updateNote)))
	postRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(updateDosage))
	postRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(authWrapper(updatePhysician)))
	postRouter.Handle("/api/accounts/login", handlerWrapper(login))
	postRouter.Handle("/api/admin/faqs", handlerWrapper(updateFAQ))
	postRouter.Handle("/api/admin/videos", handlerWrapper(updateVideo))
	postRouter.Handle("/api/admin/videos/quizzes", handlerWrapper(updateQuiz))

	// PUT Requests for Creating
	putRouter := router.Methods("PUT").Subrouter()
	putRouter.Handle("/api/accounts/patients", handlerWrapper(createPatient))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(createDosage))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(createScheduledDosages))
	putRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(createNote))
	putRouter.Handle("/api/accounts/physicians", handlerWrapper(createPhysician))
	putRouter.Handle("/api/admin/videos", handlerWrapper(createVideo))
	putRouter.Handle("/api/admin/faqs", handlerWrapper(createFAQ))
	putRouter.Handle("/api/admin/medicines", handlerWrapper(createMedicine))
	putRouter.Handle("/api/admin/videos/quizzes", handlerWrapper(createQuiz))

	// DELETE Requests for Deleting
	deleteRouter := router.Methods("DELETE").Subrouter()
	deleteRouter.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(deletePatient))
	deleteRouter.Handle("/api/accounts/patients/{id:[0-9]+}/notes/{note_id:[0-9]+}", handlerWrapper(authWrapper(deleteNote)))
	deleteRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages/scheduled", handlerWrapper(deleteScheduledDosage))
	deleteRouter.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(deleteDosage))
	deleteRouter.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(deletePhysician))
	deleteRouter.Handle("/api/admin/medicines/{id:[0-9]+}", handlerWrapper(deleteMedicine))
	deleteRouter.Handle("/api/admin/faqs", handlerWrapper(deleteFAQ))
	deleteRouter.Handle("/api/admin/videos", handlerWrapper(deleteVideo))
	deleteRouter.Handle("/api/admin/videos/quizzes", handlerWrapper(deleteQuiz))

	// Starting the router
	err = http.ListenAndServe(listenLocation, router)
	if err != nil {
		log.Fatal(err)
	}
}
