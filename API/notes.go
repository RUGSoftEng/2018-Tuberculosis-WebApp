package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	http "net/http"
	"time"
)

// CREATE
func createNote(r *http.Request, ar *APIResponse) {
	// verify patient
	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during JSON decoding.")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start new transaction.")
		return
	}

	vars := mux.Vars(r)
	patientID := vars["id"]
	_, err = tx.Exec(
		`INSERT INTO Notes (patient_id, question, day) VALUES (?, ?, ?)`,
		patientID, note.Note, note.CreatedAt)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to insert note into the database.")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
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
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	notes := []NoteReturn{}
	for rows.Next() {
		var note, date string
		var id int
		err = rows.Scan(&id, &note, &date)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		notes = append(notes, NoteReturn{id, note, date})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}

	ar.setResponse(notes)
}

//DELETE

func deleteNote(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	patientID := vars["id"]
	noteID := vars["note_id"]
	_, err := db.Exec("DELETE FROM Notes WHERE id=? and patient_id=?", noteID, patientID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	ar.StatusCode = http.StatusOK

}

//POST

func modifyNote(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	patientID := vars["id"]
	noteID := vars["note_id"]
	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start new transaction.")
		return
	}

	_, err = tx.Exec(`UPDATE Notes SET
                          question = ?,
                          day = ?
                          WHERE id = ?`, note.Note, time.Now(), noteID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to insert note into the database.")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusCreated

}
