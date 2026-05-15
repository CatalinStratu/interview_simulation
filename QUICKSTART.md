# Quick Start Guide

## Prerequisites Check

1. **Install Go** (if not installed):
   ```bash
   # macOS
   brew install go

   # Verify
   go version
   ```

2. **Install MySQL** (if not installed):
   ```bash
   # macOS
   brew install mysql
   brew services start mysql
   ```

## 5-Minute Setup

### Step 1: Create Database
```bash
mysql -u root -p
```

In MySQL prompt:
```sql
CREATE DATABASE interview_chatbot;
EXIT;
```

### Step 2: Configure Environment
```bash
cp .env.example .env
```

Edit `.env` and set your MySQL password:
```
DB_PASSWORD=your_mysql_password
```

### Step 3: Install Dependencies
```bash
go mod download
```

### Step 4: Run Application
```bash
go run cmd/server/main.go
```

You should see:
```
Database connection established
Database migrations completed successfully
Server starting on http://localhost:8080
Interview interface: http://localhost:8080/
Admin dashboard: http://localhost:8080/admin
```

### Step 5: Test the Application

1. **Open interview interface**: http://localhost:8080/
2. **Start an interview** with 3-5 questions
3. **Record voice answers** (allow microphone access)
4. **View admin dashboard**: http://localhost:8080/admin
5. **Add remarks** to the completed session

## Using the Makefile

For convenience, use the included Makefile:

```bash
# Install dependencies
make install

# Run the application
make run

# Build binary
make build

# Create database
make db-create
```

## Architecture Overview

### Backend (Go)
- **Router**: Gorilla Mux
- **Database**: MySQL with prepared statements
- **Services**: Question, Session, Voice, Remark, TTS
- **Handlers**: RESTful API endpoints

### Frontend (Vanilla JS)
- **Interview Interface**: Voice recording and playback
- **Admin Dashboard**: Question management and session review
- **TTS**: Text-to-speech for questions
- **MediaRecorder API**: Voice recording

### Database Schema
- `questions` - Interview questions library
- `interview_sessions` - Active and completed sessions
- `voice_answers` - Recorded audio responses
- `remarks` - Admin feedback and ratings
- `session_questions` - Junction table for session questions

## Key Features

- Voice-enabled interview simulation
- Text-to-speech question delivery
- Voice answer recording and storage
- Admin dashboard for question management
- Session review and feedback system
- No authentication (simplified for development)

## Next Steps

### For Development
1. Review code in `/internal` directory
2. Customize questions in database or admin UI
3. Test with different question types
4. Experiment with TTS features

### For Production
1. Add authentication (JWT, OAuth)
2. Implement proper TTS service (Google Cloud, AWS Polly)
3. Add cloud storage for audio files (S3, GCS)
4. Enable HTTPS (required for microphone access)
5. Add logging and monitoring
6. Implement rate limiting
7. Add input validation and sanitization

## Troubleshooting

**Go not found?**
- Install Go from https://golang.org/dl/
- Add to PATH: `export PATH=$PATH:/usr/local/go/bin`

**MySQL connection failed?**
- Check MySQL is running: `brew services list`
- Verify credentials in `.env`
- Ensure database exists: `mysql -u root -p -e "SHOW DATABASES;"`

**Port 8080 in use?**
- Change port in `.env`: `SERVER_PORT=8081`

**Microphone not working?**
- Allow microphone access in browser
- Use Chrome/Firefox (best support)
- HTTPS required in production

## Project Structure

```
interview-chatbot/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── internal/
│   ├── database/        # Database layer
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── web/
│   ├── static/          # CSS, JavaScript
│   └── templates/       # HTML templates
├── go.mod               # Dependencies
├── Makefile            # Build commands
└── README.md           # Documentation
```

## Support

For detailed documentation, see:
- `README.md` - Full project documentation
- `SETUP.md` - Detailed setup instructions
- Code comments in source files

Enjoy your interview chatbot!
