#!/usr/bin/env bash
# One-time setup for the local Romanian MMS-TTS sidecar.
# Usage:
#   bash tts_mms/setup.sh
#   source tts_mms/.venv/bin/activate && python tts_mms/server.py
set -euo pipefail

cd "$(dirname "$0")"

if [ ! -d ".venv" ]; then
  python3 -m venv .venv
fi

# shellcheck disable=SC1091
source .venv/bin/activate
pip install --upgrade pip
pip install -r requirements.txt

echo
echo "Setup complete."
echo "Start the sidecar:"
echo "  source tts_mms/.venv/bin/activate && python tts_mms/server.py"