package main

import (
	"database/sql"
	"github.com/pkg/errors"
)

// ScheduledDosage : Describes an instance of a dosage for the schedule
type ScheduledDosage struct {
	Dosage Dosage `json:"dosage"`
	Day    string `json:"date"`
	Taken  bool   `json:"taken"`
}

// Dosage : Describes the time and number of pills is associated with a medicine
// This can be different for different people
type Dosage struct {
	IntakeMoment  string   `json:"intake_moment"`
	NumberOfPills int      `json:"amount"`
	Medicine      Medicine `json:"medicine"`
}

// Medicine : Data for a medicine
type Medicine struct {
	Name string `json:"name"`
}

// Note : A note from a patient for their physician
type Note struct {
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

// Patient : Information of a patient
type Patient struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	APIToken string `json:"api_token"`
}

// Physician : Information of a patient
type Physician struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	APIToken      string `json:"api_token"`
	Email         string `json:"email"`
	CreationToken string `json:"creation_token"`
}

// Video : A video with their reference
type Video struct {
	Topic     string `json:"topic"`
	Title     string `json:"title"`
	Reference string `json:"reference"`
	Language  string `json:"language"`
}

// VideoQuiz : The video alongside it's paired quizzes
type VideoQuiz struct {
	Video   Video  `json:"video"`
	Quizzes []Quiz `json:"quizzes"`
}

// Quiz : The quiz, belongs to a video.
// contains a list of answers, the first answer is always the correct answer.
type Quiz struct {
	Question string   `json:"question"`
	Answers  []string `json:"answers"`
}

// FAQ : Describes a Frequently Asked Question
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
}

// UserValidation : A set of values needed for authenticate a user
type UserValidation struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWToken : Contains the ingredients to access restricted areas
type JWToken struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

// PatientInfo : Identifies a patient through his/her public data
type PatientInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PatientOverview : A complete overview of a patient
type PatientOverview struct {
	Username       string `json:"username"`
	Name           string `json:"name"`
	PhysicianName  string `json:"physician_name"`
	PhysicianEmail string `json:"email"`
}

// PhysicianOverview : A complete overview of a physician
type PhysicianOverview struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func errorWithRollback(err error, tx *sql.Tx) error {
	err2 := tx.Rollback()
	if err2 != nil {
		err = errors.New(err.Error() + "\n Rollback failed:" + err2.Error())
	}
	return err
}
