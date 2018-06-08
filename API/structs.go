package main

// ScheduledDosage : Describes an instance of a dosage for the schedule
type ScheduledDosage struct {
	Dosage Dosage `json:"dosage"`
	Day    string `json:"date"`
	Taken  bool   `json:"taken"`
}

// InputScheduledDosage : Struct for handling the input for scheduled dosages
type InputScheduledDosage struct {
	Medicine Medicine `json:"medicine"`
	Days     []string `json:"days"`
}

// Dosage : Describes the time and number of pills is associated with a medicine
// This can be different for different people
type Dosage struct {
	IntakeIntervalStart string   `json:"intake_interval_start"`
	IntakeIntervalEnd   string   `json:"intake_interval_end"`
	NumberOfPills       int      `json:"amount"`
	Medicine            Medicine `json:"medicine"`
}

// UpdateDosage : Used for updating the dosage: updates old with new
type UpdateDosage struct {
	OldDosage Dosage `json:"old_dosage"`
	NewDosage Dosage `json:"new_dosage"`
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

// NoteReturn : A object for returning Notes
type NoteReturn struct {
	ID        int    `json:"id"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

// Video : A video with their reference
type Video struct {
	Topic     string `json:"topic"`
	Title     string `json:"title"`
	Reference string `json:"reference"`
	Language  string `json:"language"`
}

// UpdateVideo : Struct used for updating Video. Specifies the old video (identifier)
//  + what the new video should be.
type UpdateVideo struct {
	Topic string `json:"topic"`
	Title string `json:"title"`
	Video Video  `json:"video"`
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

// ReferencedQuiz : Used for requests where a quiz is needed with its reference video (e.g. creating or deleting)
type ReferencedQuiz struct {
	Video Video `json:"video"`
	Quiz  Quiz  `json:"quiz"`
}

// UpdateQuiz : Used when updating the quiz.
// Contains original question + video reference + updated quiz
type UpdateQuiz struct {
	Video    Video  `json:"video"`
	Question string `json:"question"`
	Quiz     Quiz   `json:"quiz"`
}

// FAQ : Describes a Frequently Asked Question
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Language string `json:"language"`
}

// UpdateFAQ : Struct used for updating FAQ. Specifies the old question (identifier)
//  + what the new faq should be.
type UpdateFAQ struct {
	Question string `json:"question"`
	FAQ      FAQ    `json:"faq"`
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

// Patient : Information of a patient
type Patient struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
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
