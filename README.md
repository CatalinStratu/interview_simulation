# Interview Simulation Chatbot

A voice-enabled chatbot application for conducting interview simulations.

## Features

- Voice-based interview sessions with text-to-speech questions
- Voice answer recording and storage
- Admin dashboard for managing interview questions
- Remarks/feedback system for each interview session
- MySQL database for persistent storage

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher

## Installation

1. Install Go from https://golang.org/dl/

2. Install MySQL and create a database:
```sql
CREATE DATABASE interview_chatbot;
```

3. Install dependencies:
```bash
go mod download
```

4. Configure database connection in `config/config.go`

5. Run the application:
```bash
go run cmd/server/main.go
```

## Usage

### Admin Dashboard
- Navigate to `http://localhost:8080/admin`
- Add interview questions
- View all interview sessions
- Add remarks to sessions

### Interview Interface
- Navigate to `http://localhost:8080/`
- Start a new interview session
- Listen to questions (text-to-speech)
- Record voice answers
- Submit interview

## Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── database/        # Database initialization and migrations
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── web/
│   ├── static/          # CSS, JS files
│   └── templates/       # HTML templates
└── config/              # Configuration
```

## API Endpoints

### Interview Endpoints
- `GET /` - Interview interface
- `POST /api/sessions/start` - Start new interview session
- `GET /api/sessions/:id/next-question` - Get next question
- `POST /api/sessions/:id/answer` - Submit voice answer
- `POST /api/sessions/:id/complete` - Complete interview session

### Admin Endpoints
- `GET /admin` - Admin dashboard
- `GET /api/admin/questions` - List all questions
- `POST /api/admin/questions` - Add new question
- `DELETE /api/admin/questions/:id` - Delete question
- `GET /api/admin/sessions` - List all sessions
- `POST /api/admin/sessions/:id/remarks` - Add remarks to session

### TTS Endpoint
- `GET /api/tts/:question_id` - Get audio for question

## Database Schema

See `internal/database/schema.sql` for complete schema.

## License

MIT
