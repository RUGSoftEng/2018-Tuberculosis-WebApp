package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	http "net/http"
	"github.com/gorilla/mux"
	"strconv"
)

type fn func(*http.Request, chan APIResponse, chan error)

// HashPassword : placeholder function for hasing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash : compares a given unhashed password and hashed password
func CheckPasswordHash(password, hash string, errorChan chan error) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		errorChan <- errors.Wrap(err, "Authentication failed")
		return false
	}
	return true
}

//This function validates a password against a specific user, and issues a JWT Token
func login(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
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
	if !CheckPasswordHash(cred.Password, password, errorChan) {
		log.Println("Mismatching credentials")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": cred.Username,
		"password": cred.Password})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to generate JWT token")
		return
	}
	responseChan <- APIResponse{JWToken{Token: tokenString}, http.StatusOK}
	return
}

func parseToken(in JWToken, errorChan chan error, responseChan chan APIResponse, id int) bool {
	content := in.Token
	token, err := jwt.Parse(content, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("There was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		errorChan <- errors.Wrap(err, "Invalid token")
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user UserValidation
		mapstructure.Decode(claims, &user)
		var pwd string
		err := db.QueryRow(`SELECT pass_hash FROM Accounts WHERE username=?`, user.Username).Scan(&pwd)
		if err != nil {
			errorChan <- errors.Wrap(err, "Database failure")
			return false
		}
		var readId int
		err = db.QueryRow(`SELECT id FROM Accounts WHERE username=?`, user.Username).Scan(&readId)
		if err != nil {
			errorChan <- errors.Wrap(err, "Database failure")
			return false
		}
		if id != readId{
			errorChan <- errors.Wrap(errors.New("Wrong credentials"), "Wrong credentials")
			return false
		}
		if !CheckPasswordHash(user.Password, pwd, errorChan) {
			errorChan <- errors.New("Invalid credentials")
			return false
		} else {
			return true
		}
		return false
	}
	errorChan <- errors.New("Invalid token")
	return false
}

// Token authentication will probably be embedded in all the request that are give access
// to restricted contents, this functions is only for test purposes, but it uses the
// tokenParse() function that will do the core of the work


// I think this function shield potentially return an error if the credentials are invalid, instead of the boolean
func authenticate(r *http.Request, responseChan chan APIResponse, errorChan chan error) bool {
  vars := mux.Vars(r)
  id1 := vars["id"]
	id, _ := strconv.Atoi(id1)
	token := r.Header.Get("access_token")
	pass := JWToken{Token:token}
	log.Println(pass)
	if parseToken(pass, errorChan, responseChan, id) {
		return true
	} else {
		log.Println("Access denied")
		return false
	}
}

func authWrapper(handler func(r *http.Request, responseChan chan APIResponse, errorChan chan error)) func(*http.Request, chan APIResponse, chan error) {
	return func(req *http.Request, resChan chan APIResponse, errChan chan error) {
		log.Println("Are we here?")
		if !authenticate(req, resChan, errChan) {
			log.Println("Are we here in the if statement :p?")
			errChan <- errors.New("Invalid Credentials")
			return
		}
		handler(req, resChan, errChan)
	}
}
