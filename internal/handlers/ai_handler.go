package handlers

import (
	"interview-chatbot/internal/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AIHandler struct {
	freeAIService *services.FreeAIService
}

func NewAIHandler(freeAIService *services.FreeAIService) *AIHandler {
	return &AIHandler{
		freeAIService: freeAIService,
	}
}

// AnalyzeAnswer - Analyze a single voice answer
func (h *AIHandler) AnalyzeAnswer(w http.ResponseWriter, r *http.Request) {
	if h.freeAIService == nil {
		respondWithError(w, http.StatusServiceUnavailable, "AI service not configured")
		return
	}

	vars := mux.Vars(r)
	answerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	// Run analysis with Llama
	err = h.freeAIService.AnalyzeVoiceAnswerFree(answerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to analyze answer: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Answer analyzed successfully",
	})
}

// AnalyzeSession - Analyze all answers in a session
func (h *AIHandler) AnalyzeSession(w http.ResponseWriter, r *http.Request) {
	if h.freeAIService == nil {
		respondWithError(w, http.StatusServiceUnavailable, "AI service not configured")
		return
	}

	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	// Run analysis with Llama
	err = h.freeAIService.AnalyzeSessionFree(sessionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to analyze session: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Session analyzed successfully",
	})
}
