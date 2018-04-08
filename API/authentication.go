package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	http "net/http"
)

// HashPassword : placeholder function for hasing
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash : compares a given unhashed password and hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// it's better to return the error here. otherwise you know there was a error, but you don't have the error message
	return err == nil
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
	if !CheckPasswordHash(cred.Password, password) {
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

func parseToken(in JWToken, errorChan chan error, responseChan chan APIResponse) {
	content := in.Token
	token, _ := jwt.Parse(content, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user UserValidation
		mapstructure.Decode(claims, &user)
		var pwd string
		err := db.QueryRow(`SELECT pass_hash FROM Accounts WHERE username=?`, user.Username).Scan(&pwd)
		if err != nil {
			errorChan <- errors.Wrap(err, "Database failure")
			return
		}
		if !CheckPasswordHash(user.Password, pwd) {
			errorChan <- errors.New("Invalid token")
		} else {
			responseChan <- APIResponse{"You're authenticated", http.StatusOK}
		}
		return
	}
	errorChan <- errors.New("Invalid token")
}

// Token authentication will probably be embedded in all the request that are give access
// to restricted contents, this functions is only for test purposes, but it uses the
// tokenParse() function that will do the core of the work
func authenticate(r *http.Request, responseChan chan APIResponse, errorChan chan error) {
	pass := JWToken{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&pass)
	if err != nil {
		errorChan <- errors.Wrap(err, "Error while decoding")
	}
	parseToken(pass, errorChan, responseChan)
	return
}
