package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
func createFAQ(r *http.Request, ar *APIResponse) {
	newFAQ := FAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newFAQ)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode JSON.")
		return
	}
	if !isCorrectLanguage(newFAQ.Language) {
		ar.setErrorAndStatus(http.StatusBadRequest, errors.New(""), "Invalid Language")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`INSERT INTO FAQ (question, answer, language) VALUES (?, ?, ?)`, newFAQ.Question, newFAQ.Answer, newFAQ.Language)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}

	ar.setStatus(http.StatusCreated)
}

// RETRIEVE
func getFAQs(r *http.Request, ar *APIResponse) {
	lang, err := parseLanguage(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "")
		return
	}

	faqs := []FAQ{}
	rows, err := db.Query(`SELECT question, answer FROM FAQ WHERE language = ?`, lang)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during query")
		return
	}

	for rows.Next() {
		var question, answer string
		err = rows.Scan(&question, &answer)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		faqs = append(faqs, FAQ{question, answer, lang})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}

	ar.setResponse(faqs)
}
