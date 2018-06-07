package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	http "net/http"
	"time"
)

// CREATE
func createScheduledDosages(r *http.Request, ar *APIResponse) {
	// Scan Patient ID
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	// Decode Request Body to JSON
	in := InputScheduledDosage{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&in)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode JSON.")
		return
	}

	// Start Database Transaction
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	// Query Medicine ID
	var medicineID int
	err = tx.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`, in.Medicine.Name).Scan(&medicineID)
	if err != nil {
		if err == sql.ErrNoRows {
			ar.setErrorAndStatus(http.StatusNotFound, err, "Unknown medicine.")
		} else {
			ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to execute query.")
		}
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	// Query Dosage ID
	var dosageID int
	err = tx.QueryRow(`SELECT id FROM Dosages WHERE medicine_id = ? AND patient_id = ?`, medicineID, patientID).Scan(&dosageID)
	if err != nil {
		if err == sql.ErrNoRows {
			ar.setErrorAndStatus(http.StatusNotFound, err, "Unknown dosage.")
		} else {
			ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to execute query.")
		}
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	// Add a ScheduledDosage for each given day
	for _, day := range in.Days {
		_, err := tx.Exec(`INSERT INTO ScheduledDosages VALUES (?, ?, ?)`, dosageID, day, false)
		if err != nil {
			ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
			return
		}
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}
	ar.StatusCode = http.StatusCreated
}

// RETRIEVE
func retrieveScheduledDosages(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	from := r.URL.Query().Get("from")
	until := r.URL.Query().Get("until")
	const dform = "2006-01-02" // specifies YYYY-MM-DD format
	startDate, err := time.Parse(dform, from)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Error wrong from date, expected: yyyy-mm-dd")
		return
	}
	endDate, err := time.Parse(dform, until)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Error wrong until date, expected: yyyy-mm-dd")
		return
	}

	rows, err := db.Query(`SELECT amount, med_name, day, intake_interval_start, intake_interval_end, taken
                               FROM ScheduledDosages as SD JOIN 
                                 (SELECT Dosages.id, amount, intake_interval_start, intake_interval_end, med_name 
                                  FROM Dosages JOIN Medicines 
                                     ON Dosages.medicine_id = Medicines.id
                                  WHERE patient_id = ?) as DM
                               ON SD.dosage = DM.id
                               WHERE day BETWEEN ? AND ?`,
		patientID, startDate.Format(dform), endDate.Format(dform))
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	dosages := []ScheduledDosage{}
	for rows.Next() {
		var amount int
		var medicine, day, intakeIntervalStart, intakeIntervalEnd string
		var taken bool
		err = rows.Scan(&amount, &medicine, &day, &intakeIntervalStart, &intakeIntervalEnd, &taken)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, ScheduledDosage{
			Dosage{intakeIntervalStart, intakeIntervalEnd, amount, Medicine{medicine}},
			day,
			taken,
		})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}
	ar.setResponse(dosages)
}

// UPDATE
func updateScheduledDosage(r *http.Request, ar *APIResponse) {
	// Scan patient ID
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	// Read input scheduled dosages
	scheduledDosages := []ScheduledDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&scheduledDosages)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Unable to decode given json data")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Unable to start database transaction")
		return
	}

	for _, scheduledDosage := range scheduledDosages {
		dosageID, err := queryDosageID(scheduledDosage.Dosage, patientID)
		if err != nil {
			ar.setErrorAndStatus(http.StatusBadRequest, err, "")
			return
		}
		_, err = tx.Exec("UPDATE ScheduledDosages SET taken = ? WHERE dosage = ? AND day = ?",
			scheduledDosage.Taken, dosageID, scheduledDosage.Day)
		if err != nil {
			ar.setError(errorWithRollback(err, tx), "Something went wrong during SQL Update query")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		ar.setError(err, "Unable to commit changes to the database")
		return
	}
}

// DELETE
func deleteScheduledDosage(r *http.Request, ar *APIResponse) {
	// Scan patient ID
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Cannot convert patient id to an integer")
		return
	}

	// Read input scheduled dosages
	scheduledDosage := ScheduledDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&scheduledDosage)

	dosageID, err := queryDosageID(scheduledDosage.Dosage, patientID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Unable to start database transaction")
		return
	}

	_, err = tx.Exec(`DELETE FROM ScheduledDosages WHERE dosage = ? AND day = ?`, dosageID, scheduledDosage.Day)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Database failure")
		return
	}

	if err := tx.Commit(); err != nil {
		ar.setError(err, "Unable to commit changes to the database")
		return
	}
}
