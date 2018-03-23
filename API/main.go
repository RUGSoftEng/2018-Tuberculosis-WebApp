package main

import (
	"time"
	"log"
	http "net/http"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"encoding/json"
)

var (
	db *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("mysql", "database info")
	if err != nil {
		log.Printf("encountered error while connecting to database: %v", err)
	}

	router := mux.NewRouter()
	router.Handle("/api/your extension", handlerWrapper(exampleHandler))
	router.Handle("/api/pushPatient", handlerWrapper(pushPatient))
	router.Handle("/api/deletePatient", handlerWrapper(deletePatient))
	router.Handle("/api/modifyPatient", handlerWrapper(modifyPatient))
	http.ListenAndServe("portNumber", router)
}

func handlerWrapper(handler func(r *http.Request, responseChan chan []byte, errorChan chan error)) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		responseChan := make(chan []byte)
		errorChan := make(chan error)

		go handler(r, responseChan, errorChan)

		time.After(2 * time.Second)

		select {
		case body := <- responseChan:
			w.Header().Set("Content-Type", "application/json")
			w.Write(body)
		case err := <- errorChan:
			if err != nil {
				log.Printf("Server error: %v", err);
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
		case <- time.After(1 * time.Millisecond):
			log.Printf("Response timeout")
		}
		return
	})
}

func exampleHandler(r *http.Request, responseChan chan []byte, errorChan chan error) {
	ID := 0
	apiToken := r.Header.Get("api_token")

	// This is a join example for a patient call, change to physician it is a call only a physician can make
	// remove join part if it is a call able for both
	err := db.QueryRow(`	SELECT id 
								FROM patient AS pa 
									INNER JOIN account AS acc 
									ON pa.id = acc.id  
								WHERE acc.api_token = ?`,
		apiToken).Scan(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			errorChan <- errors.Wrap(err, "no valid login credentials")
			return
		}
		errorChan <- errors.Wrap(err, "encountered error during query")
		return
	}

	// if you are going to insert multiple things in the database do this using a transaction.
	// see insertPatient

	// do your own querries,
	// if you encounter a "err != nil" send it to the errorChan in the above matter
	// if all goed well, marshal your results and sen them to responseChan

	// End for a get function
	// responseChan <- "your marshalled data"

	// End for a succesfull push or put function
	// errorChan <- nil
}

// expects a json file containing the new patient and a url encoded physician token
func pushPatient(r *http.Request, responseChan chan []byte, errorChan chan error){
	physicianId := 0
	physicianToken := r.URL.Query().Get("physician_token")

	// In general this will check the api_token
	err := db.QueryRow(`	SELECT id 
								FROM physician  
								WHERE token = ?`,
								physicianToken).Scan(physicianId)
	if err != nil{
		if err == sql.ErrNoRows {
			log.Printf("Physician not found")
			return
		}
		errorChan <- errors.Wrap(err, "Encountered database problem")
		return
	}

	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(patient)
	if err  != nil {
		errorChan <- err
		return
	}

	patient.Password, err = HashPassword(patient.Password)
	if err != nil{
		errorChan <- errors.Wrap(err, "Hashing failed")
		return
	}

	// If you are going to do multiple (on each other depending) execs it is better to first start a transaction.
	// This insure that if one of them fail the others won't be done
	// Ill add a example here
	// Also don't forget to check errors.
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	result, err := tx.Exec(`INSERT  INTO account (name, username, pass_hash) VALUES(?, ?, ?)`, patient.Name, patient.Username, patient.Password) // name is reserved keyword
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}
	id, err := result.LastInsertId() //this gets the id that would be created for above insert
	_, err = tx.Exec(`INSERT INTO patient (id, physician_id) VALUES(?, ?)`,  id, physicianId) //physician is not necessery here, however, it is a be easier to read
	if err != nil {
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()//actually commits the changes to the database
}

func deletePatient(r *http.Request, responseChan chan []byte, errorChan chan error){
  Id := r.URL.Query().Get("id")
  tx, err := db.Begin()
  if err != nil {
	  errorChan <- errors.Wrap(err, "failed to start transaction")
	  return
	}
	_ , err = tx.Exec(`DELETE FROM patient WHERE id=?`,Id )
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`DELETE FROM account WHERE id=?`, Id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()
}

func modifyPatient(r *http.Request, responseChan chan []byte, errorChan chan error){
	patient := Patient{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(patient)
	if err  != nil {
		errorChan <- err
		return
	}
	patient.Password, err = HashPassword(patient.Password)
	if err != nil{
		errorChan <- errors.Wrap(err, "Hashing failed")
		return
	}

	// Using a transaction because I don't know whether we are going to have to add
	// query for a possible change of physician (or how to do that)
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	tx.Exec(`UPDATE account SET 
                 name = ?,
                 username = ?,
                 pass_hash = ?`, patient.Name, patient.Username, patient.Password )
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()

}

// placeHolderFunction
func HashPassword(password string) (string, error) {
	return password, nil
}

// See slack message
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// it's better to return the error here. otherwise you know there was a error, but you don't have the error message
	return err == nil
}
