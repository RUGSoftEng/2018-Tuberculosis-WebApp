package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	http "net/http"
)

// CREATE
func createMedicine(r *http.Request, ar *APIResponse) {
	medicine := Medicine{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&medicine)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`INSERT INTO Medicines(med_name) VALUES (?)`, medicine.Name)
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

// RETRIEVE: ?

// UPDATE: ?

// DELETE
func deleteMedicine(r *http.Request, ar *APIResponse) {
	medID, err := getURLVariable(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`DELETE FROM Medicines WHERE id=?`, medID)
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
