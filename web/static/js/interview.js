const API_BASE = '/api';

let currentSession = null;
let currentQuestion = null;
let mediaRecorder = null;
let audioChunks = [];
let recordedBlob = null;

// Screen elements
const startScreen = document.getElementById('startScreen');
const interviewScreen = document.getElementById('interviewScreen');
const completeScreen = document.getElementById('completeScreen');

// Start screen elements
const startBtn = document.getElementById('startBtn');
const candidateName = document.getElementById('candidateName');
const candidateEmail = document.getElementById('candidateEmail');
const numQuestions = document.getElementById('numQuestions');

// Interview screen elements
const progressFill = document.getElementById('progressFill');
const currentQuestionEl = document.getElementById('currentQuestion');
const totalQuestionsEl = document.getElementById('totalQuestions');
const questionText = document.getElementById('questionText');
const playAudioBtn = document.getElementById('playAudioBtn');
const questionAudio = document.getElementById('questionAudio');
const recordBtn = document.getElementById('recordBtn');
const stopBtn = document.getElementById('stopBtn');
const playBackBtn = document.getElementById('playBackBtn');
const playback = document.getElementById('playback');
const submitAnswerBtn = document.getElementById('submitAnswerBtn');
const recordingStatus = document.getElementById('recordingStatus');

// Complete screen elements
const sessionCode = document.getElementById('sessionCode');
const restartBtn = document.getElementById('restartBtn');

// Event listeners
startBtn.addEventListener('click', startInterview);
recordBtn.addEventListener('click', startRecording);
stopBtn.addEventListener('click', stopRecording);
playBackBtn.addEventListener('click', playRecording);
submitAnswerBtn.addEventListener('click', submitAnswer);
playAudioBtn.addEventListener('click', playQuestionAudio);
restartBtn.addEventListener('click', () => location.reload());

async function startInterview() {
    const data = {
        candidate_name: candidateName.value || '',
        candidate_email: candidateEmail.value || '',
        num_questions: parseInt(numQuestions.value)
    };

    try {
        const response = await fetch(`${API_BASE}/sessions/start`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Failed to start session');

        currentSession = await response.json();
        sessionCode.textContent = currentSession.session_code;
        totalQuestionsEl.textContent = currentSession.total_questions;

        showScreen('interview');
        await loadNextQuestion();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to start interview. Please try again.');
    }
}

async function loadNextQuestion() {
    try {
        const response = await fetch(`${API_BASE}/sessions/${currentSession.id}/next-question`);
        if (!response.ok) throw new Error('Failed to load question');

        const data = await response.json();

        if (data.is_complete) {
            await completeInterview();
            return;
        }

        currentQuestion = data.question;
        currentQuestionEl.textContent = data.question_index;
        questionText.textContent = currentQuestion.question_text;

        const progress = (data.question_index / data.total_questions) * 100;
        progressFill.style.width = `${progress}%`;

        // Load audio
        const audioUrl = `${API_BASE}/tts/${currentQuestion.id}?text=${encodeURIComponent(currentQuestion.question_text)}`;
        questionAudio.src = audioUrl;

        resetRecordingUI();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to load question. Please try again.');
    }
}

function playQuestionAudio() {
    questionAudio.play();
}

async function startRecording() {
    try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true });

        audioChunks = [];
        mediaRecorder = new MediaRecorder(stream);

        mediaRecorder.ondataavailable = (event) => {
            audioChunks.push(event.data);
        };

        mediaRecorder.onstop = () => {
            recordedBlob = new Blob(audioChunks, { type: 'audio/webm' });
            playback.src = URL.createObjectURL(recordedBlob);
            playback.style.display = 'block';
            playBackBtn.disabled = false;
            submitAnswerBtn.disabled = false;

            recordingStatus.classList.remove('recording');
            recordingStatus.querySelector('.status-text').textContent = 'Recording completed. You can play it back or submit.';
        };

        mediaRecorder.start();
        recordBtn.disabled = true;
        stopBtn.disabled = false;
        submitAnswerBtn.disabled = true;
        playBackBtn.disabled = true;

        recordingStatus.classList.add('recording');
        recordingStatus.querySelector('.status-text').textContent = 'Recording... Speak your answer now';
    } catch (error) {
        console.error('Error accessing microphone:', error);
        alert('Failed to access microphone. Please allow microphone access and try again.');
    }
}

function stopRecording() {
    if (mediaRecorder && mediaRecorder.state !== 'inactive') {
        mediaRecorder.stop();
        mediaRecorder.stream.getTracks().forEach(track => track.stop());
        stopBtn.disabled = true;
        recordBtn.disabled = false;
    }
}

function playRecording() {
    playback.play();
}

async function submitAnswer() {
    if (!recordedBlob) {
        alert('Please record your answer first');
        return;
    }

    const formData = new FormData();
    formData.append('question_id', currentQuestion.id);
    formData.append('audio', recordedBlob, 'answer.webm');

    try {
        submitAnswerBtn.disabled = true;
        submitAnswerBtn.textContent = 'Submitting...';

        const response = await fetch(`${API_BASE}/sessions/${currentSession.id}/answer`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) throw new Error('Failed to submit answer');

        await loadNextQuestion();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to submit answer. Please try again.');
        submitAnswerBtn.disabled = false;
        submitAnswerBtn.textContent = 'Submit Answer';
    }
}

async function completeInterview() {
    try {
        await fetch(`${API_BASE}/sessions/${currentSession.id}/complete`, {
            method: 'POST'
        });

        showScreen('complete');
    } catch (error) {
        console.error('Error:', error);
        showScreen('complete');
    }
}

function resetRecordingUI() {
    recordBtn.disabled = false;
    stopBtn.disabled = true;
    playBackBtn.disabled = true;
    submitAnswerBtn.disabled = true;
    submitAnswerBtn.textContent = 'Submit Answer';
    playback.style.display = 'none';
    recordedBlob = null;
    audioChunks = [];

    recordingStatus.classList.remove('recording');
    recordingStatus.querySelector('.status-text').textContent = 'Ready to record your answer';
}

function showScreen(screen) {
    startScreen.classList.remove('active');
    interviewScreen.classList.remove('active');
    completeScreen.classList.remove('active');

    switch(screen) {
        case 'start':
            startScreen.classList.add('active');
            break;
        case 'interview':
            interviewScreen.classList.add('active');
            break;
        case 'complete':
            completeScreen.classList.add('active');
            break;
    }
}
