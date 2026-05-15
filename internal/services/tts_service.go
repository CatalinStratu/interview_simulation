package services

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TTSService handles text-to-speech conversion.
//
// For Romanian, it calls a local Facebook MMS-TTS sidecar
// (see tts_mms/server.py). If the sidecar is unreachable it transparently
// falls back to Google Translate TTS so questions still play.
type TTSService struct {
	cacheDir string
	mmsURL   string
	mmsHTTP  *http.Client
}

func NewTTSService(cacheDir string) *TTSService {
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create TTS cache directory: %v", err))
	}
	mmsURL := os.Getenv("MMS_TTS_URL")
	if mmsURL == "" {
		mmsURL = "http://127.0.0.1:5005/tts"
	}
	return &TTSService{
		cacheDir: cacheDir,
		mmsURL:   mmsURL,
		mmsHTTP:  &http.Client{Timeout: 60 * time.Second},
	}
}

// GenerateSpeech converts text to speech and returns the audio file path.
// Defaults to Romanian (routed through the MMS sidecar).
func (s *TTSService) GenerateSpeech(text string, questionID int) (string, error) {
	return s.GenerateSpeechWithLang(text, questionID, "ro")
}

// GenerateSpeechWithLang converts text to speech with the specified language.
func (s *TTSService) GenerateSpeechWithLang(text string, questionID int, lang string) (string, error) {
	if strings.EqualFold(lang, "ro") {
		if path, err := s.generateMMS(text, questionID, lang); err == nil {
			return path, nil
		} else {
			log.Printf("TTS: MMS sidecar failed (%v), falling back to Google", err)
		}
	}
	return s.generateGoogle(text, questionID, lang)
}

func (s *TTSService) generateMMS(text string, questionID int, lang string) (string, error) {
	hash := md5.Sum([]byte("facebook-mms|" + lang + "|" + text))
	filename := fmt.Sprintf("question_%d_facebook-mms_%s_%x.wav", questionID, lang, hash)
	filePath := filepath.Join(s.cacheDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	body, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.mmsURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "audio/wav")

	resp, err := s.mmsHTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("sidecar returned %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filePath, audio, 0644); err != nil {
		return "", err
	}
	return filePath, nil
}

func (s *TTSService) generateGoogle(text string, questionID int, lang string) (string, error) {
	hash := md5.Sum([]byte("google|" + lang + "|" + text))
	filename := fmt.Sprintf("question_%d_google_%s_%x.mp3", questionID, lang, hash)
	filePath := filepath.Join(s.cacheDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	ttsURL := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&client=tw-ob&tl=%s&q=%s",
		lang, url.QueryEscape(text))

	resp, err := http.Get(ttsURL)
	if err != nil {
		return "", fmt.Errorf("failed to download TTS audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("TTS API returned status: %d", resp.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create audio file: %w", err)
	}
	defer file.Close()

	if _, err = io.Copy(file, resp.Body); err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save audio file: %w", err)
	}
	return filePath, nil
}

// GetAudioPath returns the cached audio file path if it exists.
func (s *TTSService) GetAudioPath(questionID int, text string) (string, bool) {
	hash := md5.Sum([]byte("facebook-mms|ro|" + text))
	filename := fmt.Sprintf("question_%d_facebook-mms_ro_%x.wav", questionID, hash)
	filePath := filepath.Join(s.cacheDir, filename)
	if _, err := os.Stat(filePath); err == nil {
		return filePath, true
	}
	return "", false
}
