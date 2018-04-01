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
  "github.com/dgrijalva/jwt-go"
	"fmt"
)

var (
	db *sql.DB
)

func main() {
	var err error
	rootpasswd, dbname, listen_location := "pass", "database", "localhost:8080" // just some values
	fmt.Scanf("%s", &rootpasswd)
	fmt.Scanf("%s", &dbname)
	fmt.Scanf("%s", &listen_location)
	db, err = sql.Open("mysql", "root:" + rootpasswd + "@/" + dbname)

	if err != nil {
		log.Printf("encountered error while connecting to database: %v", err)
	}

	log.Printf("Connected to database '%s', and listening on '%s'...", dbname, listen_location)
	router := mux.NewRouter()
	router.Handle("/api/your extension", handlerWrapper(exampleHandler))

	// GET Requests for Retrieving
	get_router := router.Methods("GET").Subrouter()
	get_router.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(getDosages))
	get_router.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(getNotes))
	get_router.Handle("/api/general/videos/{topic}", handlerWrapper(getVideoByTopic))
	
	// POST Requests for Updating
	post_router := router.Methods("POST").Subrouter()
	post_router.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(modifyPatient))
	post_router.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(modifyPhysician))
  post_router.Handle("/api/accounts/login", handlerWrapper(login))

	// PUT Requests for Creating
	put_router := router.Methods("PUT").Subrouter()
	put_router.Handle("/api/accounts/patients", handlerWrapper(pushPatient))
	put_router.Handle("/api/accounts/physicians", handlerWrapper(pushPhysician))
	put_router.Handle("/api/accounts/patients/{id:[0-9]+}/dosages", handlerWrapper(pushDosages))
	put_router.Handle("/api/accounts/patients/{id:[0-9]+}/notes", handlerWrapper(addNote))
	put_router.Handle("/api/general/videos", handlerWrapper(addVideo))
	
	// DELETE Requests for Deleting
	delete_router := router.Methods("DELETE").Subrouter()
	delete_router.Handle("/api/accounts/patients/{id:[0-9]+}", handlerWrapper(deletePatient))	
	delete_router.Handle("/api/accounts/physicians/{id:[0-9]+}", handlerWrapper(deletePhysician))

	// Starting the router
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
			log.Println("We here?")
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
		case <- time.After(5 * time.Second):
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
								FROM Patients AS pa 
									INNER JOIN Accounts AS acc 
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
  patient := Patient{}
  dec := json.NewDecoder(r.Body)
  err := dec.Decode(&patient)
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to decode incoming JSON")
    return
  }
  patient.Password, err = HashPassword(patient.Password)
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to hash password")
    return
  }
  tx, err := db.Begin()
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to start transaction")
    return
  }
  role := "patient"
  result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, patient.Name, patient.Username, patient.Password, role)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  id, err := result.LastInsertId()
  creation_token := r.URL.Query().Get("token")
  log.Println(creation_token)
  var physician_id int
  err = tx.QueryRow(`SELECT id FROM Physicians WHERE token=?`, creation_token).Scan(&physician_id)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  log.Println(physician_id)
  _, err = tx.Exec(`INSERT INTO Patients VALUES(?,?)`, id, physician_id)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  errorChan <- tx.Commit()
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
  result, err := tx.Exec(`INSERT INTO Accounts (name, username, pass_hash, role)
                                VALUES(?, ?, ?, ?)`, physician.Name, physician.Username, physician.Password, role)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  id, err := result.LastInsertId()
  _, err = tx.Exec(`INSERT INTO Physicians VALUES(?, ?, ?)`,
                   id, physician.Email, physician.CreationToken)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }

  errorChan <- tx.Commit()

}

func deletePatient(r *http.Request, responseChan chan []byte, errorChan chan error){
  vars := mux.Vars(r)
  Id := vars["id"]
  tx, err := db.Begin()
  if err != nil {
	  errorChan <- errors.Wrap(err, "failed to start transaction")
	  return
	}
  _ , err = tx.Exec(`DELETE FROM Notes WHERE patient_Id=?`,Id )
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  _ , err = tx.Exec(`DELETE FROM Dosages WHERE patient_id=?`,Id )
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  _ , err = tx.Exec(`DELETE FROM Patients WHERE id=?`,Id )
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, Id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()
}

func deletePhysician(r *http.Request, responseChan chan []byte, errorChan chan error){
  vars := mux.Vars(r)
  Id := vars["id"]
  log.Println(Id)
  tx, err := db.Begin()
  if err != nil{
    errorChan <- errors.Wrap(err, "Failed to start transaction")
    return
  }
  _, err = tx.Exec(`DELETE FROM Physicians  WHERE id=?`, Id)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }
  _, err = tx.Exec(`DELETE FROM Accounts WHERE id=?`, Id)
  if err != nil{
    errorChan <- err
    tx.Rollback()
    return
  }

  errorChan <- tx.Commit()
}

func modifyPatient(r *http.Request, responseChan chan []byte, errorChan chan error){
  vars := mux.Vars(r)
  id := vars["id"]
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
	tx.Exec(`UPDATE Accounts SET 
                 name = ?,
                 pass_hash = ?
                 WHERE id = ?`, patient.Name, patient.Password, id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}

	errorChan <- tx.Commit()

}

func modifyPhysician(r *http.Request, responseChan chan []byte, errorChan chan error){
  vars := mux.Vars(r)
  id := vars["id"]
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
	_, err = tx.Exec(`UPDATE Accounts SET
                          name = ?,
                          pass_hash = ?
                          WHERE id=?`, physician.Name, physician.Password, id)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`UPDATE Physicians SET
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
	vars := mux.Vars(r)
	patient_id := vars["id"]

	from  := r.URL.Query().Get("from") // maybe check if specified
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
                               FROM Dosages JOIN Medicines 
                               ON Dosages.medicine_id = Medicines.id
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
	vars := mux.Vars(r)
	patient_id := vars["id"]
	
	rows, err := db.Query(`SELECT question, day FROM Notes WHERE patient_Id = ?`, patient_id)
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
	vars := mux.Vars(r)
	patient_id := vars["id"]

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
		`INSERT INTO Notes (patient_id, question, day) VALUES (?, ?, ?)`,
		patient_id, note.Note, note.CreatedAt)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to insert note into the database")
		return
	}

	errorChan <- trans.Commit()
	return
}

func pushDosages(r *http.Request, responseChan chan []byte, errorChan chan error) {
	vars := mux.Vars(r)
	patient_id := vars["id"]
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
	err = tx.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`, dosage.Medicine).Scan(&medicine_id)
	if err != nil{
		if err == sql.ErrNoRows{
			errorChan <- errors.Wrap(err, "Unknown medicine")
		} else{
			errorChan <- errors.Wrap(err, "Failed to execute query")
		}
		tx.Rollback()
		return
	}
	_, err = tx.Exec(`INSERT INTO Dosages (amount, patient_id, medicine_id, day, intake_time) VALUES (?, ?, ?, ?, ?)`,
		dosage.NumberOfPills, patient_id, medicine_id, dosage.Day, dosage.IntakeMoment)
	if err != nil{
		errorChan <- err
		tx.Rollback()
		return
	}
	
	errorChan <- tx.Commit()
}


func getVideoByTopic(r *http.Request, responseChan chan []byte, errorChan chan error) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	rows, err := db.Query(`SELECT topic, title, reference FROM Videos WHERE topic = ?`, topic)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when querying the database")
		return
	}

	videos := []Video{}
	for rows.Next() {
		var topic, title, reference string
		err = rows.Scan(&topic, &title, &reference)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		videos = append(videos, Video{topic, title, reference})
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	json_values, err := json.Marshal(videos)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when converting to JSON")
		return
	}
	responseChan <- json_values
	errorChan <- nil
	return
}

func addVideo(r *http.Request, responseChan chan []byte, errorChan chan error) {
	video := Video{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&video)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}
	
	trans, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start new transaction")
		return
	}
	_, err = trans.Exec(`INSERT INTO Videos (topic, title, reference) VALUES (?, ?, ?)`,
		video.Topic, video.Title, video.Reference)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to insert video into the database")
		return
	}

	errorChan <- trans.Commit()
	return
}

// placeHolderFunction
func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

// See slack message
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// it's better to return the error here. otherwise you know there was a error, but you don't have the error message
	return err == nil
}

//This function validates a password against a specific user, and issues a JWT Token

func login(r *http.Request, responseChan chan []byte, errorChan chan error){
  cred := UserValidation{}
  err := json.NewDecoder(r.Body).Decode(&cred)
  if err != nil {
  errorChan <- errors.Wrap(err, "Failed to decode user credentials")
    return
  }
  var password string
  err = db.QueryRow(`SELECT pass_hash FROM Accounts WHERE username=?`, cred.Username).Scan(&password)
  if err != nil {
  	errorChan <- errors.Wrap(err, "Database failure")
    return
  }
  if !CheckPasswordHash(cred.Password, password){
    log.Println("Mismatching credentials")
    return
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
                 "username": cred.Username,
                 "password": cred.Password,})
  tokenString, err := token.SignedString([]byte("secret"))
  if err != nil {
   errorChan <- errors.Wrap(err, "Failed to generate JWT token")
    return
  }
  jsonToSend,err := json.Marshal(JWToken{Token:tokenString})
  if err != nil {
  errorChan <- errors.Wrap(err, "Failed to encode token")
    return
  }
  responseChan <- jsonToSend
  errorChan <- nil
  return
}
