package main

// Dosage : Descibes the dosage of a certain medicine for a certain day
type Dosage struct {
	Day           string `json:"date"`
	IntakeMoment  string `json:"intake_moment"`
	NumberOfPills int    `json:"amount"`
	Medicine      string `json:"medicine"`
	Taken         bool   `json:"taken"`
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
}
