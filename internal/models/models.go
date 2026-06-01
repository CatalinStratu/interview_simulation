package models

import (
	"time"
)

type Question struct {
	ID               int       `json:"id"`
	QuestionText     string    `json:"question_text"`
	QuestionType     string    `json:"question_type"`
	DifficultyLevel  string    `json:"difficulty_level"`
	ExpectedDuration int       `json:"expected_duration"`
	Position         int       `json:"position"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsActive         bool      `json:"is_active"`
}

// ReorderQuestionsRequest carries the new top-to-bottom ordering of question
// ids set from the admin UI.
type ReorderQuestionsRequest struct {
	IDs []int `json:"ids"`
}

type InterviewSession struct {
	ID                int        `json:"id"`
	SessionCode       string     `json:"session_code"`
	CandidateName     *string    `json:"candidate_name,omitempty"`
	CandidateEmail    *string    `json:"candidate_email,omitempty"`
	Status            string     `json:"status"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	TotalQuestions    int        `json:"total_questions"`
	QuestionsAnswered int        `json:"questions_answered"`
}

type VoiceAnswer struct {
	ID            int       `json:"id"`
	SessionID     int       `json:"session_id"`
	QuestionID    int       `json:"question_id"`
	AudioFilePath string    `json:"audio_file_path"`
	AudioDuration *int      `json:"audio_duration,omitempty"`
	Transcription *string   `json:"transcription,omitempty"`
	AnsweredAt    time.Time `json:"answered_at"`
}

type Remark struct {
	ID         int       `json:"id"`
	SessionID  int       `json:"session_id"`
	RemarkText string    `json:"remark_text"`
	Rating     *int      `json:"rating,omitempty"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SessionQuestion struct {
	ID            int       `json:"id"`
	SessionID     int       `json:"session_id"`
	QuestionID    int       `json:"question_id"`
	QuestionOrder int       `json:"question_order"`
	AskedAt       time.Time `json:"asked_at"`
}

// DTOs for API requests/responses
type CreateQuestionRequest struct {
	QuestionText     string `json:"question_text"`
	QuestionType     string `json:"question_type"`
	DifficultyLevel  string `json:"difficulty_level"`
	ExpectedDuration int    `json:"expected_duration"`
}

type StartSessionRequest struct {
	CandidateName  string `json:"candidate_name,omitempty"`
	CandidateEmail string `json:"candidate_email,omitempty"`
	NumQuestions   int    `json:"num_questions"`
}

type AddRemarkRequest struct {
	RemarkText string `json:"remark_text"`
	Rating     *int   `json:"rating,omitempty"`
}

type NextQuestionResponse struct {
	Question      *Question `json:"question,omitempty"`
	QuestionIndex int       `json:"question_index"`
	TotalQuestions int      `json:"total_questions"`
	IsComplete    bool      `json:"is_complete"`
}
