package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"log"
	http "net/http"
)

// CREATE
func pushPhysician(r *http.Request, ar *APIResponse) {
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Failed to decode incoming JSON.")
		return
	}
	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		ar.setError(err, "Failed to hash password")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction")
		return
	}
	role := "physician"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, physician.Name, physician.Username, physician.Password, role)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}
	id, err := result.LastInsertId()
	_, err = tx.Exec(`INSERT INTO Physicians VALUES(?, ?, ?)`,
		id, physician.Email, physician.CreationToken)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
		return
	}

	ar.StatusCode = http.StatusCreated
}

// UPDATE
func modifyPhysician(r *http.Request, ar *APIResponse) {
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Failed to decode incoming JSON")
		return
	}

	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		ar.setError(err, "Hashing failed.")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction.")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	_, err = tx.Exec(`UPDATE Accounts SET
                          name = ?,
                          pass_hash = ?
                          WHERE id=?`, physician.Name, physician.Password, id)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}
	_, err = tx.Exec(`UPDATE Physicians SET
                          email = ?,
                          token = ?
                          WHERE id = ?`, physician.Email, physician.CreationToken, id)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
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
		ar.setError(err, "Failed to start transaction.")
		return
	}
	_, err = tx.Exec(`DELETE FROM Physicians  WHERE id=?`, id)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}
	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, id)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
		return
	}
}
