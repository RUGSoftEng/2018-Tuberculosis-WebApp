package main

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	http "net/http"
	"strconv"
	"strings"
)

//HashPassword hashes the given string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//This function validates a password against a  specific user, and issues a JWT Token
func login(r *http.Request, ar *APIResponse) {
	cred := UserValidation{}
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to decode user credentials")
		return
	}
	var password string
	err = db.QueryRow(`SELECT pass_hash FROM Accounts WHERE username=?`, cred.Username).Scan(&password)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(cred.Password))
	if err != nil {
		ar.setErrorAndStatus(http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": cred.Username,
		"password": cred.Password})
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to generate JWT token")
		return
	}
	var tokenID int
	err = db.QueryRow(`SELECT id FROM Accounts WHERE username=?`, cred.Username).Scan(&tokenID)
	var salt int
	err = db.QueryRow(`SELECT api_token FROM Accounts WHERE username=?`, cred.Username).Scan(&salt)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	tokenString = Encode(tokenString, salt)
	ar.Data = JWToken{Token: tokenString, ID: tokenID}
}

func parseToken(in JWToken, ar *APIResponse) {
	content := in.Token
	id := in.ID
	var salt int
	err := db.QueryRow(`SELECT api_token FROM Accounts WHERE id=?`, id).Scan(&salt)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	content = Decode(content, salt)
	token, err := jwt.Parse(content, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("There was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		ar.setErrorAndStatus(http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	ar.setErrorAndStatus(http.StatusUnauthorized, err, "Unauthorized")

	var user UserValidation
	err = mapstructure.Decode(claims, &user)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	var pwd string
	err = db.QueryRow(`SELECT pass_hash FROM Accounts WHERE username=?`, user.Username).Scan(&pwd)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	var readID int
	err = db.QueryRow(`SELECT id FROM Accounts WHERE username=?`, user.Username).Scan(&readID)
	if err == sql.ErrNoRows {
		ar.setErrorAndStatus(http.StatusUnauthorized, err, "Unauthorized")
		return
	}
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	var physicianID = -1
	var physicianPwd = "default"
	var physicianUser = "default"
	err = db.QueryRow(`SELECT physician_id FROM Patients WHERE id=?`, id).Scan(&physicianID)
	if err != nil && err != sql.ErrNoRows {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	err = db.QueryRow(`SELECT username FROM Accounts WHERE id=?`, physicianID).Scan(&physicianUser)
	if err != nil && err != sql.ErrNoRows {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	err = db.QueryRow(`SELECT pass_hash FROM Accounts WHERE id=?`, physicianID).Scan(&physicianPwd)
	if err != nil && err != sql.ErrNoRows {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Database failure")
		return
	}
	// DO NOT REMOVE THE FOLLOWING LINES
	log.Println(id, readID, physicianID)
	if id != readID && readID != physicianID {
		ar.setErrorAndStatus(http.StatusUnauthorized, errors.New("Wrong user"), "Unauthorized")
		return
	}
	// UP UNTIL HERE
}

// Token authentication will probably be embedded in all the request that are give access
// to restricted contents, this functions is only for test purposes, but it uses the
// tokenParse() function that will do the core of the work

// I think this function shield potentially return an error if the credentials are invalid, instead of the boolean
func authenticate(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	id1 := vars["id"]
	id, err := strconv.Atoi(id1)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Conversion failed")
		return
	}
	token := r.Header.Get("access_token")
	pass := JWToken{Token: token, ID: id}
	parseToken(pass, ar)
}

// Rotate Latin letters by the shift amount.
func rotate(text string, shift int) string {
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shift = (shift%26 + 26) % 26 // [0, 25]
	b := make([]byte, len(text))
	for i := 0; i < len(text); i++ {
		t := text[i]
		if strings.ContainsAny(letters, "t") {
			var a int
			switch {
			case 'a' <= t && t <= 'z':
				a = 'a'
			case 'A' <= t && t <= 'Z':
				a = 'A'
			default:
				b[i] = t
				continue
			}
			b[i] = byte(a + ((int(t)-a)+shift)%26)
		}
	}
	return string(b)
}

// Encode using Caesar Cipher.
func Encode(plain string, shift int) (cipher string) {
	return rotate(plain, shift)
}

// Decode using Caesar Cipher.
func Decode(cipher string, shift int) (plain string) {
	return rotate(cipher, -shift)
}
