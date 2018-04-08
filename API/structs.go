package main

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
}

// UserValidation : ADD DOCUMENTATION
type UserValidation struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWToken : ADD DOCUMENTATION
type JWToken struct {
	Token string `json:"token"`
}

// APIResponse : Type used by the Response Channel
// in the handlerWrapper (does not need json tags)
type APIResponse struct {
	Data       interface{}
	StatusCode int
}
