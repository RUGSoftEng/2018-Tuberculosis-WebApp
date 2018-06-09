package main

import (
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	http "net/http"
	"strconv"
	"strings"
)

// Retrieves the id variable from the url + converts the variable to an integer
func getPatientIDVariable(r *http.Request) (patientID int, err error) {
	patientID, err = strconv.Atoi(mux.Vars(r)["id"])
	return
}

// Retrieves the value of the requested URL variable. Gives an error when variable does not exist
func getURLVariable(r *http.Request, variable string) (v string, err error) {
	v = mux.Vars(r)[variable]
	err = nil
	if v == "" {
		err = errors.New(ErrEmptyVariable)
	}
	return
}

// Calls getURLVariable function and converts the returned string into and integer
func getURLVariableInt(r *http.Request, variable string) (v int, err error) {
	str, err := getURLVariable(r, variable)
	if err != nil {
		return -1, err
	}
	v, err = strconv.Atoi(str)
	return
}

// Function which retrieves the id number from the given medicine
// Gives error when something goes wrong, e.g. when the medicine doesn't exist
func queryMedicineID(medicine Medicine) (medicineID int, err error) {
	medicineID = -1
	row := db.QueryRow(`SELECT id FROM Medicines WHERE med_name = ?`,
		medicine.Name)
	if err != nil {
		return
	}
	err = row.Scan(&medicineID)
	if err != nil {
		return
	}
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
	row := db.QueryRow(`SELECT id FROM Dosages WHERE patient_id = ? AND medicine_id = ?`,
		patientID, medicineID)
	if err != nil {
		return
	}
	err = row.Scan(&dosageID)
	if err != nil {
		return
	}
	return
}

func queryVideoID(video Video) (videoID int, err error) {
	videoID = -1
	row := db.QueryRow(`SELECT id FROM Videos WHERE topic = ? AND title = ?`, video.Topic, video.Title)
	if err != nil {
		return
	}
	err = row.Scan(&videoID)
	if err != nil {
		return
	}
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
