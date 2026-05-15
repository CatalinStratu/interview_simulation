# Changelog

## [1.1.0] - AI Analysis Feature

### Added
- **AI-Powered Answer Analysis**
  - OpenAI Whisper integration for speech-to-text
  - GPT-4 integration for answer analysis
  - Automatic transcription of voice answers
  - Detailed feedback with scoring (0-100)
  - Batch analysis for entire sessions

- **New API Endpoints**
  - `POST /api/ai/analyze/answer/{id}` - Analyze single answer
  - `POST /api/ai/analyze/session/{id}` - Analyze entire session

- **Database Enhancements**
  - `transcription_status` column in voice_answers
  - `ai_analysis` column for storing feedback
  - `ai_score` column for numerical scores
  - `analyzed_at` timestamp

- **UI Improvements**
  - "🤖 AI Analyze" button in admin dashboard
  - Display transcriptions in session details
  - Show AI analysis and scores
  - Visual feedback during analysis

- **Documentation**
  - `AI_ANALYSIS.md` - Complete AI feature guide
  - Setup instructions for OpenAI API
  - Cost estimation and optimization tips

### Changed
- Updated `go.mod` with OpenAI client dependency
- Enhanced admin dashboard to display AI results
- Improved session detail view with analysis

### Technical
- New service: `ai_service.go`
- New handler: `ai_handler.go`
- New schema: `ai_schema.sql`
- Environment variable: `OPENAI_API_KEY`

---

## [1.0.0] - Initial Release

### Added
- Voice-enabled interview chatbot
- Text-to-speech question delivery
- Voice answer recording
- MySQL database with complete schema
- Admin dashboard for question management
- Session management and tracking
- Remarks and feedback system
- RESTful API with Gorilla Mux
- Responsive web interface
- 10 sample interview questions

### Features
- **Interview Interface**
  - Start new sessions
  - Listen to TTS questions
  - Record voice answers
  - Progress tracking
  - Session completion

- **Admin Dashboard**
  - Add/delete questions
  - View all sessions
  - Add remarks to sessions
  - Session detail view

- **Database**
  - Questions table
  - Interview sessions
  - Voice answers storage
  - Remarks system
  - Session questions tracking

### Documentation
- `README.md` - Project overview
- `SETUP.md` - Detailed setup guide
- `QUICKSTART.md` - 5-minute quick start
- `Makefile` - Convenience commands
