package models

// Answer represents an answer given by a user to a question.
type Answer struct {
	ID         int64  // Unique identifier for the answer
	QuestionID int64  // Identifier of the associated question
	Text       string // Text of the answer
}

// AnswerSubmission represents a submitted answer with JSON tags.
type AnswerSubmission struct {
	QuestionID int64  `json:"question_id"` // Identifier of the associated question
	Text       string `json:"text"`        // Text of the answer
}

// SignAnswersRequest represents a request to sign answers.
type SignAnswersRequest struct {
	Answers []AnswerSubmission `json:"answers"` // List of submitted answers
}
