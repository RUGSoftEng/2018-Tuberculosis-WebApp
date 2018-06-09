package main

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
	"strconv"
	"strings"
)

// Retrieves the value of the requested URL variable. Gives an error when variable does not exist
func getURLVariable(r *http.Request, variable string) (v string, err error) {
	v = mux.Vars(r)[variable]
	err = nil
	if v == "" {
		err = errors.New(ErrEmptyVariable)
	}
	return
}

func getURLVariableInt(r *http.Request, variable string) (id int, err error) {
	id = -1
	v, err := getURLVariable(r, variable)
	if err != nil {
		return
	}
	id, err = strconv.Atoi(v)
	return
}

// Function which retrieves the id number from the given medicine
// Gives error when something goes wrong, e.g. when the medicine doesn't exist
func queryMedicineID(medicine Medicine) (medicineID int, err error) {
	medicineID = -1
	err = db.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		medicine.Name).Scan(&medicineID)
	return
}

// Function which retrieves the id number from the given dosage + patient
// Gives error when something goes wrong, e.g. when the dosage doesn't exist
func queryDosageID(dosage Dosage, patientID int) (dosageID int, err error) {
	dosageID = -1

	// Scan Medicine ID
	medicineID, err := queryMedicineID(dosage.Medicine)
	if err != nil {
		return
	}

	// Scan Dosage ID
	err = db.QueryRow(`SELECT id FROM Dosages WHERE patient_id = ? AND medicine_id = ?`,
		patientID, medicineID).Scan(&dosageID)
	return
}

func queryVideoID(video Video) (videoID int, err error) {
	videoID = -1
	err = db.QueryRow(`SELECT id FROM Videos WHERE topic = ? AND title = ?`,
		video.Topic, video.Title).Scan(&videoID)
	return
}

func queryQuizzes(videoID int) (quizzes []Quiz, err error) {
	rows, err := db.Query(`SELECT question, answers FROM Quizzes WHERE video = ?`, videoID)
	if err != nil {
		return
	}
	for rows.Next() {
		var question, answers string
		err = rows.Scan(&question, &answers)
		if err != nil {
			return
		}
		splittedAnswers := strings.Split(answers, ":")
		quizzes = append(quizzes, Quiz{question, splittedAnswers})
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
