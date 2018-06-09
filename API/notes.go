package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	http "net/http"
)

// CREATE
func createNote(r *http.Request, ar *APIResponse) {
	// verify patient
	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(
		`INSERT INTO Notes (patient_id, question, day) VALUES (?, ?, ?)`,
		patientID, note.Note, note.CreatedAt)
	if err != nil {
		ar.setErrorDBInsert(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusCreated)
}

// RETRIEVE
func retrieveNotes(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	rows, err := db.Query(`SELECT id, question, day FROM Notes WHERE patient_id = ?`, patientID)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	notes := []NoteReturn{}
	for rows.Next() {
		var note, date string
		var id int
		err = rows.Scan(&id, &note, &date)
		if err != nil {
			ar.setErrorDBScan(err)
			return
		}
		notes = append(notes, NoteReturn{id, note, date})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorDBAfter(err)
		return
	}

	ar.setResponse(notes)
}

// UPDATE
func updateNote(r *http.Request, ar *APIResponse) {
	noteID, err := getURLVariable(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	note := Note{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&note)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`UPDATE Notes SET
                          question = ?,
                          day = ?
                          WHERE id = ?`, note.Note, note.CreatedAt, noteID)
	if err != nil {
		ar.setErrorDBUpdate(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusUpdated)
}

// DELETE
func deleteNote(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	noteID, err := getURLVariable(r, "note_id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = db.Exec("DELETE FROM Notes WHERE id=? and patient_id=?", noteID, patientID)
	if err != nil {
		ar.setErrorDBDelete(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusDeleted)
}
