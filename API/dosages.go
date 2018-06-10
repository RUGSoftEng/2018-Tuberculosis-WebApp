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

// RETRIEVE
func retrieveDosages(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	dosages := []Dosage{}
	rows, err := db.Query(`SELECT amount, med_name, intake_interval_start, intake_interval_end
                          FROM Dosages JOIN Medicines 
                             ON Dosages.medicine_id = Medicines.id 
                          WHERE patient_id = ?`, patientID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	for rows.Next() {
		var amount int
		var medicine, intakeIntervalStart, intakeIntervalEnd string
		err = rows.Scan(&amount, &medicine, &intakeIntervalStart, &intakeIntervalEnd)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, Dosage{
			intakeIntervalStart,
			intakeIntervalEnd,
			amount,
			Medicine{medicine}})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}
	ar.setResponse(dosages)
}

// UPDATE
func updateDosage(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	dosage := UpdateDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&dosage)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Unable to decode given json data")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Unable to start database transaction")
		return
	}

	oldMedicineID, err := queryMedicineID(dosage.OldDosage.Medicine)
	if err != nil {
		ar.setError(err, "Error during Medicine Query")
		return
	}
	newMedicineID, err := queryMedicineID(dosage.NewDosage.Medicine)
	if err != nil {
		ar.setError(err, "Error during Medicine Query")
		return
	}

	_, err = tx.Exec(`UPDATE Dosages SET medicine_id = ?, amount = ?, 
                          intake_interval_start = ?, intake_interval_end = ? 
                          WHERE patient_id = ? AND medicine_id = ?`,
		newMedicineID, dosage.NewDosage.NumberOfPills,
		dosage.NewDosage.IntakeIntervalStart,
		dosage.NewDosage.IntakeIntervalEnd,
		patientID, oldMedicineID)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Something went wrong during SQL Update query")
		return
	}

	if err := tx.Commit(); err != nil {
		ar.setError(err, "Unable to commit changes to the database")
		return
	}
}

// DELETE
func deleteDosage(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	dosage := Dosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&dosage)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Unable to decode given json data")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Unable to start database transaction")
		return
	}

	medicineID, err := queryMedicineID(dosage.Medicine)
	if err != nil {
		ar.setError(err, "Error during Medicine Query")
		return
	}

	_, err = tx.Exec(`DELETE FROM Dosages WHERE patient_id = ? AND medicine_id = ?`,
		patientID, medicineID)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Something went wrong during SQL Update query")
		return
	}

	if err := tx.Commit(); err != nil {
		ar.setError(err, "Unable to commit changes to the database")
		return
	}
}
