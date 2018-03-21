package main

import (
  "time"
  "log"
  http "net/http"
  "database/sql"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
  "golang.org/x/crypto/bcrypt"
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
  router.Handle("/api/createAccount", handlerWrapper(createAccountPatient))
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
      log.Printf("Server error: %v", err);
      http.Error(w, err.Error(), http.StatusInternalServerError)
    case <- time.After(1 * time.Millisecond):
      log.Printf("Response timeout")
    }
    return
  })
}

func exampleHandler(r *http.Request, responseChan chan []byte, errorChan chan error) {
  username := r.Header.Get("username")
  token := r.Header.Get("token")

  // expend this query to include the patient, or physician according to your api call
  rows, err := db.Query("SELECT * FROM account WHERE username=? AND password=?", username, token)
  if err != nil {
    if err == sql.ErrNoRows {
      errorChan <- errors.Wrap(err, "no valid login credentials")
      return
    }
    errorChan <- errors.Wrap(err, "encountered error during query")
    return
  }
  rows.Close()

  // do your own querries,
  // if you encounter a "err != nil" send it to the errorChan in the above matter
  // if all goed well, marshal your results and sen them to responseChan
  // responseChan <- "your marshalled data"

  return //obsolete
}

func createAccountPatient(r *http.Request, responseChan chan []byte, errorChan chan error){
  token := r.Header.Get("token")
  rows, err := db.Query("SELECT Id FROM Physician WHERE Token = ?", token)
  if err != nil{
    errorChan <- errors.Wrap(err, "Encountered database problem")
  }
  if !rows.Next(){
    log.Printf("Token not found")
    return
  }
  row := db.QueryRow("SELECT MAX(Id) FROM Account")
  var id int
  row.Scan(&id)
  id += 1
  name := r.Header.Get("name")
  username := r.Header.Get("username")
  password, err := HashPassword(r.Header.Get("password"))
  if err != nil{
    errorChan <- errors.Wrap(err, "Hashing failed")
  }
  role := r.Header.Get("role")
  db.Exec("INSERT INTO Account VALUES(?, ?, ?, ?, ?)", id, name, username, password, role)
  db.Exec("INSERT INTO Patient VALUES(?)", id)
}

func HashPassword(password string) (string, error) {
  bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
  return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  return err == nil
}
