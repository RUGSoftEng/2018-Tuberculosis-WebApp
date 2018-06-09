package main

import (
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
		ar.setErrorJSON(err)
		return
	}

	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	var medicineID int
	err = db.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		dosage.Medicine.Name).Scan(&medicineID)
	if err != nil {
		ar.setErrorDBScan(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`INSERT INTO Dosages (patient_id, medicine_id, amount,
                          intake_interval_start, intake_interval_end) 
                          VALUES (?, ?, ?, ?, ?)`,
		patientID, medicineID, dosage.NumberOfPills,
		dosage.IntakeIntervalStart, dosage.IntakeIntervalEnd)
	if err != nil {
		ar.setErrorDBInsert(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}
	ar.setStatus(StatusCreated)
}

// RETRIEVE
func retrieveDosages(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	dosages := []Dosage{}
	rows, err := db.Query(`SELECT amount, med_name, intake_interval_start, intake_interval_end
                          FROM Dosages JOIN Medicines 
                             ON Dosages.medicine_id = Medicines.id 
                          WHERE patient_id = ?`, patientID)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	for rows.Next() {
		var amount int
		var medicine, intakeIntervalStart, intakeIntervalEnd string
		err = rows.Scan(&amount, &medicine, &intakeIntervalStart, &intakeIntervalEnd)
		if err != nil {
			ar.setErrorDBScan(err)
			return
		}
		dosages = append(dosages, Dosage{
			intakeIntervalStart,
			intakeIntervalEnd,
			amount,
			Medicine{medicine}})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorDBAfter(err)
		return
	}
	ar.setResponse(dosages)
}

// UPDATE
func updateDosage(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	dosage := UpdateDosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&dosage)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	oldMedicineID, err := queryMedicineID(dosage.OldDosage.Medicine)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}
	newMedicineID, err := queryMedicineID(dosage.NewDosage.Medicine)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
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
		ar.setErrorDBUpdate(err, tx)
		return
	}

	if err := tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}
	ar.setStatus(StatusUpdated)
}

// DELETE
func deleteDosage(r *http.Request, ar *APIResponse) {
	patientID, err := getPatientIDVariable(r)
	if err != nil {
		ar.setErrorVariable(err)
		return
	}

	dosage := Dosage{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&dosage)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	medicineID, err := queryMedicineID(dosage.Medicine)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	_, err = tx.Exec(`DELETE FROM Dosages WHERE patient_id = ? AND medicine_id = ?`,
		patientID, medicineID)
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
