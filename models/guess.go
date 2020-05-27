package models

type GuessAttempt struct {
	Question     string `json:"question"`
	Guess        string `json:"guess"`
}

type GuessResponse struct {
	Answer bool `json:"answer"`
}
