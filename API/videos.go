package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
func addVideo(r *http.Request, ar *APIResponse) {
	video := Video{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&video)
	if err != nil {
		ar.StatusCode = http.StatusBadRequest
		ar.Error = errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to start new transaction")
		return
	}
	_, err = tx.Exec(`INSERT INTO Videos (topic, title, reference) VALUES (?, ?, ?)`,
		video.Topic, video.Title, video.Reference)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		err = errors.Wrap(err, "Failed to insert video into the database")
		ar.Error = errorWithRollback(err, tx)
		return
	}

	if err = tx.Commit(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Failed to commit changes to database.")
		return
	}
	ar.StatusCode = http.StatusCreated
}

// RETRIEVE
func getTopics(r *http.Request, ar *APIResponse) {
	rows, err := db.Query(`SELECT DISTINCT topic FROM Videos`)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error when querying the database")
		return
	}

	var topics []string
	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		topics = append(topics, topic)
	}
	if err = rows.Err(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}
	ar.Data = topics
}

// RETRIEVE
func getVideoByTopic(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	rows, err := db.Query(`SELECT topic, title, reference FROM Videos WHERE topic = ?`, topic)
	if err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error when querying the database")
		return
	}

	videos := []Video{}
	for rows.Next() {
		var topic, title, reference string
		err = rows.Scan(&topic, &title, &reference)
		if err != nil {
			ar.StatusCode = http.StatusInternalServerError
			ar.Error = errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		videos = append(videos, Video{topic, title, reference})
	}
	if err = rows.Err(); err != nil {
		ar.StatusCode = http.StatusInternalServerError
		ar.Error = errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	ar.Data = videos
}
