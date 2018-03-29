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
	_ "github.com/go-sql-driver/mysql" // anonymous import
  //"go/token"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // anonymous import
)

var (
	db *sql.DB
)

func main() {
	var err error
	rootpasswd, dbname, listen_location := "pass", "database", "localhost:8080" // just some values
	fmt.Print("MySql Root Password: ")
	fmt.Scanf("%s", &rootpasswd)
	fmt.Print("\nMySql Database Name: ")
	fmt.Scanf("%s", &dbname)
	fmt.Print("\nRouter Listen Location: ")
	fmt.Scanf("%s", &listen_location)
	db, err = sql.Open("mysql", "root:" + rootpasswd + "@/" + dbname)

	if err != nil {
		log.Printf("encountered error while connecting to database: %v", err)
	}

	log.Printf("Connected to database '%s', and listening on '%s'...", dbname, listen_location)
	router := mux.NewRouter()
	router.Handle("/api/your extension", handlerWrapper(exampleHandler))
	router.Handle("/api/pushPatient", handlerWrapper(pushPatient))
	router.Handle("/api/deletePatient", handlerWrapper(deletePatient))
	router.Handle("/api/modifyPatient", handlerWrapper(modifyPatient))
	router.Handle("/api/pushPhysician", handlerWrapper(pushPhysician))
	router.Handle("/api/deletePhysician", handlerWrapper(deletePhysician))
	router.Handle("/api/modifyPhysician", handlerWrapper(modifyPhysician))
	router.Handle("/api/getDosages/", handlerWrapper(getDosages))
	router.Handle("/api/pushDosages", handlerWrapper(pushDosages))
	router.Handle("/api/getNotes/", handlerWrapper(getNotes))
	router.Handle("/api/addNote/", handlerWrapper(addNote))
	http.ListenAndServe(listen_location, router)
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

func pushPhysician(r *http.Request, responseChan chan []byte, errorChan chan error){
  physician := Physician{}
  dec := json.NewDecoder(r.Body)
  err := dec.Decode(&physician)
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to decode incoming JSON")
    return
  }
  physician.Password, err = HashPassword(physician.Password)
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to hash password")
    return
  }
  tx, err := db.Begin()
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to start transaction")
    return
  }
  role := "physician"
  result, err := tx.Exec(`INSERT INTO account (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, physician.Name, physician.Username, physician.Password, role)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  id, err := result.LastInsertId()
  _, err = tx.Exec(`INSERT INTO physician VALUES(?, ?, ?)`,
                   id, physician.Email, physician.CreationToken)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }

  errorChan <- tx.Commit()

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

func deletePhysician(r *http.Request, responseChan chan []byte, errorChan chan error){
  Id := r.URL.Query().Get("id")
  tx, err := db.Begin()
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to start transaction")
    return
  }
  _, err = tx.Exec(`DELETE FROM physician  WHERE id=?`, Id)
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
	err := dec.Decode(&patient)
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
                 pass_hash = ?
                 WHERE username = ?`, patient.Name, patient.Username, patient.Username)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()

}

func modifyPhysician(r *http.Request, responseChan chan []byte, errorChan chan error){
	physician := Physician{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&physician)
	if err  != nil {
		errorChan <- err
		return
	}
	physician.Password, err = HashPassword(physician.Password)
	if err != nil{
		errorChan <- errors.Wrap(err, "Hashing failed")
		return
	}
	tx, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "failed to start transaction")
		return
	}
	var id int
	err = db.QueryRow("SELECT id FROM account WHERE username=?", physician.Username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			errorChan <- errors.Wrap(err, "No matched rows")
		} else {
		errorChan <- errors.Wrap(err, "Database failure")
		return
  	}
	}
	_, err = tx.Exec(`UPDATE account SET
                          name = ?,
                          pass_hash = ?
                          WHERE id=?`, physician.Name, physician.Password, id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`UPDATE physician SET
                          email = ?,
                          token = ?
                          WHERE id = ?`, physician.Email, physician.CreationToken, id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()

}

// start/ end dates might be optional ?
//  Possible defaults:
//     start_date = [current_day]
//     end_date   = start_date + 1 month
func getDosages(r *http.Request, responseChan chan []byte, errorChan chan error) {
	// verify patient ?
	patient_id := r.URL.Query().Get("patient_id")

	from  := r.URL.Query().Get("from")
	until := r.URL.Query().Get("until")
	const dform = "2006-01-02" // specifies YYYY-MM-DD format
	start_date, err := time.Parse(dform, from)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error in parsing starting date")
		return
	}
	end_date, err := time.Parse(dform, until)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error in parsing end time")
		return
	}

	rows, err := db.Query(`SELECT amount, med_name, day, intake_time 
                               FROM dosage JOIN medicine 
                               ON dosage.medicine_id = medicine.id
                               WHERE patient_id = ? AND (day BETWEEN ? AND ?)
                               `,
		patient_id, start_date.Format(dform), end_date.Format(dform)) //add AND day BETWEEN ? AND ?
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during query")
		return
	}
	
	dosages := []Dosage{}
	for rows.Next() {
		var amount int
		var medicine, day, intake_time string
		err = rows.Scan(&amount, &medicine, &day, &intake_time)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		dosages = append(dosages, Dosage{day, intake_time, amount, medicine, false}) //false for now
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	json_values, err := json.Marshal(dosages)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when converting to JSON")
		return
	}
	responseChan <- json_values
	errorChan <- nil
	return
}

// Possible to also add a time interval?
// Or all 'untreated' notes
func getNotes(r *http.Request, responseChan chan []byte, errorChan chan error) {
	// verify patient
	patient_id := r.URL.Query().Get("patient_id")
	rows, err := db.Query(`SELECT question, day FROM note WHERE patient_Id = ?`, patient_id)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during query")
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note, date string
		err = rows.Scan(&note, &date)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		notes = append(notes, Note{note, date})
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	json_values, err := json.Marshal(notes)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when converting to JSON")
		return
	}
	responseChan <- json_values
	errorChan <- nil
	return
}

func addNote(r *http.Request, responseChan chan []byte, errorChan chan error) {
	// verify patient

	patient_id := r.URL.Query().Get("patient_id")
	note := Note{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&note)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}
	
	trans, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start new transaction")
		return
	}
	_, err = trans.Exec(
		`INSERT INTO note (patient_id, question, day) VALUES (?, ?, ?)`,
		patient_id, note.Note, note.CreatedAt)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to insert note into the database")
		return
	}

	errorChan <- trans.Commit()
	return
}

func pushDosages(r *http.Request, responseChan chan []byte, errorChan chan error) {
  patient_id := r.URL.Query().Get("patient_id")
  dosage := Dosage{}
  var medicine_id int
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
  err = tx.QueryRow(`SELECT id FROM medicine WHERE med_name = ?`, dosage.Medicine).Scan(&medicine_id)
  if err != nil{
    if err == sql.ErrNoRows{
      errorChan <- errors.Wrap(err, "Unknown medicine")
    } else{
      errorChan <- errors.Wrap(err, "Failed to execute query")
    }
    tx.Rollback()
    return
  }
  _, err = tx.Exec(`INSERT INTO dosage (amount, patient_id, medicine_id, day, intake_time) VALUES (?, ?, ?, ?, ?)`,
                                              dosage.NumberOfPills, patient_id, medicine_id, dosage.Day, dosage.IntakeMoment)
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