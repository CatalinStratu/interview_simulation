# Quick Start: FREE AI Analysis (2 Minutes)

## 🎯 Goal
Add AI analysis to your interview chatbot **for free** using Ollama.

## ⚡ Quick Setup (3 commands)

```bash
# 1. Install Ollama (free, local AI)
curl -fsSL https://ollama.com/install.sh | sh

# 2. Download AI model (~2GB, one time)
ollama pull llama3.2

# 3. Start Ollama
ollama serve
```

**That's it!** Your `.env` is already configured with `USE_FREE_AI=true`.

## ▶️ Run Your App

```bash
# Terminal 1: Keep Ollama running
ollama serve

# Terminal 2: Start your app
go run cmd/server/main.go
```

You should see:
```
Using FREE AI (Ollama) for analysis
Server starting on http://localhost:8080
```

## 🧪 Test It

1. Go to http://localhost:8080/
2. Complete an interview (record voice answers)
3. Go to http://localhost:8080/admin
4. Click **"🤖 AI Analyze"** on your session
5. Wait 30-60 seconds
6. Click **"View Details"** to see AI feedback!

## 📦 What You Get (FREE)

- ✅ Speech-to-text transcription (optional with whisper.cpp)
- ✅ Detailed answer analysis
- ✅ Strengths identification
- ✅ Improvement suggestions
- ✅ Numerical scoring (0-100)
- ✅ **All FREE, runs locally, no API keys**

## 🚀 Alternative: Auto-Install Script

```bash
./install_ollama.sh
```

This script:
1. Installs Ollama
2. Downloads llama3.2 model
3. Tests everything
4. Shows you next steps

## 💡 Other Free Models

Try different models:

```bash
# Fast and small (~1.5GB)
ollama pull phi3

# Excellent quality (~4GB)
ollama pull mistral

# Google's model (~5GB)
ollama pull gemma2

# Best quality (~7GB)
ollama pull llama3
```

Change model in `.env`:
```bash
OLLAMA_MODEL=mistral
```

## 🔧 Troubleshooting

### "Ollama not running"
```bash
# Start it
ollama serve
```

### "Model not found"
```bash
# Download it
ollama pull llama3.2
```

### "Analysis failed"
```bash
# Check Ollama is running
curl http://localhost:11434/api/version

# Should return: {"version":"0.x.x"}
```

## 📊 Cost Comparison

| Solution | Setup Time | Monthly Cost | Quality |
|----------|-----------|--------------|---------|
| **Ollama (Free)** | 2 minutes | $0 | Very Good |
| OpenAI | 1 minute | $30+ | Excellent |
| Gemini Free | 2 minutes | $0 | Excellent |

## 🎓 How It Works

```
Voice Answer
    ↓
Ollama AI (running locally on your Mac)
    ↓
Analysis, Feedback, Score
    ↓
Saved to Database
```

**No data leaves your computer!**

## ⚙️ Configuration

Your `.env` is already set up:

```bash
USE_FREE_AI=true              # Use free Ollama
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.2         # Which model to use
```

## 📚 Learn More

- **Ollama**: https://ollama.com
- **Models**: https://ollama.com/library
- **Full guide**: See `FREE_AI_SETUP.md`

## 🆘 Need Help?

Common issues:

1. **Port 11434 in use**: Another Ollama instance running
2. **Slow first time**: Model downloading (2GB)
3. **Analysis slow**: Normal, ~30-60 seconds for good analysis

## 🎉 You're Done!

Enjoy **free, unlimited** AI analysis with no API keys, no costs, and full privacy!

Questions? Check `FREE_AI_SETUP.md` for detailed docs.
