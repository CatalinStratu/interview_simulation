package services

import (
	"database/sql"
	"fmt"
	"interview-chatbot/internal/models"
	"math/rand"
	"time"
)

type SessionService struct {
	db              *sql.DB
	questionService *QuestionService
}

func NewSessionService(db *sql.DB, questionService *QuestionService) *SessionService {
	return &SessionService{
		db:              db,
		questionService: questionService,
	}
}

func (s *SessionService) StartSession(req models.StartSessionRequest) (*models.InterviewSession, error) {
	// Generate unique session code
	sessionCode := generateSessionCode()

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Always use every active question, of every type and difficulty, served in
	// the fixed order they are set in the DB (by id) — one after another, never
	// random and never limited to a subset. NumQuestions is ignored on purpose.
	questions, err := s.questionService.GetOrderedQuestions(0)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("no active questions available")
	}

	// Create session
	query := `INSERT INTO interview_sessions (session_code, candidate_name, candidate_email, total_questions)
	          VALUES (?, ?, ?, ?)`

	var candidateName, candidateEmail *string
	if req.CandidateName != "" {
		candidateName = &req.CandidateName
	}
	if req.CandidateEmail != "" {
		candidateEmail = &req.CandidateEmail
	}

	result, err := tx.Exec(query, sessionCode, candidateName, candidateEmail, len(questions))
	if err != nil {
		return nil, err
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert session questions
	for i, q := range questions {
		query := `INSERT INTO session_questions (session_id, question_id, question_order)
		          VALUES (?, ?, ?)`
		_, err := tx.Exec(query, sessionID, q.ID, i+1)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetSessionByID(int(sessionID))
}

func (s *SessionService) GetSessionByID(id int) (*models.InterviewSession, error) {
	query := `SELECT id, session_code, candidate_name, candidate_email, status,
	          started_at, completed_at, total_questions, questions_answered
	          FROM interview_sessions WHERE id = ?`

	var session models.InterviewSession
	err := s.db.QueryRow(query, id).Scan(&session.ID, &session.SessionCode,
		&session.CandidateName, &session.CandidateEmail, &session.Status,
		&session.StartedAt, &session.CompletedAt, &session.TotalQuestions,
		&session.QuestionsAnswered)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) GetSessionByCode(code string) (*models.InterviewSession, error) {
	query := `SELECT id, session_code, candidate_name, candidate_email, status,
	          started_at, completed_at, total_questions, questions_answered
	          FROM interview_sessions WHERE session_code = ?`

	var session models.InterviewSession
	err := s.db.QueryRow(query, code).Scan(&session.ID, &session.SessionCode,
		&session.CandidateName, &session.CandidateEmail, &session.Status,
		&session.StartedAt, &session.CompletedAt, &session.TotalQuestions,
		&session.QuestionsAnswered)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionService) GetNextQuestion(sessionID int) (*models.NextQuestionResponse, error) {
	session, err := s.GetSessionByID(sessionID)
	if err != nil {
		return nil, err
	}

	if session.Status != "in_progress" {
		return &models.NextQuestionResponse{
			Question:       nil,
			QuestionIndex:  session.QuestionsAnswered,
			TotalQuestions: session.TotalQuestions,
			IsComplete:     true,
		}, nil
	}

	// Get the next unanswered question
	query := `SELECT q.id, q.question_text, q.question_type, q.difficulty_level,
	          q.expected_duration, q.created_at, q.updated_at, q.is_active
	          FROM session_questions sq
	          JOIN questions q ON sq.question_id = q.id
	          LEFT JOIN voice_answers va ON va.session_id = sq.session_id AND va.question_id = sq.question_id
	          WHERE sq.session_id = ? AND va.id IS NULL
	          ORDER BY sq.question_order
	          LIMIT 1`

	var question models.Question
	err = s.db.QueryRow(query, sessionID).Scan(&question.ID, &question.QuestionText,
		&question.QuestionType, &question.DifficultyLevel, &question.ExpectedDuration,
		&question.CreatedAt, &question.UpdatedAt, &question.IsActive)

	if err == sql.ErrNoRows {
		// No more questions
		return &models.NextQuestionResponse{
			Question:       nil,
			QuestionIndex:  session.QuestionsAnswered,
			TotalQuestions: session.TotalQuestions,
			IsComplete:     true,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return &models.NextQuestionResponse{
		Question:       &question,
		QuestionIndex:  session.QuestionsAnswered + 1,
		TotalQuestions: session.TotalQuestions,
		IsComplete:     false,
	}, nil
}

func (s *SessionService) CompleteSession(sessionID int) error {
	query := `UPDATE interview_sessions SET status = 'completed', completed_at = NOW()
	          WHERE id = ?`
	_, err := s.db.Exec(query, sessionID)
	return err
}

func (s *SessionService) GetAllSessions() ([]models.InterviewSession, error) {
	query := `SELECT id, session_code, candidate_name, candidate_email, status,
	          started_at, completed_at, total_questions, questions_answered
	          FROM interview_sessions ORDER BY started_at DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.InterviewSession
	for rows.Next() {
		var session models.InterviewSession
		err := rows.Scan(&session.ID, &session.SessionCode, &session.CandidateName,
			&session.CandidateEmail, &session.Status, &session.StartedAt,
			&session.CompletedAt, &session.TotalQuestions, &session.QuestionsAnswered)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func generateSessionCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	rand.Seed(time.Now().UnixNano())
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
