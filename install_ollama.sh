#!/bin/bash

echo "🚀 Installing FREE AI tools for Interview Chatbot"
echo "================================================="
echo ""

# Check if Ollama is installed
if command -v ollama &> /dev/null; then
    echo "✅ Ollama is already installed"
else
    echo "📦 Installing Ollama..."
    curl -fsSL https://ollama.com/install.sh | sh

    if [ $? -eq 0 ]; then
        echo "✅ Ollama installed successfully"
    else
        echo "❌ Failed to install Ollama"
        exit 1
    fi
fi

echo ""
echo "🤖 Downloading AI model (llama3.2 - ~2GB)..."
echo "This may take a few minutes..."

# Start Ollama in background if not running
if ! pgrep -x "ollama" > /dev/null; then
    echo "Starting Ollama service..."
    ollama serve &
    OLLAMA_PID=$!
    sleep 3
fi

# Download model
ollama pull llama3.2

if [ $? -eq 0 ]; then
    echo "✅ Model downloaded successfully"
else
    echo "❌ Failed to download model"
    exit 1
fi

echo ""
echo "🧪 Testing Ollama..."
ollama run llama3.2 "Say 'Hello! I am ready to analyze interviews!'" --verbose

echo ""
echo "================================================="
echo "✅ FREE AI Setup Complete!"
echo ""
echo "Next steps:"
echo "1. Make sure .env has: USE_FREE_AI=true"
echo "2. Start Ollama: ollama serve"
echo "3. Run your app: go run cmd/server/main.go"
echo "4. Click '🤖 AI Analyze' in admin dashboard"
echo ""
echo "🎉 Enjoy free AI analysis with no costs!"
