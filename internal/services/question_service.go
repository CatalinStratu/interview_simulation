package services

import (
	"database/sql"
	"fmt"
	"interview-chatbot/internal/models"
)

type QuestionService struct {
	db *sql.DB
}

func NewQuestionService(db *sql.DB) *QuestionService {
	return &QuestionService{db: db}
}

func (s *QuestionService) GetAllQuestions() ([]models.Question, error) {
	query := `SELECT id, question_text, question_type, difficulty_level,
	          expected_duration, position, created_at, updated_at, is_active
	          FROM questions WHERE is_active = TRUE ORDER BY position ASC, id ASC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, &q.DifficultyLevel,
			&q.ExpectedDuration, &q.Position, &q.CreatedAt, &q.UpdatedAt, &q.IsActive)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}

func (s *QuestionService) GetQuestionByID(id int) (*models.Question, error) {
	query := `SELECT id, question_text, question_type, difficulty_level,
	          expected_duration, position, created_at, updated_at, is_active
	          FROM questions WHERE id = ?`

	var q models.Question
	err := s.db.QueryRow(query, id).Scan(&q.ID, &q.QuestionText, &q.QuestionType,
		&q.DifficultyLevel, &q.ExpectedDuration, &q.Position, &q.CreatedAt, &q.UpdatedAt, &q.IsActive)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("question not found")
	}
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (s *QuestionService) CreateQuestion(req models.CreateQuestionRequest) (*models.Question, error) {
	// New questions go to the end of the list.
	query := `INSERT INTO questions (question_text, question_type, difficulty_level, expected_duration, position)
	          VALUES (?, ?, ?, ?, (SELECT next_pos FROM (SELECT COALESCE(MAX(position), 0) + 1 AS next_pos FROM questions) AS t))`

	result, err := s.db.Exec(query, req.QuestionText, req.QuestionType,
		req.DifficultyLevel, req.ExpectedDuration)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetQuestionByID(int(id))
}

func (s *QuestionService) DeleteQuestion(id int) error {
	query := `UPDATE questions SET is_active = FALSE WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}

// UpdateQuestionOrder persists a new top-to-bottom ordering. The given ids are
// assigned positions 1..N in the order received; the same order is then used by
// both the admin list and the interview (chestionar).
func (s *QuestionService) UpdateQuestionOrder(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE questions SET position = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, id := range ids {
		if _, err := stmt.Exec(i+1, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrderedQuestions returns `count` active questions in the order set in the
// admin (by position) — one after another, never random. When count <= 0 it
// returns every active question (no LIMIT).
func (s *QuestionService) GetOrderedQuestions(count int) ([]models.Question, error) {
	base := `SELECT id, question_text, question_type, difficulty_level,
	          expected_duration, position, created_at, updated_at, is_active
	          FROM questions WHERE is_active = TRUE ORDER BY position ASC, id ASC`

	return s.fetchQuestions(base, count)
}

// fetchQuestions runs the given base query, applying a LIMIT when count > 0.
func (s *QuestionService) fetchQuestions(base string, count int) ([]models.Question, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if count > 0 {
		rows, err = s.db.Query(base+" LIMIT ?", count)
	} else {
		rows, err = s.db.Query(base)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType, &q.DifficultyLevel,
			&q.ExpectedDuration, &q.Position, &q.CreatedAt, &q.UpdatedAt, &q.IsActive)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}
