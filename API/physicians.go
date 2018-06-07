package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"log"
	http "net/http"
)

// CREATE
func createPhysician(r *http.Request, ar *APIResponse) {
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Failed to decode incoming JSON.")
		return
	}
	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to hash password")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction")
		return
	}
	role := "physician"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, physician.Name, physician.Username, physician.Password, role)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}
	id, err := result.LastInsertId()
	_, err = tx.Exec(`INSERT INTO Physicians VALUES(?, ?, ?)`,
		id, physician.Email, physician.CreationToken)
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
func retrievePatients(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	physicianID := vars["id"]
	rows, err := db.Query(`SELECT Accounts.id, Accounts.name 
                              FROM Accounts INNER JOIN Patients 
                              ON Accounts.id=Patients.id AND Patients.physician_id=?`, physicianID)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	patients := []PatientInfo{}
	for rows.Next() {
		var name string
		var id int
		err = rows.Scan(&id, &name)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		patients = append(patients, PatientInfo{id, name})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}

	ar.setResponse(patients)

}

// RETRIEVE
func retrievePyByID(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	doc := PhysicianOverview{}
	var name string
	var user string
	err := db.QueryRow(`SELECT name, username FROM Accounts WHERE id=?`, id).Scan(&name, &user)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	var email string
	var token string
	err = db.QueryRow(`SELECT email,token FROM Physicians WHERE id=?`, id).Scan(&email, &token)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	doc.Name = name
	doc.Username = user
	doc.Email = email
	doc.Token = token
	ar.setResponse(doc)
}

// RETRIEVE
func retrievePyByUsername(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	username := vars["username"]
	doc := PhysicianOverview{}
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
	var email string
	var token string
	err = db.QueryRow(`SELECT email,token FROM Physicians WHERE id=?`, id).Scan(&email, &token)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
	doc.Name = name
	doc.Username = username
	doc.Email = email
	doc.Token = token
	ar.setResponse(doc)
}

// UPDATE
func updatePhysician(r *http.Request, ar *APIResponse) {
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode incoming JSON")
		return
	}

	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Hashing failed.")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	_, err = tx.Exec(`UPDATE Accounts SET
                          name = ?,
                          pass_hash = ?
                          WHERE id=?`, physician.Name, physician.Password, id)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}
	_, err = tx.Exec(`UPDATE Physicians SET
                          email = ?,
                          token = ?
                          WHERE id = ?`, physician.Email, physician.CreationToken, id)
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
func deletePhysician(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println(id)
	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}
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
