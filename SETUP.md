# Setup Guide - Interview Chatbot

## Prerequisites

### 1. Install Go
Download and install Go 1.21 or higher from [https://golang.org/dl/](https://golang.org/dl/)

For macOS:
```bash
brew install go
```

Verify installation:
```bash
go version
```

### 2. Install MySQL
Download and install MySQL 8.0 or higher from [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)

For macOS:
```bash
brew install mysql
brew services start mysql
```

## Database Setup

### 1. Create Database
```bash
mysql -u root -p
```

In MySQL shell:
```sql
CREATE DATABASE interview_chatbot;
EXIT;
```

### 2. Configure Database Credentials
Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

Edit `.env` with your MySQL credentials:
```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=interview_chatbot
SERVER_PORT=8080
```

## Application Setup

### 1. Install Dependencies
```bash
go mod download
```

### 2. Run the Application
```bash
go run cmd/server/main.go
```

The application will:
- Connect to MySQL database
- Run migrations automatically (create tables and insert sample questions)
- Start the web server on port 8080

### 3. Access the Application

**Interview Interface:**
[http://localhost:8080/](http://localhost:8080/)

**Admin Dashboard:**
[http://localhost:8080/admin](http://localhost:8080/admin)

## Usage Guide

### For Candidates (Interview Interface)

1. Navigate to [http://localhost:8080/](http://localhost:8080/)
2. Enter your name and email (optional)
3. Select number of questions
4. Click "Start Interview"
5. For each question:
   - Click "Play Question" to hear the question via text-to-speech
   - Click "Start Recording" to record your answer
   - Click "Stop Recording" when done
   - Click "Submit Answer" to move to next question
6. Complete all questions
7. View your session code for reference

### For Admins (Admin Dashboard)

1. Navigate to [http://localhost:8080/admin](http://localhost:8080/admin)

**Manage Questions:**
- View all active questions
- Add new questions with type, difficulty, and expected duration
- Delete questions

**View Sessions:**
- See all interview sessions
- View session details including:
  - Candidate information
  - Progress and completion status
  - Voice answers
  - Add remarks and ratings

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── config/
│   └── config.go                # Configuration management
├── internal/
│   ├── database/
│   │   ├── database.go          # Database connection
│   │   └── schema.sql           # Database schema and migrations
│   ├── handlers/
│   │   ├── question_handler.go  # Question API handlers
│   │   ├── session_handler.go   # Session API handlers
│   │   ├── page_handler.go      # Page rendering handlers
│   │   └── utils.go             # Handler utilities
│   ├── models/
│   │   └── models.go            # Data models and DTOs
│   └── services/
│       ├── question_service.go  # Question business logic
│       ├── session_service.go   # Session business logic
│       ├── voice_service.go     # Voice recording management
│       ├── remark_service.go    # Remarks management
│       └── tts_service.go       # Text-to-speech service
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css        # Styles
│   │   └── js/
│   │       ├── interview.js     # Interview interface logic
│   │       └── admin.js         # Admin dashboard logic
│   └── templates/
│       ├── interview.html       # Interview page template
│       └── admin.html           # Admin page template
├── go.mod                       # Go module definition
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
└── README.md                    # Project documentation
```

## Features

### Voice Recording
- Uses browser's MediaRecorder API
- Records in WebM format
- Allows playback before submission
- Stores recordings on server

### Text-to-Speech
- Converts question text to speech
- Uses Google Translate TTS (free tier)
- Caches audio files for performance
- For production, consider:
  - Google Cloud Text-to-Speech
  - AWS Polly
  - Azure Speech Services

### Database Schema
- **questions**: Interview questions with metadata
- **interview_sessions**: Session tracking
- **voice_answers**: Audio recordings storage
- **remarks**: Admin feedback on sessions
- **session_questions**: Questions asked per session

## API Endpoints

### Interview Endpoints
```
POST   /api/sessions/start                  Start new session
GET    /api/sessions/:id                    Get session details
GET    /api/sessions/:id/next-question      Get next question
POST   /api/sessions/:id/answer             Submit voice answer
POST   /api/sessions/:id/complete           Complete session
GET    /api/sessions/:id/answers            Get all answers
```

### Admin Endpoints
```
GET    /api/admin/questions                 List questions
POST   /api/admin/questions                 Create question
DELETE /api/admin/questions/:id             Delete question
GET    /api/admin/sessions                  List sessions
GET    /api/admin/sessions/:id              Get session details
POST   /api/admin/sessions/:id/remarks      Add remark
GET    /api/admin/sessions/:id/remarks      Get remarks
```

### TTS Endpoint
```
GET    /api/tts/:id?text=...                Get question audio
```

## Troubleshooting

### MySQL Connection Error
- Ensure MySQL is running: `brew services list`
- Check credentials in `.env` file
- Verify database exists: `SHOW DATABASES;`

### Port Already in Use
Change the port in `.env`:
```
SERVER_PORT=8081
```

### Microphone Access Denied
- Allow microphone access in browser settings
- Use HTTPS in production (required for MediaRecorder API)

### TTS Not Working
- Check internet connection (Google Translate TTS requires internet)
- For production, implement proper TTS service with API keys

## Production Considerations

1. **Security:**
   - Add authentication/authorization
   - Implement HTTPS
   - Validate and sanitize all inputs
   - Use prepared statements (already implemented)

2. **TTS Service:**
   - Replace Google Translate TTS with proper service
   - Add API key management
   - Implement rate limiting

3. **File Storage:**
   - Consider cloud storage (S3, Google Cloud Storage)
   - Implement file size limits
   - Add virus scanning

4. **Performance:**
   - Add database indexes
   - Implement caching
   - Use connection pooling

5. **Monitoring:**
   - Add logging
   - Implement error tracking
   - Add metrics and monitoring

## License

MIT
