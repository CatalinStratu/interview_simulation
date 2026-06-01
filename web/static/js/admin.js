const API_BASE = '/api';

// Tab elements
const questionsTabBtn = document.getElementById('questionsTabBtn');
const sessionsTabBtn = document.getElementById('sessionsTabBtn');
const questionsTab = document.getElementById('questionsTab');
const sessionsTab = document.getElementById('sessionsTab');

// Questions elements
const questionText = document.getElementById('questionText');
const questionType = document.getElementById('questionType');
const difficultyLevel = document.getElementById('difficultyLevel');
const expectedDuration = document.getElementById('expectedDuration');
const addQuestionBtn = document.getElementById('addQuestionBtn');
const questionsList = document.getElementById('questionsList');

// Sessions elements
const sessionsList = document.getElementById('sessionsList');
const sessionModal = document.getElementById('sessionModal');
const sessionDetails = document.getElementById('sessionDetails');
const remarkText = document.getElementById('remarkText');
const remarkRating = document.getElementById('remarkRating');
const addRemarkBtn = document.getElementById('addRemarkBtn');
const remarksList = document.getElementById('remarksList');

let currentSessionId = null;
let currentQuestions = [];

// Event listeners
questionsTabBtn.addEventListener('click', () => switchTab('questions'));
sessionsTabBtn.addEventListener('click', () => switchTab('sessions'));
addQuestionBtn.addEventListener('click', addQuestion);
addRemarkBtn.addEventListener('click', addRemark);

// Modal close
const closeModal = document.querySelector('.close');
closeModal.addEventListener('click', () => {
    sessionModal.style.display = 'none';
});

window.addEventListener('click', (event) => {
    if (event.target === sessionModal) {
        sessionModal.style.display = 'none';
    }
});

// Initialize
loadQuestions();

function switchTab(tab) {
    questionsTabBtn.classList.remove('active');
    sessionsTabBtn.classList.remove('active');
    questionsTab.classList.remove('active');
    sessionsTab.classList.remove('active');

    if (tab === 'questions') {
        questionsTabBtn.classList.add('active');
        questionsTab.classList.add('active');
        loadQuestions();
    } else {
        sessionsTabBtn.classList.add('active');
        sessionsTab.classList.add('active');
        loadSessions();
    }
}

async function loadQuestions() {
    try {
        const response = await fetch(`${API_BASE}/admin/questions`);
        if (!response.ok) throw new Error('Failed to load questions');

        const questions = await response.json();
        displayQuestions(questions);
    } catch (error) {
        console.error('Error:', error);
        questionsList.innerHTML = '<p class="loading">Failed to load questions</p>';
    }
}

function displayQuestions(questions) {
    currentQuestions = questions || [];

    if (currentQuestions.length === 0) {
        questionsList.innerHTML = '<p class="loading">No questions available</p>';
        return;
    }

    questionsList.innerHTML = currentQuestions.map((q, i) => `
        <div class="question-item">
            <div class="question-order">
                <button class="btn-move" title="Move up" onclick="moveQuestion(${i}, -1)" ${i === 0 ? 'disabled' : ''}>▲</button>
                <span class="order-num">${i + 1}</span>
                <button class="btn-move" title="Move down" onclick="moveQuestion(${i}, 1)" ${i === currentQuestions.length - 1 ? 'disabled' : ''}>▼</button>
            </div>
            <div class="question-body">
                <h4>${escapeHtml(q.question_text)}</h4>
                <div class="question-meta">
                    <span class="badge badge-${q.question_type}">${q.question_type}</span>
                    <span class="badge badge-${q.difficulty_level}">${q.difficulty_level}</span>
                    <span>Duration: ${q.expected_duration}s</span>
                </div>
                <button class="btn-delete" onclick="deleteQuestion(${q.id})">Delete</button>
            </div>
        </div>
    `).join('');
}

// Move a question up (-1) or down (+1), persist the new order, and re-render.
async function moveQuestion(index, direction) {
    const target = index + direction;
    if (target < 0 || target >= currentQuestions.length) return;

    // Swap locally and re-render immediately for snappy feedback.
    const reordered = currentQuestions.slice();
    [reordered[index], reordered[target]] = [reordered[target], reordered[index]];
    displayQuestions(reordered);

    try {
        const response = await fetch(`${API_BASE}/admin/questions/reorder`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ ids: reordered.map(q => q.id) })
        });
        if (!response.ok) throw new Error('Failed to save order');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to save order');
        await loadQuestions(); // reload the real order on failure
    }
}

async function addQuestion() {
    const data = {
        question_text: questionText.value,
        question_type: questionType.value,
        difficulty_level: difficultyLevel.value,
        expected_duration: parseInt(expectedDuration.value)
    };

    if (!data.question_text) {
        alert('Please enter a question');
        return;
    }

    try {
        addQuestionBtn.disabled = true;
        addQuestionBtn.textContent = 'Adding...';

        const response = await fetch(`${API_BASE}/admin/questions`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Failed to add question');

        questionText.value = '';
        await loadQuestions();
        alert('Question added successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add question');
    } finally {
        addQuestionBtn.disabled = false;
        addQuestionBtn.textContent = 'Add Question';
    }
}

async function deleteQuestion(id) {
    if (!confirm('Are you sure you want to delete this question?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/admin/questions/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete question');

        await loadQuestions();
        alert('Question deleted successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to delete question');
    }
}

async function loadSessions() {
    try {
        const response = await fetch(`${API_BASE}/admin/sessions`);
        if (!response.ok) throw new Error('Failed to load sessions');

        const sessions = await response.json();
        displaySessions(sessions);
    } catch (error) {
        console.error('Error:', error);
        sessionsList.innerHTML = '<p class="loading">Failed to load sessions</p>';
    }
}

function displaySessions(sessions) {
    if (!sessions || sessions.length === 0) {
        sessionsList.innerHTML = '<p class="loading">No sessions available</p>';
        return;
    }

    sessionsList.innerHTML = sessions.map(s => {
        const candidateName = s.candidate_name || 'Anonymous';
        const startedAt = new Date(s.started_at).toLocaleString();
        const completedAt = s.completed_at ? new Date(s.completed_at).toLocaleString() : 'In Progress';

        return `
            <div class="session-item">
                <h4>Session: ${s.session_code}</h4>
                <div class="session-meta">
                    <span><strong>Candidate:</strong> ${escapeHtml(candidateName)}</span>
                    <span class="badge badge-${s.status}">${s.status}</span>
                    <span><strong>Progress:</strong> ${s.questions_answered}/${s.total_questions}</span>
                </div>
                <div class="session-meta">
                    <span><strong>Started:</strong> ${startedAt}</span>
                    <span><strong>Completed:</strong> ${completedAt}</span>
                </div>
                <button class="btn-view" onclick="viewSession(${s.id})">View Details</button>
                <button class="btn-analyze" onclick="analyzeSession(${s.id})">🤖 AI Analyze</button>
            </div>
        `;
    }).join('');
}

async function viewSession(sessionId) {
    currentSessionId = sessionId;

    try {
        // Load session details, answers, and the question catalog in parallel
        const [sessionResponse, answersResponse, questionsResponse] = await Promise.all([
            fetch(`${API_BASE}/admin/sessions/${sessionId}`),
            fetch(`${API_BASE}/sessions/${sessionId}/answers`),
            fetch(`${API_BASE}/admin/questions`),
        ]);

        if (!sessionResponse.ok) throw new Error('Failed to load session');
        const session = await sessionResponse.json();
        const answers = answersResponse.ok ? await answersResponse.json() : [];
        const questions = questionsResponse.ok ? await questionsResponse.json() : [];

        // Map question_id -> question_text for display
        const questionTextById = new Map((questions || []).map(q => [q.id, q.question_text]));

        // Load remarks
        await loadRemarks(sessionId);

        // Display session details
        const candidateName = session.candidate_name || 'Anonymous';
        const candidateEmail = session.candidate_email || 'Not provided';

        sessionDetails.innerHTML = `
            <div class="session-info">
                <p><strong>Session Code:</strong> ${session.session_code}</p>
                <p><strong>Candidate:</strong> ${escapeHtml(candidateName)}</p>
                <p><strong>Email:</strong> ${escapeHtml(candidateEmail)}</p>
                <p><strong>Status:</strong> <span class="badge badge-${session.status}">${session.status}</span></p>
                <p><strong>Questions:</strong> ${session.questions_answered}/${session.total_questions}</p>
                <p><strong>Started:</strong> ${new Date(session.started_at).toLocaleString()}</p>
            </div>

            <h3>Voice Answers (${answers.length})</h3>
            ${answers.length > 0 ? answers.map((a, i) => {
                const qText = questionTextById.get(a.question_id) || `Question #${a.question_id}`;
                return `
                <div class="remark-item">
                    <p><strong>Q${i + 1}.</strong> ${escapeHtml(qText)}</p>
                    <p><strong>Answered:</strong> ${new Date(a.answered_at).toLocaleString()}</p>
                    <audio controls preload="none" src="${API_BASE}/admin/answers/${a.id}/audio" style="width:100%;margin:6px 0;">
                        Your browser does not support audio playback.
                    </audio>
                    ${a.transcription ? `
                        <p><strong>Transcription:</strong> ${escapeHtml(a.transcription)}</p>
                    ` : ''}
                    ${a.ai_analysis ? `
                        <div style="background: #f0f8ff; padding: 10px; margin-top: 10px; border-radius: 5px;">
                            <p><strong>AI Analysis:</strong></p>
                            <pre style="white-space: pre-wrap;">${escapeHtml(a.ai_analysis)}</pre>
                            ${a.ai_score ? `<p><strong>Score:</strong> ${a.ai_score}/100</p>` : ''}
                        </div>
                    ` : '<p><em>No AI analysis yet</em></p>'}
                </div>
            `;
            }).join('') : '<p>No answers recorded yet</p>'}
        `;

        sessionModal.style.display = 'block';
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to load session details');
    }
}

async function loadRemarks(sessionId) {
    try {
        const response = await fetch(`${API_BASE}/admin/sessions/${sessionId}/remarks`);
        if (!response.ok) throw new Error('Failed to load remarks');

        const remarks = await response.json();

        if (remarks && remarks.length > 0) {
            remarksList.innerHTML = `
                <h3>Remarks (${remarks.length})</h3>
                ${remarks.map(r => `
                    <div class="remark-item">
                        <p>${escapeHtml(r.remark_text)}</p>
                        ${r.rating ? `<p><strong>Rating:</strong> ${r.rating}/5</p>` : ''}
                        <p><small>By ${r.created_by} on ${new Date(r.created_at).toLocaleString()}</small></p>
                    </div>
                `).join('')}
            `;
        } else {
            remarksList.innerHTML = '<p>No remarks yet</p>';
        }
    } catch (error) {
        console.error('Error:', error);
        remarksList.innerHTML = '<p>Failed to load remarks</p>';
    }
}

async function addRemark() {
    if (!currentSessionId) return;

    const data = {
        remark_text: remarkText.value,
        rating: remarkRating.value ? parseInt(remarkRating.value) : null
    };

    if (!data.remark_text) {
        alert('Please enter a remark');
        return;
    }

    try {
        addRemarkBtn.disabled = true;
        addRemarkBtn.textContent = 'Adding...';

        const response = await fetch(`${API_BASE}/admin/sessions/${currentSessionId}/remarks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) throw new Error('Failed to add remark');

        remarkText.value = '';
        remarkRating.value = '';
        await loadRemarks(currentSessionId);
        alert('Remark added successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add remark');
    } finally {
        addRemarkBtn.disabled = false;
        addRemarkBtn.textContent = 'Add Remark';
    }
}

async function analyzeSession(sessionId) {
    if (!confirm('This will analyze all voice answers using AI. This may take a few minutes and will use OpenAI API credits. Continue?')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/ai/analyze/session/${sessionId}`, {
            method: 'POST'
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to analyze session');
        }

        alert('AI analysis completed! Reload the session to see results.');
        await loadSessions();
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to analyze session: ' + error.message);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
