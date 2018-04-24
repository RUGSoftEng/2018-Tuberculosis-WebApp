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
func pushDosage(r *http.Request, ar *APIResponse) {

	dosage := Dosage{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&dosage)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Failed to decode JSON")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "failed to start transaction")
		return
	}

	var medicineID int
	err = tx.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		dosage.Medicine.Name).Scan(&medicineID)
	if err != nil {
		if err == sql.ErrNoRows {
			ar.StatusCode = http.StatusNotFound
			ar.Error = errors.Wrap(err, "Unknown medicine")
		} else {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Failed to execute query")
		}
		err = tx.Rollback()
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Rollback Failed")
		}
		return
	}

	vars := mux.Vars(r)
	patientID := vars["id"]
	_, err = tx.Exec(`INSERT INTO Dosages (patient_id, medicine_id, amount, intake_time) 
                          VALUES (?, ?, ?, ?)`,
		patientID, medicineID, dosage.NumberOfPills, dosage.IntakeMoment)
	if err != nil {
		ar.Error = err
		err = tx.Rollback()
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Rollback Failed")
		}
		return
	}

	if err = tx.Commit(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
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
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Error wrong from date, expected: yyyy-mm-dd")
		return
	}
	endDate, err := time.Parse(dform, until)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Error wrong until date, expected: yyyy-mm-dd")
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
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error during query")
		return
	}

	dosages := []ScheduledDosage{}
	for rows.Next() {
		var amount int
		var medicine, day, intakeTime string
		var taken bool
		err = rows.Scan(&amount, &medicine, &day, &intakeTime, &taken)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, ScheduledDosage{
			Dosage{intakeTime, amount, Medicine{medicine}},
			day,
			taken,
		})
	}
	if err = rows.Err(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}
	ar.Data = dosages
}
