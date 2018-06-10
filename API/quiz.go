package main

import (
	"encoding/json"
	http "net/http"
	"strings"
)

// CREATE
func createQuiz(r *http.Request, ar *APIResponse) {
	newQuiz := ReferencedQuiz{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newQuiz)
	if err != nil {
		ar.setErrorAndStatus(http.StatusBadRequest, err, "Unexpected error during JSON decoding")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorAndStatus(http.StatusInternalServerError, err, "Failed to start new transaction")
		return
	}

	videoID, err := queryVideoID(newQuiz.Video)
	if err != nil {
		ar.setError(err, "Error during querying video")
	}

	unseparatedAnswers := strings.Join(newQuiz.Quiz.Answers, ":")
	_, err = tx.Exec(`INSERT INTO Quizzes (video, question, answers) VALUES (?, ?, ?)`,
		videoID, newQuiz.Quiz.Question, unseparatedAnswers)
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

// RETRIEVE: getVideoByTopic

// UPDATE
func updateQuiz(r *http.Request, ar *APIResponse) {
	newQuiz := UpdateQuiz{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newQuiz)
	if err != nil {
		ar.setError(err, "Error during JSON parse, expected an UpdateVideo struct")
		return
	}

	videoID, err := queryVideoID(newQuiz.Video)
	if err != nil {
		ar.setError(err, "Error during querying video")
	}
	unseparatedAnswers := strings.Join(newQuiz.Quiz.Answers, ":")

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`UPDATE Quizzes SET question = ?, answers= ? WHERE video = ? AND question = ?`,
		newQuiz.Quiz.Question, unseparatedAnswers, videoID, newQuiz.Question)
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
func deleteQuiz(r *http.Request, ar *APIResponse) {
	quiz := ReferencedQuiz{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&quiz)
	if err != nil {
		ar.setError(err, "Error during JSON parse, expected an UpdateVideo struct")
		return
	}

	videoID, err := queryVideoID(quiz.Video)
	if err != nil {
		ar.setError(err, "Error during querying video")
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setError(err, "Failed to start transaction.")
		return
	}

	_, err = tx.Exec(`DELETE FROM Quizzes WHERE video = ? AND question = ?`, videoID, quiz.Quiz.Question)
	if err != nil {
		ar.setError(errorWithRollback(err, tx), "Database failure")
		return
	}
	if err = tx.Commit(); err != nil {
		ar.setError(err, "Failed to commit changes to database.")
		return
	}
}
