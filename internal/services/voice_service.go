package services

import (
	"database/sql"
	"fmt"
	"interview-chatbot/internal/models"
	"io"
	"os"
	"path/filepath"
	"time"
)

type VoiceService struct {
	db             *sql.DB
	uploadDir      string
	sessionService *SessionService
}

func NewVoiceService(db *sql.DB, uploadDir string, sessionService *SessionService) *VoiceService {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
	return &VoiceService{
		db:             db,
		uploadDir:      uploadDir,
		sessionService: sessionService,
	}
}

func (s *VoiceService) SaveVoiceAnswer(sessionID, questionID int, audioData io.Reader, filename string) (*models.VoiceAnswer, error) {
	// Generate unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("session_%d_question_%d_%d%s", sessionID, questionID, timestamp, ext)
	filePath := filepath.Join(s.uploadDir, newFilename)

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy audio data to file
	_, err = io.Copy(file, audioData)
	if err != nil {
		return nil, fmt.Errorf("failed to save audio data: %w", err)
	}

	// Save to database
	query := `INSERT INTO voice_answers (session_id, question_id, audio_file_path)
	          VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, sessionID, questionID, filePath)
	if err != nil {
		// Delete the file if database insert fails
		os.Remove(filePath)
		return nil, err
	}

	// Update session questions_answered count
	updateQuery := `UPDATE interview_sessions SET questions_answered = questions_answered + 1
	                WHERE id = ?`
	_, err = s.db.Exec(updateQuery, sessionID)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetVoiceAnswerByID(int(id))
}

func (s *VoiceService) GetVoiceAnswerByID(id int) (*models.VoiceAnswer, error) {
	query := `SELECT id, session_id, question_id, audio_file_path, audio_duration,
	          transcription, answered_at FROM voice_answers WHERE id = ?`

	var answer models.VoiceAnswer
	err := s.db.QueryRow(query, id).Scan(&answer.ID, &answer.SessionID, &answer.QuestionID,
		&answer.AudioFilePath, &answer.AudioDuration, &answer.Transcription, &answer.AnsweredAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("voice answer not found")
	}
	if err != nil {
		return nil, err
	}

	return &answer, nil
}

func (s *VoiceService) GetAnswersBySession(sessionID int) ([]models.VoiceAnswer, error) {
	query := `SELECT id, session_id, question_id, audio_file_path, audio_duration,
	          transcription, answered_at FROM voice_answers WHERE session_id = ?
	          ORDER BY answered_at`

	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.VoiceAnswer
	for rows.Next() {
		var answer models.VoiceAnswer
		err := rows.Scan(&answer.ID, &answer.SessionID, &answer.QuestionID,
			&answer.AudioFilePath, &answer.AudioDuration, &answer.Transcription, &answer.AnsweredAt)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
