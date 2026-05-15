/**
 * Voice Interview Controller
 * Automated hands-free interview flow with Romanian language support
 * Stop word: "STOP" - triggers next question after 2 seconds
 */

const VoiceInterview = (function() {
    'use strict';

    // Configuration
    const API_BASE = '/api';
    const MIN_RECORDING_MS = 2000;
    const SILENCE_LEVEL = 0.015;
    const STOP_WORD_DELAY = 2000; // 2 seconds after "STOP"
    const LANGUAGE = 'ro-RO'; // Romanian
    const STOP_WORDS = ['stop', 'gata', 'terminat', 'următoarea']; // Romanian stop words

    // State
    const state = {
        session: null,
        question: null,
        silenceThreshold: 3000,
        isRecording: false,
        recordingStart: null,
        stopWordDetected: false,
        speechRecognition: null
    };

    // Audio
    let mediaRecorder = null;
    let audioChunks = [];
    let audioContext = null;
    let analyser = null;
    let mediaStream = null;
    let stopWordTimeout = null;

    // DOM Elements
    const elements = {};

    /**
     * Initialize the application
     */
    function init() {
        console.log('Voice Interview: Initializing...');
        try {
            cacheElements();
            bindEvents();
            initSpeechRecognition();
            console.log('Voice Interview: Ready!');
        } catch (err) {
            console.error('Voice Interview: Init error:', err);
        }
    }

    /**
     * Initialize Web Speech API for stop word detection
     */
    function initSpeechRecognition() {
        const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition;

        if (!SpeechRecognition) {
            console.warn('Speech Recognition not supported');
            return;
        }

        state.speechRecognition = new SpeechRecognition();
        state.speechRecognition.continuous = true;
        state.speechRecognition.interimResults = true;
        state.speechRecognition.lang = LANGUAGE;

        state.speechRecognition.onresult = (event) => {
            if (!state.isRecording || state.stopWordDetected) return;

            for (let i = event.resultIndex; i < event.results.length; i++) {
                const transcript = event.results[i][0].transcript.toLowerCase().trim();
                console.log('Heard:', transcript);

                // Check for stop words
                const hasStopWord = STOP_WORDS.some(word => transcript.includes(word));

                if (hasStopWord) {
                    state.stopWordDetected = true;
                    setStatus('recording', 'STOP detectat! Se trece la următoarea întrebare...');
                    elements.stopWordIndicator.classList.remove('hidden');

                    // Start countdown
                    let countdown = 2;
                    elements.stopWordCountdown.textContent = countdown;

                    stopWordTimeout = setInterval(() => {
                        countdown--;
                        elements.stopWordCountdown.textContent = countdown;

                        if (countdown <= 0) {
                            clearInterval(stopWordTimeout);
                            stopRecording();
                        }
                    }, 1000);

                    break;
                }
            }
        };

        state.speechRecognition.onerror = (event) => {
            console.warn('Speech recognition error:', event.error);
            // Restart on recoverable errors
            if (event.error === 'no-speech' || event.error === 'audio-capture') {
                if (state.isRecording) {
                    try {
                        state.speechRecognition.start();
                    } catch (e) {}
                }
            }
        };

        state.speechRecognition.onend = () => {
            // Restart if still recording
            if (state.isRecording && !state.stopWordDetected) {
                try {
                    state.speechRecognition.start();
                } catch (e) {}
            }
        };
    }

    /**
     * Cache DOM elements
     */
    function cacheElements() {
        elements.setupScreen = document.getElementById('setupScreen');
        elements.interviewScreen = document.getElementById('interviewScreen');
        elements.completeScreen = document.getElementById('completeScreen');
        elements.startBtn = document.getElementById('startBtn');
        elements.candidateName = document.getElementById('candidateName');
        elements.candidateEmail = document.getElementById('candidateEmail');
        elements.numQuestions = document.getElementById('numQuestions');
        elements.silenceThreshold = document.getElementById('silenceThreshold');
        elements.statusDisplay = document.getElementById('statusDisplay');
        elements.statusText = document.getElementById('statusText');
        elements.questionText = document.getElementById('questionText');
        elements.currentQ = document.getElementById('currentQ');
        elements.totalQ = document.getElementById('totalQ');
        elements.progressBar = document.getElementById('progressBar');
        elements.visualizerSection = document.getElementById('visualizerSection');
        elements.visualizer = document.getElementById('visualizer');
        elements.voiceWave = document.getElementById('voiceWave');
        elements.silenceIndicator = document.getElementById('silenceIndicator');
        elements.silenceCountdown = document.getElementById('silenceCountdown');
        elements.sessionCode = document.getElementById('sessionCode');
        elements.stopWordIndicator = document.getElementById('stopWordIndicator');
        elements.stopWordCountdown = document.getElementById('stopWordCountdown');
    }

    /**
     * Bind event listeners
     */
    function bindEvents() {
        if (elements.startBtn) {
            elements.startBtn.addEventListener('click', startInterview);
            console.log('Voice Interview: Start button bound');
        } else {
            console.error('Voice Interview: Start button not found!');
        }
    }

    /**
     * Show specific screen
     */
    function showScreen(screenName) {
        document.querySelectorAll('.voice-screen').forEach(s => s.classList.remove('active'));

        const screen = {
            setup: elements.setupScreen,
            interview: elements.interviewScreen,
            complete: elements.completeScreen
        }[screenName];

        if (screen) screen.classList.add('active');
    }

    /**
     * Update status display
     */
    function setStatus(statusState, text) {
        elements.statusDisplay.dataset.state = statusState;
        elements.statusText.textContent = text;
    }

    /**
     * Start the interview
     */
    async function startInterview() {
        console.log('Voice Interview: Start button clicked');

        // Request microphone permission
        try {
            console.log('Voice Interview: Requesting microphone...');
            const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
            stream.getTracks().forEach(t => t.stop());
            console.log('Voice Interview: Microphone access granted');
        } catch (err) {
            console.error('Voice Interview: Microphone error:', err);
            alert('Accesul la microfon este necesar. Vă rugăm să permiteți și să încercați din nou.');
            return;
        }

        // Get settings
        state.silenceThreshold = parseInt(elements.silenceThreshold.value) * 1000;

        const payload = {
            candidate_name: elements.candidateName.value || '',
            candidate_email: elements.candidateEmail.value || '',
            num_questions: parseInt(elements.numQuestions.value)
        };

        try {
            elements.startBtn.disabled = true;
            showScreen('interview');
            setStatus('loading', 'Se pornește interviul...');

            const response = await fetch(`${API_BASE}/sessions/start`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!response.ok) throw new Error('Failed to start session');

            state.session = await response.json();
            elements.totalQ.textContent = state.session.total_questions;
            elements.sessionCode.textContent = state.session.session_code;

            await delay(1000);
            await loadNextQuestion();

        } catch (err) {
            console.error('Start error:', err);
            alert('Nu s-a putut porni interviul. Vă rugăm să încercați din nou.');
            showScreen('setup');
            elements.startBtn.disabled = false;
        }
    }

    /**
     * Load next question
     */
    async function loadNextQuestion() {
        try {
            setStatus('loading', 'Se încarcă întrebarea...');

            const response = await fetch(`${API_BASE}/sessions/${state.session.id}/next-question`);
            if (!response.ok) throw new Error('Failed to load question');

            const data = await response.json();

            if (data.is_complete) {
                await completeInterview();
                return;
            }

            state.question = data.question;
            elements.currentQ.textContent = data.question_index;
            elements.questionText.textContent = state.question.question_text;

            const progress = (data.question_index / data.total_questions) * 100;
            elements.progressBar.style.width = `${progress}%`;

            await delay(500);
            await playQuestion();

        } catch (err) {
            console.error('Load question error:', err);
            setStatus('error', 'Eroare la încărcarea întrebării');
            await delay(2000);
            await loadNextQuestion();
        }
    }

    /**
     * Play question audio via TTS
     */
    async function playQuestion() {
        return new Promise((resolve) => {
            setStatus('speaking', 'Ascultați întrebarea...');
            elements.voiceWave.classList.remove('hidden');

            const audioUrl = `${API_BASE}/tts/${state.question.id}?text=${encodeURIComponent(state.question.question_text)}&lang=ro`;
            const audio = new Audio(audioUrl);

            const onEnd = async () => {
                elements.voiceWave.classList.add('hidden');
                setStatus('idle', 'Pregătiți-vă să răspundeți...');
                await delay(1500);
                resolve();
                startRecording();
            };

            audio.onended = onEnd;
            audio.onerror = () => {
                console.warn('TTS audio failed');
                onEnd();
            };

            audio.play().catch(() => {
                setStatus('idle', 'Citiți întrebarea de mai sus...');
                delay(4000).then(onEnd);
            });
        });
    }

    /**
     * Start recording
     */
    async function startRecording() {
        try {
            setStatus('recording', 'Vorbiți răspunsul... (spuneți "STOP" pentru a termina)');
            elements.visualizerSection.classList.remove('hidden');
            state.stopWordDetected = false;
            elements.stopWordIndicator.classList.add('hidden');

            if (stopWordTimeout) {
                clearInterval(stopWordTimeout);
                stopWordTimeout = null;
            }

            mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true });

            // Setup audio analysis
            audioContext = new AudioContext();
            const source = audioContext.createMediaStreamSource(mediaStream);
            analyser = audioContext.createAnalyser();
            analyser.fftSize = 256;
            source.connect(analyser);

            // Setup recorder
            audioChunks = [];
            mediaRecorder = new MediaRecorder(mediaStream);

            mediaRecorder.ondataavailable = (e) => audioChunks.push(e.data);

            mediaRecorder.onstop = async () => {
                state.isRecording = false;
                elements.visualizerSection.classList.add('hidden');
                elements.silenceIndicator.classList.add('hidden');
                elements.stopWordIndicator.classList.add('hidden');

                // Stop speech recognition
                if (state.speechRecognition) {
                    try {
                        state.speechRecognition.stop();
                    } catch (e) {}
                }

                mediaStream.getTracks().forEach(t => t.stop());
                if (audioContext) {
                    audioContext.close();
                    audioContext = null;
                }

                if (audioChunks.length === 0) {
                    setStatus('error', 'Nu s-a capturat audio, se reîncearcă...');
                    await delay(1500);
                    startRecording();
                    return;
                }

                const blob = new Blob(audioChunks, { type: 'audio/webm' });
                await submitAnswer(blob);
            };

            mediaRecorder.start();
            state.isRecording = true;
            state.recordingStart = Date.now();

            // Start speech recognition for stop word
            if (state.speechRecognition) {
                try {
                    state.speechRecognition.start();
                } catch (e) {
                    console.warn('Could not start speech recognition:', e);
                }
            }

            drawVisualizer();
            detectSilence();

        } catch (err) {
            console.error('Recording error:', err);
            setStatus('error', 'Eroare microfon');
            await delay(2000);
            startRecording();
        }
    }

    /**
     * Stop recording
     */
    function stopRecording() {
        if (stopWordTimeout) {
            clearInterval(stopWordTimeout);
            stopWordTimeout = null;
        }

        if (state.speechRecognition) {
            try {
                state.speechRecognition.stop();
            } catch (e) {}
        }

        if (mediaRecorder && mediaRecorder.state !== 'inactive') {
            setStatus('processing', 'Se procesează...');
            mediaRecorder.stop();
        }
    }

    /**
     * Detect silence to auto-stop recording
     */
    function detectSilence() {
        const dataArray = new Uint8Array(analyser.frequencyBinCount);
        let silentFrames = 0;
        const framesNeeded = Math.ceil(state.silenceThreshold / 50);

        function check() {
            if (!state.isRecording || state.stopWordDetected) return;

            analyser.getByteFrequencyData(dataArray);
            const avg = dataArray.reduce((a, b) => a + b, 0) / dataArray.length;
            const volume = avg / 255;

            const elapsed = Date.now() - state.recordingStart;

            if (volume < SILENCE_LEVEL && elapsed > MIN_RECORDING_MS) {
                silentFrames++;

                const silentMs = silentFrames * 50;
                const remaining = Math.ceil((state.silenceThreshold - silentMs) / 1000);

                if (remaining > 0 && remaining <= Math.ceil(state.silenceThreshold / 1000)) {
                    elements.silenceIndicator.classList.remove('hidden');
                    elements.silenceCountdown.textContent = remaining;

                    const progress = (silentMs / state.silenceThreshold) * 100;
                    elements.silenceIndicator.querySelector('.silence-bar').style.width = `${progress}%`;
                }

                if (silentFrames >= framesNeeded) {
                    stopRecording();
                    return;
                }
            } else {
                silentFrames = 0;
                elements.silenceIndicator.classList.add('hidden');
            }

            requestAnimationFrame(check);
        }

        check();
    }

    /**
     * Draw audio visualizer
     */
    function drawVisualizer() {
        const canvas = elements.visualizer;
        const ctx = canvas.getContext('2d');

        // Set canvas size
        canvas.width = canvas.offsetWidth * 2;
        canvas.height = canvas.offsetHeight * 2;
        ctx.scale(2, 2);

        const WIDTH = canvas.offsetWidth;
        const HEIGHT = canvas.offsetHeight;

        function draw() {
            if (!state.isRecording || !analyser) return;

            requestAnimationFrame(draw);

            const dataArray = new Uint8Array(analyser.frequencyBinCount);
            analyser.getByteFrequencyData(dataArray);

            ctx.fillStyle = '#1a1a2e';
            ctx.fillRect(0, 0, WIDTH, HEIGHT);

            const barCount = 50;
            const barWidth = WIDTH / barCount;
            const step = Math.floor(dataArray.length / barCount);

            for (let i = 0; i < barCount; i++) {
                const value = dataArray[i * step];
                const barHeight = (value / 255) * HEIGHT * 0.9;

                const hue = 250 + (value / 255) * 60;
                ctx.fillStyle = `hsl(${hue}, 70%, 60%)`;

                const x = i * barWidth;
                const y = HEIGHT - barHeight;
                const radius = barWidth / 3;

                ctx.beginPath();
                ctx.roundRect(x + 1, y, barWidth - 2, barHeight, radius);
                ctx.fill();
            }
        }

        draw();
    }

    /**
     * Submit answer
     */
    async function submitAnswer(blob) {
        setStatus('processing', 'Se trimite răspunsul...');

        const formData = new FormData();
        formData.append('question_id', state.question.id);
        formData.append('audio', blob, 'answer.webm');

        try {
            const response = await fetch(`${API_BASE}/sessions/${state.session.id}/answer`, {
                method: 'POST',
                body: formData
            });

            if (!response.ok) throw new Error('Submit failed');

            setStatus('success', 'Răspuns trimis!');
            await delay(1200);
            await loadNextQuestion();

        } catch (err) {
            console.error('Submit error:', err);
            setStatus('error', 'Eroare la trimitere, se reîncearcă...');
            await delay(2000);
            await submitAnswer(blob);
        }
    }

    /**
     * Complete interview
     */
    async function completeInterview() {
        try {
            await fetch(`${API_BASE}/sessions/${state.session.id}/complete`, {
                method: 'POST'
            });
        } catch (err) {
            console.error('Complete error:', err);
        }

        showScreen('complete');

        // Speak completion message in Romanian
        if ('speechSynthesis' in window) {
            const msg = new SpeechSynthesisUtterance(
                'Interviul s-a încheiat. Vă mulțumim pentru răspunsuri.'
            );
            msg.lang = LANGUAGE;
            speechSynthesis.speak(msg);
        }
    }

    /**
     * Utility: delay
     */
    function delay(ms) {
        return new Promise(r => setTimeout(r, ms));
    }

    // Initialize on DOM ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    return { state };
})();