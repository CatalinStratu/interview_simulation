# Free AI Analysis Setup

Use **completely free** AI models instead of paid OpenAI API. No costs, runs locally!

## Option 1: Ollama (Recommended - 100% Free)

### What is Ollama?
- Run AI models **locally** on your Mac
- **Completely free**, no API keys
- **Private** - data never leaves your computer
- Fast and easy to use

### Installation

**1. Install Ollama:**
```bash
# macOS
curl -fsSL https://ollama.com/install.sh | sh

# Or download from: https://ollama.com/download
```

**2. Start Ollama:**
```bash
ollama serve
```

**3. Download AI Model:**
```bash
# Llama 3.2 (recommended, ~2GB)
ollama pull llama3.2

# Or other free models:
ollama pull mistral        # Fast and good
ollama pull gemma2         # Google's model
ollama pull phi3           # Small and fast
```

**4. Test It:**
```bash
ollama run llama3.2 "Analyze this interview answer: I have 5 years of Python experience..."
```

### Configure Application

**Update .env:**
```bash
# Use free local AI instead of OpenAI
USE_FREE_AI=true
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2
```

**Update go.mod:**
```bash
go mod tidy
```

**Restart App:**
```bash
go run cmd/server/main.go
```

### Usage
Same as before! Just click "🤖 AI Analyze" in admin dashboard. Now it uses Ollama (free) instead of OpenAI (paid).

---

## Option 2: Google Gemini (Free Tier)

Google offers **free tier** with generous limits.

### Setup Gemini

**1. Get Free API Key:**
- Visit: https://makersuite.google.com/app/apikey
- Create free account
- Generate API key (no credit card required)

**2. Configure .env:**
```bash
USE_GEMINI=true
GEMINI_API_KEY=your_free_api_key_here
```

**3. Free Tier Limits:**
- 60 requests per minute
- 1,500 requests per day
- Completely free for personal use

---

## Option 3: Whisper.cpp (Free Speech-to-Text)

For converting voice to text locally.

### Install Whisper.cpp

```bash
# Install dependencies
brew install ffmpeg

# Clone whisper.cpp
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp

# Build
make

# Download model (base is ~150MB)
bash ./models/download-ggml-model.sh base.en

# Add to PATH
sudo cp whisper-cpp /usr/local/bin/
```

### Test It:
```bash
whisper-cpp -m models/ggml-base.en.bin -f your_audio.wav
```

---

## Comparison

| Feature | Ollama | Gemini Free | OpenAI |
|---------|--------|-------------|---------|
| **Cost** | 100% Free | Free tier | ~$0.20/session |
| **Setup** | 5 minutes | 2 minutes | 1 minute |
| **Privacy** | Full (local) | Cloud | Cloud |
| **Speed** | Fast (local) | Fast | Fast |
| **Quality** | Very Good | Excellent | Excellent |
| **Limits** | None | 1,500/day | Pay per use |
| **Internet** | Not needed* | Required | Required |

*After initial model download

---

## Complete Free Setup (Recommended)

**1. Install Everything:**
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Install Whisper.cpp (optional, for transcription)
brew install ffmpeg
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp && make
bash ./models/download-ggml-model.sh base.en
```

**2. Start Services:**
```bash
# Terminal 1: Start Ollama
ollama serve

# Terminal 2: Download model (one time)
ollama pull llama3.2

# Terminal 3: Start your app
cd /Users/catalinstratu/GolandProjects/interviuew_simulationChatbot
go run cmd/server/main.go
```

**3. Configure .env:**
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=interview_chatbot

# Server
SERVER_PORT=8080

# Free AI Configuration
USE_FREE_AI=true
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2

# Optional: If you want to use Gemini instead
# USE_GEMINI=true
# GEMINI_API_KEY=your_free_key
```

**4. Update Code to Use Free AI:**

Edit `cmd/server/main.go` line 36:
```go
// Change from:
aiService := services.NewAIService(os.Getenv("OPENAI_API_KEY"), db.DB, voiceService, questionService)

// To:
var aiService *services.AIHandler
if os.Getenv("USE_FREE_AI") == "true" {
    freeAI := services.NewFreeAIService(db.DB, voiceService, questionService)
    aiService = handlers.NewAIHandler(freeAI)
} else {
    // Fallback to OpenAI if configured
    aiService = services.NewAIService(os.Getenv("OPENAI_API_KEY"), db.DB, voiceService, questionService)
}
```

---

## Available Free Models (Ollama)

```bash
# Small & Fast (~2GB)
ollama pull llama3.2        # Meta's latest
ollama pull phi3            # Microsoft, very fast

# Medium (~4GB)
ollama pull mistral         # Excellent quality
ollama pull gemma2          # Google

# Large (~7GB) - Best quality
ollama pull llama3          # Meta's best
ollama pull mixtral         # Very powerful
```

Test models:
```bash
ollama run llama3.2 "Tell me about yourself"
```

---

## Troubleshooting

### Ollama not starting?
```bash
# Check if running
ps aux | grep ollama

# Restart
pkill ollama
ollama serve
```

### Model download slow?
- Base models: ~2GB (llama3.2, mistral)
- Use faster internet or smaller model (phi3)

### Analysis quality low?
- Try larger model: `ollama pull llama3` (7GB)
- Or use Gemini free tier (cloud-based, better quality)

### Whisper transcription not working?
- Install ffmpeg: `brew install ffmpeg`
- Check audio format is supported (WebM, MP3, WAV)
- Try manual transcription as fallback

---

## Cost Comparison

**5 Interview Sessions/Day:**

| Solution | Daily Cost | Monthly Cost | Annual Cost |
|----------|-----------|--------------|-------------|
| **Ollama** | $0 | $0 | $0 |
| **Gemini Free** | $0 | $0 | $0 |
| **OpenAI** | $1.05 | $31.50 | $378 |

---

## Next Steps

1. Install Ollama: `curl -fsSL https://ollama.com/install.sh | sh`
2. Download model: `ollama pull llama3.2`
3. Update `.env`: `USE_FREE_AI=true`
4. Run app: `go run cmd/server/main.go`
5. Test: Click "🤖 AI Analyze" in admin dashboard

**It just works!** No credit cards, no API keys, no costs. 🎉

---

## Support

- Ollama Docs: https://ollama.com
- Gemini Docs: https://ai.google.dev/
- Whisper.cpp: https://github.com/ggerganov/whisper.cpp
- Issues: Create GitHub issue in your project

Enjoy free AI analysis! 🚀
