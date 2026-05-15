package handlers

import (
	"encoding/json"
	"interview-chatbot/internal/models"
	"interview-chatbot/internal/services"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	sessionService *services.SessionService
	voiceService   *services.VoiceService
	ttsService     *services.TTSService
	remarkService  *services.RemarkService
}

func NewSessionHandler(
	sessionService *services.SessionService,
	voiceService *services.VoiceService,
	ttsService *services.TTSService,
	remarkService *services.RemarkService,
) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		voiceService:   voiceService,
		ttsService:     ttsService,
		remarkService:  remarkService,
	}
}

func (h *SessionHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	var req models.StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// NumQuestions <= 0 means "all active questions" — handled downstream.

	session, err := h.sessionService.StartSession(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start session")
		return
	}

	respondWithJSON(w, http.StatusCreated, session)
}

func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := h.sessionService.GetSessionByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Session not found")
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

func (h *SessionHandler) GetNextQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	response, err := h.sessionService.GetNextQuestion(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get next question")
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *SessionHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	questionIDStr := r.FormValue("question_id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	// Get the audio file
	file, header, err := r.FormFile("audio")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Audio file is required")
		return
	}
	defer file.Close()

	// Save voice answer
	answer, err := h.voiceService.SaveVoiceAnswer(sessionID, questionID, file, header.Filename)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save answer")
		return
	}

	respondWithJSON(w, http.StatusCreated, answer)
}

func (h *SessionHandler) CompleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	if err := h.sessionService.CompleteSession(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to complete session")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Session completed successfully"})
}

func (h *SessionHandler) GetAllSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.sessionService.GetAllSessions()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch sessions")
		return
	}

	respondWithJSON(w, http.StatusOK, sessions)
}

func (h *SessionHandler) GetSessionAnswers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	answers, err := h.voiceService.GetAnswersBySession(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch answers")
		return
	}

	respondWithJSON(w, http.StatusOK, answers)
}

func (h *SessionHandler) AddRemark(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	var req models.AddRemarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.RemarkText == "" {
		respondWithError(w, http.StatusBadRequest, "Remark text is required")
		return
	}

	remark, err := h.remarkService.AddRemark(sessionID, req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to add remark")
		return
	}

	respondWithJSON(w, http.StatusCreated, remark)
}

func (h *SessionHandler) GetSessionRemarks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	remarks, err := h.remarkService.GetRemarksBySession(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch remarks")
		return
	}

	respondWithJSON(w, http.StatusOK, remarks)
}

// GetAnswerAudio serves the recorded answer audio file for a given voice_answer ID.
// Used by the admin UI to play back candidate answers.
func (h *SessionHandler) GetAnswerAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	answerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid answer ID")
		return
	}

	answer, err := h.voiceService.GetVoiceAnswerByID(answerID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Answer not found")
		return
	}

	// Force audio mime types so browsers play these in an <audio> element
	// (http.ServeFile would otherwise tag .webm as video/webm).
	switch strings.ToLower(filepath.Ext(answer.AudioFilePath)) {
	case ".webm":
		w.Header().Set("Content-Type", "audio/webm")
	case ".ogg", ".oga":
		w.Header().Set("Content-Type", "audio/ogg")
	case ".wav":
		w.Header().Set("Content-Type", "audio/wav")
	case ".mp3":
		w.Header().Set("Content-Type", "audio/mpeg")
	case ".m4a", ".mp4":
		w.Header().Set("Content-Type", "audio/mp4")
	}

	http.ServeFile(w, r, answer.AudioFilePath)
}

func (h *SessionHandler) GetQuestionAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid question ID")
		return
	}

	questionText := r.URL.Query().Get("text")
	if questionText == "" {
		respondWithError(w, http.StatusBadRequest, "Question text is required")
		return
	}

	// Generate or retrieve cached audio
	audioPath, err := h.ttsService.GenerateSpeech(questionText, questionID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate audio")
		return
	}

	// Serve the audio file
	http.ServeFile(w, r, audioPath)
}
