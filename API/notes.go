package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
func addNote(r *http.Request, ar *APIResponse) {
	// verify patient
	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to start new transaction")
		return
	}

	vars := mux.Vars(r)
	patientID := vars["id"]
	_, err = tx.Exec(
		`INSERT INTO Notes (patient_id, question, day) VALUES (?, ?, ?)`,
		patientID, note.Note, note.CreatedAt)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to insert note into the database")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusCreated
}

// RETRIEVE
// Possible to also add a time interval?
// Or all 'untreated' notes
func getNotes(r *http.Request, ar *APIResponse) {

	vars := mux.Vars(r)
	patientID := vars["id"]

	rows, err := db.Query(`SELECT question, day FROM Notes WHERE patient_id = ?`, patientID)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error =  errors.Wrap(err, "Unexpected error during query")
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note, date string
		err = rows.Scan(&note, &date)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error =  errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		notes = append(notes, Note{note, date})
	}
	if err = rows.Err(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error =  errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	ar.Data = notes
}
