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
		ar.setErrorAndStatus(StatusInvalidJSON, err, ErrDecodeJSON)
		return
	}
	if !isCorrectLanguage(newFAQ.Language) {
		ar.setErrorAndStatus(StatusInvalidLanguage, errors.New(ErrLang))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBTransactionStartFaillure)
		return
	}

	_, err = tx.Exec(`INSERT INTO FAQ (question, answer, language) VALUES (?, ?, ?)`, newFAQ.Question, newFAQ.Answer, newFAQ.Language)
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, errorWithRollback(err, tx), ErrDBInsert)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBCommit)
		return
	}

	ar.setStatus(StatusCreated)
}

// RETRIEVE
func retrieveFAQs(r *http.Request, ar *APIResponse) {
	lang, err := parseLanguage(r)
	if err != nil {
		ar.setErrorAndStatus(StatusInvalidLanguage, err)
		return
	}

	faqs := []FAQ{}
	rows, err := db.Query(`SELECT question, answer FROM FAQ WHERE language = ?`, lang)
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBSelect)
		return
	}

	for rows.Next() {
		var question, answer string
		err = rows.Scan(&question, &answer)
		if err != nil {
			ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBScan)
			return
		}
		faqs = append(faqs, FAQ{question, answer, lang})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBAfter)
		return
	}

	ar.setResponse(faqs)
}

// UPDATE
func updateFAQ(r *http.Request, ar *APIResponse) {
	newFAQ := UpdateFAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newFAQ)
	if err != nil {
		ar.setErrorAndStatus(StatusInvalidJSON, err, ErrDecodeJSON)
		return
	}

	if !isCorrectLanguage(newFAQ.FAQ.Language) {
		ar.setErrorAndStatus(StatusInvalidLanguage, errors.New(ErrLang))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBTransactionStartFaillure)
		return
	}

	_, err = tx.Exec(`UPDATE FAQ SET language = ?, question = ?, answer = ? WHERE question = ?`,
		newFAQ.FAQ.Language, newFAQ.FAQ.Question, newFAQ.FAQ.Answer, newFAQ.Question)
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, errorWithRollback(err, tx), ErrDBUpdate)
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBCommit)
		return
	}
}

// DELETE
func deleteFAQ(r *http.Request, ar *APIResponse) {
	faq := FAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&faq)
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDecodeJSON)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBTransactionStartFaillure)
		return
	}

	_, err = tx.Exec(`DELETE FROM FAQ WHERE question = ?`, faq.Question)
	if err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, errorWithRollback(err, tx), ErrDBDelete)
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(StatusDatabaseError, err, ErrDBCommit)
		return
	}
}
