package main

import (
	"encoding/json"
	http "net/http"
)

// CreateFAQ
// @Description Creation of a new FAQ
// @ID create-faq
// @Accept json
// @Router /api/admin/faqs [put]
// @Tag faq create put admin
func CreateFAQ(r *http.Request, ar *APIResponse) {
	newFAQ := FAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newFAQ)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}
	if !isCorrectLanguage(newFAQ.Language) {
		ar.setErrorLanguage(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`INSERT INTO FAQ (question, answer, language) VALUES (?, ?, ?)`, newFAQ.Question, newFAQ.Answer, newFAQ.Language)
	if err != nil {
		ar.setErrorDBInsert(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}

	ar.setStatus(StatusCreated)
}

// RETRIEVE
func retrieveFAQs(r *http.Request, ar *APIResponse) {
	lang, err := parseLanguage(r)
	if err != nil {
		ar.setErrorLanguage(err)
		return
	}

	faqs := []FAQ{}
	rows, err := db.Query(`SELECT question, answer FROM FAQ WHERE language = ?`, lang)
	if err != nil {
		ar.setErrorDBSelect(err)
		return
	}

	for rows.Next() {
		var question, answer string
		err = rows.Scan(&question, &answer)
		if err != nil {
			ar.setErrorDBScan(err)
			return
		}
		faqs = append(faqs, FAQ{question, answer, lang})
	}

	if err = rows.Err(); err != nil {
		ar.setErrorDBAfter(err)
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
		ar.setErrorJSON(err)
		return
	}

	if !isCorrectLanguage(newFAQ.FAQ.Language) {
		ar.setErrorLanguage(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`UPDATE FAQ SET language = ?, question = ?, answer = ? WHERE question = ?`,
		newFAQ.FAQ.Language, newFAQ.FAQ.Question, newFAQ.FAQ.Answer, newFAQ.Question)
	if err != nil {
		ar.setErrorDBUpdate(err, tx)
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}
	ar.setStatus(StatusUpdated)
}

// DELETE
func deleteFAQ(r *http.Request, ar *APIResponse) {
	faq := FAQ{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&faq)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`DELETE FROM FAQ WHERE question = ?`, faq.Question)
	if err != nil {
		ar.setErrorDBDelete(err, tx)
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setErrorDBCommit(err)
		return
	}
	ar.setStatus(StatusDeleted)
}
