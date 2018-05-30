package main

import (
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
// expects a json file containing the new patient and a url encoded physician token
func createPatient(r *http.Request, ar *APIResponse) {
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode incoming JSON.")
		return
	}

	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to hash password")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction")
		return
	}

	role := "patient"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, patient.Name, patient.Username, patient.Password, role)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok {
			ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unknown error")
			return
		}
		if me.Number == 1062 {
			ar.setErrorAndStatus(StatusDatabaseConstraintViolation, errors.New("Username already in use"), "please choose another one")
			return
		}
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	var physicianID int
	creationToken := r.URL.Query().Get("token")
	err = tx.QueryRow(`SELECT id FROM Physicians WHERE token=?`, creationToken).Scan(&physicianID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	_, err = tx.Exec(`INSERT INTO Patients VALUES(?,?)`, id, physicianID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	err = tx.Commit()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusCreated
}

// UPDATE
func updatePatient(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&patient)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode incoming JSON")
		return
	}
	patient.Password, err = HashPassword(patient.Password)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Hashing failed.")
		return
	}

	// Using a transaction because I don't know whether we are going to have to add
	// query for a possible change of physician (or how to do that)
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	_, err = tx.Exec(`UPDATE Accounts SET 
                 name = ?,
                 pass_hash = ?
                 WHERE id = ?`, patient.Name, patient.Password, id)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}
}

// DELETE
func deletePatient(r *http.Request, ar *APIResponse) {
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, id)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}
}

//GET

func retrieveByID(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	patientOV := PatientOverview{}
	var name string
	var user string
	err := db.QueryRow(`SELECT name, username FROM Accounts WHERE id=?`, id).Scan(&name, &user)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var phyID int
	err = db.QueryRow(`SELECT physician_id FROM Patients WHERE id=?`, id).Scan(&phyID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var physicianName string
	var email string
	err = db.QueryRow(`SELECT email FROM Physicians WHERE id=?`, phyID).Scan(&email)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	err = db.QueryRow(`SELECT name FROM Accounts WHERE id=?`, phyID).Scan(&physicianName)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	patientOV.Name = name
	patientOV.Username = user
	patientOV.PhysicianName = physicianName
	patientOV.PhysicianEmail = email
	ar.setResponse(patientOV)
}

func retrieveByUsername(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	username := vars["username"]
	patient := PatientOverview{}
	var name string
	err := db.QueryRow(`SELECT name FROM Accounts WHERE username=?`, username).Scan(&name)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var id int
	err = db.QueryRow(`SELECT id FROM Accounts WHERE username=?`, username).Scan(&id)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var phyID int
	err = db.QueryRow(`SELECT physician_id FROM Patients WHERE id=?`, id).Scan(&phyID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var physicianName string
	var email string
	err = db.QueryRow(`SELECT email FROM Physicians WHERE id=?`, phyID).Scan(&email)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	err = db.QueryRow(`SELECT name FROM Accounts WHERE id=?`, phyID).Scan(&physicianName)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	patient.Name = name
	patient.Username = username
	patient.PhysicianEmail = email
	patient.PhysicianName = physicianName
	ar.setResponse(patient)
}
