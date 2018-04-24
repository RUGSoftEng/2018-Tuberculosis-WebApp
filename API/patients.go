package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
// expects a json file containing the new patient and a url encoded physician token
func pushPatient(r *http.Request, ar *APIResponse) {
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Failed to decode incoming JSON")
		return
	}

	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to hash password")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to start transaction")
		return
	}

	role := "patient"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, patient.Name, patient.Username, patient.Password, role)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	var physicianID int
	creationToken := r.URL.Query().Get("token")
	err = tx.QueryRow(`SELECT id FROM Physicians WHERE token=?`, creationToken).Scan(&physicianID)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	_, err = tx.Exec(`INSERT INTO Patients VALUES(?,?)`, id, physicianID)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	err = tx.Commit()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusCreated
}

// UPDATE
func modifyPatient(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Failed to decode incoming JSON")
		return
	}
	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Hashing failed")
		return
	}

	// Using a transaction because I don't know whether we are going to have to add
	// query for a possible change of physician (or how to do that)
	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "failed to start transaction")
		return
	}
	_, err = tx.Exec(`UPDATE Accounts SET 
                 name = ?,
                 pass_hash = ?
                 WHERE id = ?`, patient.Name, patient.Password, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
}

// DELETE
func deletePatient(r *http.Request, ar *APIResponse) {
	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error =  errors.Wrap(err, "failed to start transaction")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	_, err = tx.Exec(`DELETE FROM Notes WHERE patient_id=?`, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	// Retrieve all dosage identifiers
	rows, err := tx.Query(`SELECT id FROM Dosages
                               WHERE patient_id = ?`, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	var dosageIDs []int
	for rows.Next() {
		var dosageID int
		err = rows.Scan(&id)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errorWithRollback(err, tx)
			return
		}
		dosageIDs = append(dosageIDs, dosageID)
	}
	if rows.Err() != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	// Delete all specific scheduled dosages attached to the patient
	for _, dosageID := range dosageIDs {
		_, err = tx.Exec(`DELETE FROM SchedulesDosages WHERE dosage=?`, dosageID)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errorWithRollback(err, tx)
			return
		}
	}

	_, err = tx.Exec(`DELETE FROM Dosages WHERE patient_id=?`, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	_, err = tx.Exec(`DELETE FROM Patients WHERE id=?`, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}
	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, id)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errorWithRollback(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
}
