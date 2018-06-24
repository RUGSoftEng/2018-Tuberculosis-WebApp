package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	http "net/http"
	"time"
)

// CREATE
func createScheduledDosages(r *http.Request, ar *APIResponse) {
	// Scan Patient ID
	patientID, err := getURLVariableInt(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	// Decode Request Body to JSON
	in := InputScheduledDosage{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&in)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	dosageID, err := queryDosageID(in.Dosage, patientID)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	// Start Database Transaction
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	// Add a ScheduledDosage for each given day
	for _, day := range in.Days {
		_, err := tx.Exec(`INSERT INTO ScheduledDosages VALUES (?, ?, ?)`, dosageID, day, false)
		if err != nil {
			ar.setErrorDBInsert(err, tx)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}
	ar.setStatus(StatusCreated)
}

// RETRIEVE
func retrieveScheduledDosages(r *http.Request, ar *APIResponse) {
	patientID, err := getURLVariable(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	from := r.URL.Query().Get("from")
	until := r.URL.Query().Get("until")
	const dform = "2006-01-02" // specifies YYYY-MM-DD format
	startDate, err := time.Parse(dform, from)
	if err != nil {
		ar.setErrorDate(err)
		return
	}
	endDate, err := time.Parse(dform, until)
	if err != nil {
		ar.setErrorDate(err)
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
		ar.setErrorDBSelect(err)
		return
	}

	dosages := []ScheduledDosage{}
	for rows.Next() {
		var amount int
		var medicine, day, intakeIntervalStart, intakeIntervalEnd string
		var taken bool
		err = rows.Scan(&amount, &medicine, &day, &intakeIntervalStart, &intakeIntervalEnd, &taken)
		if err != nil {
			ar.setErrorDBScan(err)
			return
		}
		dosages = append(dosages, ScheduledDosage{
			Dosage{intakeIntervalStart, intakeIntervalEnd, amount, Medicine{medicine}},
			day,
			taken,
		})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorDBAfter(err)
		return
	}

	ar.setResponse(dosages)
}

// UPDATE
func updateScheduledDosage(r *http.Request, ar *APIResponse) {
	// Scan patient ID
	patientID, err := getURLVariableInt(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	// Read input scheduled dosages
	scheduledDosages := []ScheduledDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&scheduledDosages)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	for _, scheduledDosage := range scheduledDosages {
		dosageID, err := queryDosageID(scheduledDosage.Dosage, patientID)
		if err != nil {
			ar.setErrorDBSelect(errorWithRollback(err, tx))
			return
		}
		_, err = tx.Exec("UPDATE ScheduledDosages SET taken = ? WHERE dosage = ? AND day = ?",
			scheduledDosage.Taken, dosageID, scheduledDosage.Day)
		if err != nil {
			ar.setErrorDBUpdate(err, tx)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusUpdated)
}

// DELETE
func deleteScheduledDosage(r *http.Request, ar *APIResponse) {
	// Scan patient ID
	patientID, err := getURLVariableInt(r, "id")
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	// Read input scheduled dosages
	scheduledDosage := ScheduledDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&scheduledDosage)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	dosageID, err := queryDosageID(scheduledDosage.Dosage, patientID)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`DELETE FROM ScheduledDosages WHERE dosage = ? AND day = ?`, dosageID, scheduledDosage.Day)
	if err != nil {
		ar.setErrorDBDelete(err, tx)
		return
	}

	if err := tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusDeleted)
}
