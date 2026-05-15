package services

import (
	"database/sql"
	"fmt"
	"interview-chatbot/internal/models"
)

type RemarkService struct {
	db *sql.DB
}

func NewRemarkService(db *sql.DB) *RemarkService {
	return &RemarkService{db: db}
}

func (s *RemarkService) AddRemark(sessionID int, req models.AddRemarkRequest) (*models.Remark, error) {
	query := `INSERT INTO remarks (session_id, remark_text, rating)
	          VALUES (?, ?, ?)`

	result, err := s.db.Exec(query, sessionID, req.RemarkText, req.Rating)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return s.GetRemarkByID(int(id))
}

func (s *RemarkService) GetRemarkByID(id int) (*models.Remark, error) {
	query := `SELECT id, session_id, remark_text, rating, created_by, created_at, updated_at
	          FROM remarks WHERE id = ?`

	var remark models.Remark
	err := s.db.QueryRow(query, id).Scan(&remark.ID, &remark.SessionID, &remark.RemarkText,
		&remark.Rating, &remark.CreatedBy, &remark.CreatedAt, &remark.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("remark not found")
	}
	if err != nil {
		return nil, err
	}

	return &remark, nil
}

func (s *RemarkService) GetRemarksBySession(sessionID int) ([]models.Remark, error) {
	query := `SELECT id, session_id, remark_text, rating, created_by, created_at, updated_at
	          FROM remarks WHERE session_id = ? ORDER BY created_at DESC`

	rows, err := s.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var remarks []models.Remark
	for rows.Next() {
		var remark models.Remark
		err := rows.Scan(&remark.ID, &remark.SessionID, &remark.RemarkText,
			&remark.Rating, &remark.CreatedBy, &remark.CreatedAt, &remark.UpdatedAt)
		if err != nil {
			return nil, err
		}
		remarks = append(remarks, remark)
	}

	return remarks, nil
}

func (s *RemarkService) UpdateRemark(id int, req models.AddRemarkRequest) (*models.Remark, error) {
	query := `UPDATE remarks SET remark_text = ?, rating = ?, updated_at = NOW()
	          WHERE id = ?`

	_, err := s.db.Exec(query, req.RemarkText, req.Rating, id)
	if err != nil {
		return nil, err
	}

	return s.GetRemarkByID(id)
}

func (s *RemarkService) DeleteRemark(id int) error {
	query := `DELETE FROM remarks WHERE id = ?`
	_, err := s.db.Exec(query, id)
	return err
}
