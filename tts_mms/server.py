"""Local Facebook MMS-TTS (Romanian) sidecar.

POST /tts  body: {"text": "..."}  -> audio/wav
GET  /health                       -> {"ok": true}
"""

import io
import os
import sys

import numpy as np
import scipy.io.wavfile
import torch
from flask import Flask, Response, jsonify, request
from transformers import AutoTokenizer, VitsModel

MODEL_NAME = os.environ.get("MMS_MODEL", "facebook/mms-tts-ron")
HOST = os.environ.get("MMS_HOST", "127.0.0.1")
PORT = int(os.environ.get("MMS_PORT", "5005"))

print(f"[mms-tts] loading {MODEL_NAME} ...", flush=True)
model = VitsModel.from_pretrained(MODEL_NAME)
tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
model.eval()
print("[mms-tts] model ready", flush=True)

app = Flask(__name__)


@app.get("/health")
def health():
    return jsonify(ok=True, model=MODEL_NAME)


@app.post("/tts")
def tts():
    data = request.get_json(silent=True) or {}
    text = (data.get("text") or "").strip()
    if not text:
        return jsonify(error="text required"), 400

    inputs = tokenizer(text, return_tensors="pt")
    with torch.no_grad():
        waveform = model(**inputs).waveform

    audio = waveform.squeeze().cpu().numpy()
    audio = np.clip(audio, -1.0, 1.0)
    pcm16 = (audio * 32767).astype(np.int16)

    buf = io.BytesIO()
    scipy.io.wavfile.write(buf, model.config.sampling_rate, pcm16)
    buf.seek(0)
    return Response(buf.read(), mimetype="audio/wav")


if __name__ == "__main__":
    print(f"[mms-tts] listening on http://{HOST}:{PORT}", flush=True)
    app.run(host=HOST, port=PORT, threaded=True)
