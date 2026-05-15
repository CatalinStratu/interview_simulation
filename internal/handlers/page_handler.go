package handlers

import (
	"html/template"
	"net/http"
)

type PageHandler struct {
	templates *template.Template
}

func NewPageHandler() *PageHandler {
	// Load templates
	templates := template.Must(template.ParseGlob("web/templates/*.html"))
	return &PageHandler{templates: templates}
}

func (h *PageHandler) InterviewPage(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "interview.html", nil)
}

func (h *PageHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "admin.html", nil)
}

func (h *PageHandler) VoiceOnlyPage(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "interview_voice.html", nil)
}
