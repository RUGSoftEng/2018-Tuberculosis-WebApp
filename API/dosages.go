package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	http "net/http"
	"time"
)

// CREATE
func pushDosage(r *http.Request, ar *APIResponse) {
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

	vars := mux.Vars(r)
	patientID := vars["id"]
	_, err = tx.Exec(`INSERT INTO Dosages (patient_id, medicine_id, amount, intake_time) 
                          VALUES (?, ?, ?, ?)`,
		patientID, medicineID, dosage.NumberOfPills, dosage.IntakeMoment)
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

// InputDosagesJSON : Temp struct
type InputDosagesJSON struct {
	Medicine Medicine `json:"medicine"`
	Days     []string `json:"days"`
}

// CREATE
func addScheduledDosages(r *http.Request, ar *APIResponse) {
	// Scan Patient ID
	vars := mux.Vars(r)
	patientID := vars["id"]

	// Decode Request Body to JSON
	in := InputDosagesJSON{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&in)
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
// start/ end dates might be optional ?
//  Possible defaults:
//     startDate = [current_day]
//     endDate   = startDate + 1 month
func getDosages(r *http.Request, ar *APIResponse) {

	vars := mux.Vars(r)
	patientID := vars["id"]

	from := r.URL.Query().Get("from") // maybe check if specified
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

	rows, err := db.Query(`SELECT amount, med_name, day, intake_time, taken
                               FROM ScheduledDosages as SD JOIN 
                                 (SELECT Dosages.id, amount, intake_time, med_name 
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
		var medicine, day, intakeTime string
		var taken bool
		err = rows.Scan(&amount, &medicine, &day, &intakeTime, &taken)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, ScheduledDosage{
			Dosage{intakeTime, amount, Medicine{medicine}},
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
