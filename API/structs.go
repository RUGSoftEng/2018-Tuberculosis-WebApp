package main

//import "time"

type Dosage struct {
	Day string `json:"date"`
	IntakeMoment string `json:"intake_moment"`
	NumberOfPills int `json:"amount"`
	Medicine string `json:"medicine"`
	Taken bool `json:"taken"`
}

type Note struct {
	Note string `json:"note"`
	CreatedAt string `json:"created_at"`
}

type Patient struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	ApiToken string `json:"api_token"`
}

type Physician struct {
	Username string `json:"username"`
	Name string `json:"name"`
	Password string `json:"password"`
	ApiToken string `json:"api_token"`
	Email string `json:"email"`
	CreationToken string `json:"creation_token"`
}
