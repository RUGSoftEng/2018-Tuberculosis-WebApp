package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
func addVideo(r *http.Request, responseChan chan []byte, errorChan chan error) {
	video := Video{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&video)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error during JSON decoding")
		return
	}

	trans, err := db.Begin()
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to start new transaction")
		return
	}
	_, err = trans.Exec(`INSERT INTO Videos (topic, title, reference) VALUES (?, ?, ?)`,
		video.Topic, video.Title, video.Reference)
	if err != nil {
		errorChan <- errors.Wrap(err, "Failed to insert video into the database")
		return
	}

	errorChan <- trans.Commit()
	return
}

// RETRIEVE
func getTopics(r *http.Request, responseChan chan []byte, errorChan chan error) {
	rows, err := db.Query(`SELECT DISTINCT topic FROM Videos`)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when querying the database")
		return
	}

	var topics []string
	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		topics = append(topics, topic)
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	jsonValues, err := json.Marshal(topics)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when converting to JSON")
		return
	}
	responseChan <- jsonValues
	return
}

// RETRIEVE
func getVideoByTopic(r *http.Request, responseChan chan []byte, errorChan chan error) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	rows, err := db.Query(`SELECT topic, title, reference FROM Videos WHERE topic = ?`, topic)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when querying the database")
		return
	}

	videos := []Video{}
	for rows.Next() {
		var topic, title, reference string
		err = rows.Scan(&topic, &title, &reference)
		if err != nil {
			errorChan <- errors.Wrap(err, "Unexpected error during row scanning")
			return
		}
		videos = append(videos, Video{topic, title, reference})
	}
	if err = rows.Err(); err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error after scanning rows")
		return
	}

	jsonValues, err := json.Marshal(videos)
	if err != nil {
		errorChan <- errors.Wrap(err, "Unexpected error when converting to JSON")
		return
	}
	responseChan <- jsonValues
	errorChan <- nil
	return
}
