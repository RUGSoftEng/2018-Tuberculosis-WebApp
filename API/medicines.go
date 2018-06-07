package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	http "net/http"
)

// CREATE
func createMedicine(r *http.Request, ar *APIResponse) {
	medicine := Medicine{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&medicine)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode JSON.")
		return
	}
	tx, err := db.Begin()

	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`INSERT INTO Medicines(med_name) VALUES (?)`, medicine.Name)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}
	ar.StatusCode = http.StatusCreated
}

// DELETE
func deleteMedicine(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	medID := vars["id"]

	tx, err := db.Begin()

	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`DELETE FROM Medicines WHERE id=?`, medID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusOK
}
