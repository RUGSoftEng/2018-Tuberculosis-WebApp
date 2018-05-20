package main

import (
	"encoding/json"
	http "net/http"
)

// CREATE
func addFAQ(r *http.Request, ar *APIResponse) {
	newFAQ := FAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newFAQ)
	if err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Failed to decode JSON.")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`INSERT INTO FAQ (question, answer) VALUES (?, ?)`, newFAQ.Question, newFAQ.Answer)
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
	faqs := []FAQ{}
	rows, err := db.Query(`SELECT question, answer FROM FAQ`)
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
		faqs = append(faqs, FAQ{question, answer})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}

	ar.setResponse(faqs)
}
