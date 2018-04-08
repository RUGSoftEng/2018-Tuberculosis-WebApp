package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"log"
	http "net/http"
)

// CREATE
func pushPhysician(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to decode incoming JSON")
		return
	}
	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to hash password")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start transaction")
		return
	}
	role := "physician"
	result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, physician.Name, physician.Username, physician.Password, role)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}
	id, err := result.LastInsertId()
	_, err = tx.Exec(`INSERT INTO Physicians VALUES(?, ?, ?)`,
		id, physician.Email, physician.CreationToken)
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

// UPDATE
func modifyPhysician(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	vars := mux.Vars(r)
	id := vars["id"]
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err != nil {
		errorChan <- err
		return
	}
	physician.Password, err = HashPassword(physician.Password)
	if err != nil {
		errorChan <- errors.Wrap(err, "Hashing failed")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	_, err = tx.Exec(`UPDATE Accounts SET
                          name = ?,
                          pass_hash = ?
                          WHERE id=?`, physician.Name, physician.Password, id)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`UPDATE Physicians SET
                          email = ?,
                          token = ?
                          WHERE id = ?`, physician.Email, physician.CreationToken, id)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return		
	}
	responseChan <- APIResponse{nil, http.StatusOK}
}

// DELETE
func deletePhysician(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println(id)
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start transaction")
		return
	}
	_, err = tx.Exec(`DELETE FROM Physicians  WHERE id=?`, id)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, id)
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}

	if err = tx.Commit(); err != nil {
		errorChan <- errors.Wrap(err, "Failed to commit changes to database.")
		return		
	}
	responseChan <- APIResponse{nil, http.StatusOK}
}
