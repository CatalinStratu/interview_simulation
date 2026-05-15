-- Interview Questions Table
CREATE TABLE IF NOT EXISTS questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) DEFAULT 'technical', -- technical, behavioral, etc.
    difficulty_level VARCHAR(20) DEFAULT 'medium', -- easy, medium, hard
    expected_duration INT DEFAULT 120, -- expected answer duration in seconds
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Interview Sessions Table
CREATE TABLE IF NOT EXISTS interview_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_code VARCHAR(50) UNIQUE NOT NULL,
    candidate_name VARCHAR(255),
    candidate_email VARCHAR(255),
    status VARCHAR(20) DEFAULT 'in_progress', -- in_progress, completed, abandoned
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    total_questions INT DEFAULT 0,
    questions_answered INT DEFAULT 0
);

-- Voice Answers Table
CREATE TABLE IF NOT EXISTS voice_answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    question_id INT NOT NULL,
    audio_file_path VARCHAR(500) NOT NULL,
    audio_duration INT, -- duration in seconds
    transcription TEXT, -- optional: if we add speech-to-text later
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES interview_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
);

-- Remarks Table
CREATE TABLE IF NOT EXISTS remarks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    remark_text TEXT NOT NULL,
    rating INT, -- optional rating 1-5
    created_by VARCHAR(255) DEFAULT 'admin',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES interview_sessions(id) ON DELETE CASCADE
);

-- Session Questions Junction Table (to track which questions were asked in each session)
CREATE TABLE IF NOT EXISTS session_questions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT NOT NULL,
    question_id INT NOT NULL,
    question_order INT NOT NULL,
    asked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES interview_sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    UNIQUE KEY unique_session_question (session_id, question_id)
);

-- Insert some sample questions
INSERT INTO questions (question_text, question_type, difficulty_level) VALUES
('Tell me about yourself and your professional background.', 'behavioral', 'easy'),
('What are your greatest strengths and weaknesses?', 'behavioral', 'easy'),
('Explain the concept of object-oriented programming.', 'technical', 'medium'),
('What is the difference between a process and a thread?', 'technical', 'medium'),
('How would you design a URL shortening service?', 'technical', 'hard'),
('Describe a challenging project you worked on and how you overcame obstacles.', 'behavioral', 'medium'),
('What is a RESTful API and what are its key principles?', 'technical', 'medium'),
('How do you handle conflicts in a team environment?', 'behavioral', 'medium'),
('Explain the CAP theorem in distributed systems.', 'technical', 'hard'),
('Where do you see yourself in five years?', 'behavioral', 'easy');
