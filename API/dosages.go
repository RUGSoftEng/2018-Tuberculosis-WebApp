package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
	"time"
)

// CREATE
func pushDosage(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	vars := mux.Vars(r)
	patientID := vars["id"]
	dosage := Dosage{}
	var medicineID int
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&dosage)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to decode JSON")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	err = tx.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		dosage.Medicine.Name).Scan(&medicineID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorChan <- errors.Wrap(err, "Unknown medicine")
		} else {
			errorChan <- errors.Wrap(err, "Failed to execute query")
		}
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`INSERT INTO Dosages (patient_id, medicine_id, amount, intake_time) 
                          VALUES (?, ?, ?, ?)`,
		patientID, medicineID, dosage.NumberOfPills, dosage.IntakeMoment)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return		
	}
	responseChan <- APIResponse{nil, http.StatusCreated}
}

// RETRIEVE
// start/ end dates might be optional ?
//  Possible defaults:
//     startDate = [current_day]
//     endDate   = startDate + 1 month
func getDosages(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	// verify patient ?
	vars := mux.Vars(r)
	patientID := vars["id"]

	from := r.URL.Query().Get("from") // maybe check if specified
	until := r.URL.Query().Get("until")
	const dform = "2006-01-02" // specifies YYYY-MM-DD format
	startDate, err := time.Parse(dform, from)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error in parsing starting date")
		return
	}
	endDate, err := time.Parse(dform, until)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error in parsing end time")
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
		errorChan <- errors.Wrap(err, "Unexpected error during query")
		return
	}

	dosages := []ScheduledDosage{}
	for rows.Next() {
		var amount int
		var medicine, day, intakeTime string
		var taken bool
		err = rows.Scan(&amount, &medicine, &day, &intakeTime, &taken)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, ScheduledDosage{
			Dosage{intakeTime, amount, Medicine{medicine}},
			day,
			taken,
		})
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}
	responseChan <- APIResponse{dosages, http.StatusOK}
}
