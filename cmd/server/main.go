package main

import (
	"interview-chatbot/config"
	"interview-chatbot/internal/database"
	"interview-chatbot/internal/handlers"
	"interview-chatbot/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.InitDB(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed questions (idempotent — skips when already applied)
	if err := db.RunSeed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize services
	questionService := services.NewQuestionService(db.DB)
	sessionService := services.NewSessionService(db.DB, questionService)
	voiceService := services.NewVoiceService(db.DB, "voice_answers", sessionService)
	ttsService := services.NewTTSService("tts_cache")
	remarkService := services.NewRemarkService(db.DB)

	// Initialize AI service (using free Llama via Ollama)
	log.Println("Using FREE AI (Ollama/Llama) for analysis")
	aiService := services.NewFreeAIService(db.DB, voiceService, questionService)

	// Initialize handlers
	questionHandler := handlers.NewQuestionHandler(questionService)
	sessionHandler := handlers.NewSessionHandler(sessionService, voiceService, ttsService, remarkService)
	aiHandler := handlers.NewAIHandler(aiService)
	pageHandler := handlers.NewPageHandler()

	// Setup router
	r := mux.NewRouter()

	// Page routes
	r.HandleFunc("/", pageHandler.InterviewPage).Methods("GET")
	r.HandleFunc("/admin", pageHandler.AdminPage).Methods("GET")
	r.HandleFunc("/voice", pageHandler.VoiceOnlyPage).Methods("GET")

	// API routes - Interview Session
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/sessions/start", sessionHandler.StartSession).Methods("POST")
	api.HandleFunc("/sessions/{id}", sessionHandler.GetSession).Methods("GET")
	api.HandleFunc("/sessions/{id}/next-question", sessionHandler.GetNextQuestion).Methods("GET")
	api.HandleFunc("/sessions/{id}/answer", sessionHandler.SubmitAnswer).Methods("POST")
	api.HandleFunc("/sessions/{id}/complete", sessionHandler.CompleteSession).Methods("POST")
	api.HandleFunc("/sessions/{id}/answers", sessionHandler.GetSessionAnswers).Methods("GET")

	// API routes - Admin Question Management
	api.HandleFunc("/admin/questions", questionHandler.GetAllQuestions).Methods("GET")
	api.HandleFunc("/admin/questions", questionHandler.CreateQuestion).Methods("POST")
	api.HandleFunc("/admin/questions/{id}", questionHandler.GetQuestion).Methods("GET")
	api.HandleFunc("/admin/questions/{id}", questionHandler.DeleteQuestion).Methods("DELETE")

	// API routes - Admin Session Management
	api.HandleFunc("/admin/sessions", sessionHandler.GetAllSessions).Methods("GET")
	api.HandleFunc("/admin/sessions/{id}", sessionHandler.GetSession).Methods("GET")
	api.HandleFunc("/admin/sessions/{id}/remarks", sessionHandler.AddRemark).Methods("POST")
	api.HandleFunc("/admin/sessions/{id}/remarks", sessionHandler.GetSessionRemarks).Methods("GET")

	// TTS route
	api.HandleFunc("/tts/{id}", sessionHandler.GetQuestionAudio).Methods("GET")

	// Recorded answer playback (admin)
	api.HandleFunc("/admin/answers/{id}/audio", sessionHandler.GetAnswerAudio).Methods("GET")

	// AI Analysis routes
	api.HandleFunc("/ai/analyze/answer/{id}", aiHandler.AnalyzeAnswer).Methods("POST")
	api.HandleFunc("/ai/analyze/session/{id}", aiHandler.AnalyzeSession).Methods("POST")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// CORS middleware
	r.Use(corsMiddleware)

	// Start server
	addr := ":" + cfg.Server.Port
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Interview interface: http://localhost%s/", addr)
	log.Printf("Voice-only interview: http://localhost%s/voice", addr)
	log.Printf("Admin dashboard: http://localhost%s/admin", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
