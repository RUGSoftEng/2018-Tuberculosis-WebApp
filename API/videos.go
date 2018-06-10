package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql" // anonymous import
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
)

// CREATE
func createVideo(r *http.Request, ar *APIResponse) {
	video := Video{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&video)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Unexpected error during JSON decoding")
		return
	}

	if !isCorrectLanguage(video.Language) {
		ar.setErrorAndStatus(http.StatusBadRequest, errors.New(""), "Invalid Language")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start new transaction")
		return
	}
	_, err = tx.Exec(`INSERT INTO Videos (language, topic, title, reference) VALUES (?, ?, ?, ?)`,
		video.Language, video.Topic, video.Title, video.Reference)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, errorWithRollback(err, tx), "Database failure")
		return
	}

	if err = tx.Commit(); err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to commit changes to database.")
		return
	}
	ar.StatusCode = http.StatusCreated
}

// RETRIEVE
func retrieveTopics(r *http.Request, ar *APIResponse) {
	lang, err := parseLanguage(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "")
		return
	}
	rows, err := db.Query(`SELECT DISTINCT topic FROM Videos WHERE language = ?`, lang)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error when querying the database")
		return
	}

	var topics []string
	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error during row scanning")
			return
		}
		topics = append(topics, topic)
	}
	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}
	ar.setResponse(topics)
}

// RETRIEVE
func retrieveVideoByTopic(r *http.Request, ar *APIResponse) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	lang, err := parseLanguage(r)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "")
		return
	}

	rows, err := db.Query(`SELECT id, topic, title, reference FROM Videos WHERE topic = ? AND language = ?`,
		topic, lang)
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error when querying the database")
		return
	}

	videos := []VideoQuiz{}
	for rows.Next() {
		var id int
		var topic, title, reference string
		err = rows.Scan(&id, &topic, &title, &reference)
		if err != nil {
			ar.setErrorAndStatus(http.StatusInternalServerError, err, "Unexpected error during row scanning")
			return
		}
		video := Video{topic, title, reference, lang}
		quizzes, err := queryQuizzes(id)
		if err != nil {
			ar.setError(err, "Error during querying quizzes")
		}
		videos = append(videos, VideoQuiz{video, quizzes})
	}
	if err = rows.Err(); err != nil {
		ar.setErrorAndStatus(StatusFailedOperation, err, "Unexpected error after scanning rows")
		return
	}
	ar.setResponse(videos)
}

// UPDATE
func updateVideo(r *http.Request, ar *APIResponse) {
	newVideo := UpdateVideo{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newVideo)
	if err != nil {
		ar.setError(err, "Error during JSON parse, expected an UpdateVideo struct")
		return
	}

	if !isCorrectLanguage(newVideo.Video.Language) {
		ar.setErrorAndStatus(http.StatusBadRequest, errors.New(""), "Invalid Language "+newVideo.Video.Language)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`UPDATE Videos SET topic = ?, title = ?, reference = ?, language = ? 
                 WHERE topic = ? AND title = ?`,
		newVideo.Video.Topic, newVideo.Video.Title, newVideo.Video.Reference, newVideo.Video.Language,
		newVideo.Topic, newVideo.Title)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Database failure")
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
		return
	}
}

// DELETE
func deleteVideo(r *http.Request, ar *APIResponse) {
	video := Video{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&video)
	if err != nil {
		ar.setError(err, "Error during JSON parse, expected a Video struct")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`DELETE FROM Videos WHERE topic = ? AND title = ?`,
		video.Topic, video.Title)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Database failure")
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
		return
	}
}
