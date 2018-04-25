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
func pushPatient(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to decode incoming JSON")
		return
	}
	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to hash password")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start transaction")
		return
	}
	role := "patient"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, patient.Name, patient.Username, patient.Password, role)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}
	id, err := result.LastInsertId()
	creationToken := r.URL.Query().Get("token")
	var physicianID int
	err = tx.QueryRow(`SELECT id FROM Physicians WHERE token=?`, creationToken).Scan(&physicianID)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}
	_, err = tx.Exec(`INSERT INTO Patients VALUES(?,?)`, id, physicianID)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}
	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
	responseChan <- APIResponse{nil, http.StatusCreated}
}

// UPDATE
func modifyPatient(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	vars := mux.Vars(r)
	id := vars["id"]
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		errorChan <- err
		return
	}
	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		errorChan <- errors.Wrap(err, "Hashing failed")
		return
	}

	// Using a transaction because I don't know whether we are going to have to add
	// query for a possible change of physician (or how to do that)
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	_, err = tx.Exec(`UPDATE Accounts SET 
                 name = ?,
                 pass_hash = ?
                 WHERE id = ?`, patient.Name, patient.Password, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}

	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
	responseChan <- APIResponse{nil, http.StatusOK}
}

// DELETE
func deletePatient(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	vars := mux.Vars(r)
	id := vars["id"]
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	_, err = tx.Exec(`DELETE FROM Notes WHERE patient_id=?`, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}

	// Retrieve all dosage identifiers
	rows, err := tx.Query(`SELECT id FROM Dosages
                               WHERE patient_id = ?`, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}
	var dosageIDs []int
	for rows.Next() {
		var dosageID int
		err = rows.Scan(&id)
		if err != nil {
			errorChan <- err
			err = tx.Rollback()
			if err != nil {
				errorChan <- errors.Wrap(err, "Rollback failed")
			}
			return
		}
		dosageIDs = append(dosageIDs, dosageID)
	}
	if rows.Err() != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}

	// Delete all specific scheduled dosages attached to the patient
	for _, dosageID := range dosageIDs {
		_, err = tx.Exec(`DELETE FROM SchedulesDosages WHERE dosage=?`, dosageID)
		if err != nil {
			errorChan <- err
			err = tx.Rollback()
			if err != nil {
				errorChan <- errors.Wrap(err, "Rollback failed")
			}
			return
		}
	}

	_, err = tx.Exec(`DELETE FROM Dosages WHERE patient_id=?`, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}

	_, err = tx.Exec(`DELETE FROM Patients WHERE id=?`, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}
	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, id)
	if err != nil {
		errorChan <- err
		err = tx.Rollback()
		if err != nil {
			errorChan <- errors.Wrap(err, "Rollback failed")
		}
		return
	}

	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
	responseChan <- APIResponse{nil, http.StatusOK}
}
