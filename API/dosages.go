package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	http "net/http"
)

// CREATE
func createDosage(r *http.Request, ar *APIResponse) {
	dosage := Dosage{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&dosage)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode JSON.")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	var medicineID int
	err = tx.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		dosage.Medicine.Name).Scan(&medicineID)
	if err != nil {
		if err == sql.ErrNoRows {
			ar.setErrorAndStatus(http.StatusNotFound, err, "Unknown medicine.")
		} else {
			ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to execute query.")
		}
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	_, err = tx.Exec(`INSERT INTO Dosages (patient_id, medicine_id, amount,
 intake_interval_start, intake_interval_end) 
                          VALUES (?, ?, ?, ?, ?)`,
		patientID, medicineID, dosage.NumberOfPills, dosage.IntakeIntervalStart, dosage.IntakeIntervalEnd)
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
