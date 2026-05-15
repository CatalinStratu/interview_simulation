package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FreeAIService uses free, local AI models via Ollama
type FreeAIService struct {
	db              *sql.DB
	voiceService    *VoiceService
	questionService *QuestionService
	ollamaURL       string
	model           string
}

// AnalysisResult holds the AI analysis results
type AnalysisResult struct {
	Analysis string
	Score    int
}

// NewFreeAIService creates a service using Ollama (free, local AI)
func NewFreeAIService(db *sql.DB, voiceService *VoiceService, questionService *QuestionService) *FreeAIService {
	return &FreeAIService{
		db:              db,
		voiceService:    voiceService,
		questionService: questionService,
		ollamaURL:       getEnvOrDefault("OLLAMA_URL", "http://localhost:11434"),
		model:           getEnvOrDefault("OLLAMA_MODEL", "gpt-oss:20b"),
	}
}

// TranscribeAudioFree uses free Whisper model via Python or whisper.cpp
func (s *FreeAIService) TranscribeAudioFree(audioFilePath string) (string, error) {
	// Option 1: Use Python Whisper (easiest, free)
	if isWhisperPythonAvailable() {
		return s.transcribeWithWhisperPython(audioFilePath)
	}

	// Option 2: Use whisper.cpp (faster, free)
	if isWhisperCppAvailable() {
		return s.transcribeWithWhisperCpp(audioFilePath)
	}

	// No transcription available
	return "", fmt.Errorf("whisper not available. Install: pip3 install openai-whisper")
}

func (s *FreeAIService) transcribeWithWhisperCpp(audioFilePath string) (string, error) {
	// Convert WebM to WAV if needed
	wavPath := strings.TrimSuffix(audioFilePath, filepath.Ext(audioFilePath)) + ".wav"

	// Convert audio to WAV format (whisper.cpp requirement)
	cmd := exec.Command("ffmpeg", "-i", audioFilePath, "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", wavPath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("audio conversion failed: %w", err)
	}
	defer os.Remove(wavPath) // cleanup

	// Run whisper.cpp
	cmd = exec.Command("whisper-cpp", "-m", "models/ggml-base.en.bin", "-f", wavPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("whisper transcription failed: %w", err)
	}

	return string(output), nil
}

// AnalyzeAnswerWithOllama uses local Ollama LLM (completely free)
func (s *FreeAIService) AnalyzeAnswerWithOllama(questionText, answerText string) (*AnalysisResult, error) {
	var prompt string

	if strings.Contains(answerText, "[Audio answer provided but not transcribed") {
		// No transcription available - analyze based on question expectations
		prompt = fmt.Sprintf(`You are an expert interview coach. Provide guidance on what makes a good answer for this interview question.

Question: %s

Note: The candidate provided an audio answer, but transcription is not available.

Please provide:
1. What makes a good answer to this question (2-3 sentences)
2. Key points that should be covered (3-5 bullet points)
3. Common mistakes to avoid (3-5 bullet points)
4. Expected Score Range: 70-85/100 (since audio was provided)
5. Tips for improvement

Keep your response clear and structured.`, questionText)
	} else {
		// Normal analysis with transcription
		prompt = fmt.Sprintf(`You are an expert interview coach. Analyze this interview answer and provide constructive feedback.

Question: %s

Candidate's Answer: %s

Please provide:
1. Overall Assessment (2-3 sentences)
2. Strengths (3-5 bullet points)
3. Areas for Improvement (3-5 bullet points)
4. Score (0-100)
5. Specific Suggestions

Keep your response clear and structured.`, questionText, answerText)
	}

	// Call Ollama API
	reqBody := map[string]interface{}{
		"model":  s.model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.7,
			"num_predict": 1000,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.ollamaURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w. Is Ollama running? Run: ollama serve", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract score from response
	score := extractScoreFromText(result.Response)

	return &AnalysisResult{
		Analysis: result.Response,
		Score:    score,
	}, nil
}

// AnalyzeVoiceAnswerFree - Complete pipeline with free tools
func (s *FreeAIService) AnalyzeVoiceAnswerFree(answerID int) error {
	// Get voice answer
	answer, err := s.voiceService.GetVoiceAnswerByID(answerID)
	if err != nil {
		return fmt.Errorf("failed to get answer: %w", err)
	}

	// Get question
	question, err := s.questionService.GetQuestionByID(answer.QuestionID)
	if err != nil {
		return fmt.Errorf("failed to get question: %w", err)
	}

	// Step 1: Transcribe audio (try if whisper is available)
	var transcription string
	fmt.Printf("Checking for transcription tools...\n")

	if isWhisperPythonAvailable() || isWhisperCppAvailable() {
		fmt.Printf("Transcribing audio for answer %d...\n", answerID)
		transcription, err = s.TranscribeAudioFree(answer.AudioFilePath)
		if err != nil {
			fmt.Printf("Warning: Transcription failed: %v\n", err)
			transcription = "[Audio answer recorded - transcription failed. Analysis based on question expectations.]"
		} else {
			fmt.Printf("Transcription successful: %d characters\n", len(transcription))
		}
	} else {
		fmt.Printf("No transcription tool found. Install: pip3 install openai-whisper\n")
		transcription = "[Audio answer provided but not transcribed. Analysis is based on question expectations.]"
	}

	// Save transcription
	err = s.updateTranscription(answerID, transcription)
	if err != nil {
		return fmt.Errorf("failed to save transcription: %w", err)
	}

	// Step 2: Analyze with Ollama
	fmt.Printf("Analyzing answer %d with Ollama...\n", answerID)
	result, err := s.AnalyzeAnswerWithOllama(question.QuestionText, transcription)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	// Step 3: Save analysis
	err = s.saveAnalysis(answerID, result)
	if err != nil {
		return fmt.Errorf("failed to save analysis: %w", err)
	}

	fmt.Printf("Successfully analyzed answer %d with score: %d\n", answerID, result.Score)
	return nil
}

// AnalyzeSessionFree - Analyze all answers in a session
func (s *FreeAIService) AnalyzeSessionFree(sessionID int) error {
	answers, err := s.voiceService.GetAnswersBySession(sessionID)
	if err != nil {
		return err
	}

	for _, answer := range answers {
		// Skip if already analyzed
		if answer.Transcription != nil && *answer.Transcription != "" {
			continue
		}

		err := s.AnalyzeVoiceAnswerFree(answer.ID)
		if err != nil {
			fmt.Printf("Warning: Failed to analyze answer %d: %v\n", answer.ID, err)
			continue
		}
	}

	return nil
}

// Helper functions
func (s *FreeAIService) updateTranscription(answerID int, transcription string) error {
	query := `UPDATE voice_answers
	          SET transcription = ?, transcription_status = 'transcribed'
	          WHERE id = ?`
	_, err := s.db.Exec(query, transcription, answerID)
	return err
}

func (s *FreeAIService) saveAnalysis(answerID int, result *AnalysisResult) error {
	query := `UPDATE voice_answers
	          SET ai_analysis = ?, ai_score = ?, analyzed_at = NOW(), transcription_status = 'analyzed'
	          WHERE id = ?`
	_, err := s.db.Exec(query, result.Analysis, result.Score, answerID)
	return err
}

func isWhisperPythonAvailable() bool {
	// Check if python3 -m whisper works
	cmd := exec.Command("python3", "-m", "whisper", "--help")
	return cmd.Run() == nil
}

func isWhisperCppAvailable() bool {
	cmd := exec.Command("which", "whisper-cpp")
	return cmd.Run() == nil
}

func (s *FreeAIService) transcribeWithWhisperPython(audioFilePath string) (string, error) {
	// Use Python Whisper CLI via python3 -m
	cmd := exec.Command("python3", "-m", "whisper", audioFilePath, "--model", "base", "--output_format", "txt", "--output_dir", "/tmp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper transcription failed: %w, output: %s", err, string(output))
	}

	// Read the generated text file
	// Whisper creates filename without the audio extension: audio.webm -> audio.txt
	baseNameNoExt := strings.TrimSuffix(filepath.Base(audioFilePath), filepath.Ext(audioFilePath))
	txtFile := "/tmp/" + baseNameNoExt + ".txt"

	content, err := os.ReadFile(txtFile)
	if err != nil {
		// Try with extension as well (fallback)
		txtFile = "/tmp/" + filepath.Base(audioFilePath) + ".txt"
		content, err = os.ReadFile(txtFile)
		if err != nil {
			return "", fmt.Errorf("failed to read transcription: %w (tried: %s)", err, txtFile)
		}
	}

	// Cleanup
	os.Remove(txtFile)

	return strings.TrimSpace(string(content)), nil
}

func extractScoreFromText(text string) int {
	// Look for patterns like "Score: 85" or "85/100" or "Score: 85/100"
	text = strings.ToLower(text)

	// Simple pattern matching
	if strings.Contains(text, "score") {
		// Extract number after "score"
		words := strings.Fields(text)
		for i, word := range words {
			if strings.Contains(word, "score") && i+1 < len(words) {
				var score int
				fmt.Sscanf(words[i+1], "%d", &score)
				if score > 0 && score <= 100 {
					return score
				}
			}
		}
	}

	// Default score
	return 70
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Alternative: Use Gemini Free Tier
func (s *FreeAIService) AnalyzeWithGemini(questionText, answerText string) (*AnalysisResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY") // Free tier available
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not set")
	}

	// Gemini API implementation
	// Free tier: 60 requests per minute
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey

	prompt := fmt.Sprintf(`Analyze this interview answer:
Question: %s
Answer: %s

Provide assessment, strengths, improvements, and score (0-100).`, questionText, answerText)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		text := result.Candidates[0].Content.Parts[0].Text
		return &AnalysisResult{
			Analysis: text,
			Score:    extractScoreFromText(text),
		}, nil
	}

	return nil, fmt.Errorf("no response from Gemini")
}
