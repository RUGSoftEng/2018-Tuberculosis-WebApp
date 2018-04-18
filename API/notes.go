package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
	"log"
)

// CREATE
func addNote(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	// verify patient
	vars := mux.Vars(r)
	patientID := vars["id"]

	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}

	trans, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start new transaction")
		return
	}
	_, err = trans.Exec(
		`INSERT INTO Notes (patient_id, question, day) VALUES (?, ?, ?)`,
		patientID, note.Note, note.CreatedAt)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to insert note into the database")
		return
	}

	if err = trans.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
	responseChan <- APIResponse{nil, http.StatusCreated}
}

// RETRIEVE
// Possible to also add a time interval?
// Or all 'untreated' notes
func getNotes(r *http.Request, responseChan chan APIResponse, errorChan chan error) {

	if !authenticate(r, responseChan, errorChan) {
		log.Println("You are not authenticated")
		return
	}

	vars := mux.Vars(r)
	patientID := vars["id"]

	rows, err := db.Query(`SELECT question, day FROM Notes WHERE patient_id = ?`, patientID)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during query")
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note, date string
		err = rows.Scan(&note, &date)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		notes = append(notes, Note{note, date})
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}
	responseChan <- APIResponse{notes, http.StatusOK}
}
