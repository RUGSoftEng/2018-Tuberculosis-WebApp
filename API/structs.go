package main

import "time"

type Dosage struct {
	Day string `json:"date"`
	IntakeMoment string `json:"intake_moment"`
	NumberOfPills int `json:"amount"`
	Medicine string `json:"medicine"`
	Taken bool `json:"taken"`
}

type Note struct {
	Note string `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type Patient struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	ApiToken string `json:"api_token"`
	Dosages []Dosage `json:"dosages"`
	Notes []Note `json:"note"`
}

type Physician struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	ApiToken string `json:"api_token"`
	Email string `json:"email"`
	CreationToken string `json:"creation_token"`
}
