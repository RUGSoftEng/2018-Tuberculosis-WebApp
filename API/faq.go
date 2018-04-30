package main

import (
	"github.com/pkg/errors"
	http "net/http"
)

func getFAQs(r *http.Request, ar *APIResponse) {
	faqs := []FAQ{}
	rows, err := db.Query(`SELECT question, answer FROM FAQ`)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error during query")
		return
	}

	for rows.Next() {
		var question, answer string
		err = rows.Scan(&question, &answer)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		faqs = append(faqs, FAQ{question, answer})
	}

	if err = rows.Err(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	ar.Data = faqs
}
