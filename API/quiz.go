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
		ar.setErrorJSON(err)
		return
	}

	videoID, err := queryVideoID(newQuiz.Video)
	if err != nil {
		ar.setErrorDBSelect(err)
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	unseparatedAnswers := strings.Join(newQuiz.Quiz.Answers, ":")
	_, err = tx.Exec(`INSERT INTO Quizzes (video, question, answers) VALUES (?, ?, ?)`,
		videoID, newQuiz.Quiz.Question, unseparatedAnswers)
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

// RETRIEVE: getVideoByTopic

// UPDATE
func updateQuiz(r *http.Request, ar *APIResponse) {
	newQuiz := UpdateQuiz{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&newQuiz)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	videoID, err := queryVideoID(newQuiz.Video)
	if err != nil {
		ar.setErrorDBSelect(err)
	}
	unseparatedAnswers := strings.Join(newQuiz.Quiz.Answers, ":")

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`UPDATE Quizzes SET question = ?, answers= ? WHERE video = ? AND question = ?`,
		newQuiz.Quiz.Question, unseparatedAnswers, videoID, newQuiz.Question)
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
func deleteQuiz(r *http.Request, ar *APIResponse) {
	quiz := ReferencedQuiz{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&quiz)
	if err != nil {
		ar.setErrorJSON(err)
		return
	}

	videoID, err := queryVideoID(quiz.Video)
	if err != nil {
		ar.setErrorDBSelect(err)
	}

	tx, err := db.Begin()
	if err != nil {
		ar.setErrorDBBegin(err)
		return
	}

	_, err = tx.Exec(`DELETE FROM Quizzes WHERE video = ? AND question = ?`, videoID, quiz.Quiz.Question)
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
