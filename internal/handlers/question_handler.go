package handlers

import (
	"encoding/json"
	"interview-chatbot/internal/models"
	"interview-chatbot/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type QuestionHandler struct {
	questionService *services.QuestionService
}

func NewQuestionHandler(questionService *services.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

func (h *QuestionHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questionService.GetAllQuestions()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch questions")
		return
	}

	respondWithJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	question, err := h.questionService.GetQuestionByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Question not found")
		return
	}

	respondWithJSON(w, http.StatusOK, question)
}

func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var req models.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.QuestionText == "" {
		respondWithError(w, http.StatusBadRequest, "Question text is required")
		return
	}

	if req.QuestionType == "" {
		req.QuestionType = "technical"
	}
	if req.DifficultyLevel == "" {
		req.DifficultyLevel = "medium"
	}
	if req.ExpectedDuration == 0 {
		req.ExpectedDuration = 120
	}

	question, err := h.questionService.CreateQuestion(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create question")
		return
	}

	respondWithJSON(w, http.StatusCreated, question)
}

func (h *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	if err := h.questionService.DeleteQuestion(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete question")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Question deleted successfully"})
}
