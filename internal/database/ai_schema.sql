-- Add transcription and AI analysis columns to voice_answers table
ALTER TABLE voice_answers
ADD COLUMN transcription_status VARCHAR(20) DEFAULT 'pending',
ADD COLUMN ai_analysis TEXT,
ADD COLUMN ai_score INT,
ADD COLUMN analyzed_at TIMESTAMP NULL;

-- Create index for faster queries
CREATE INDEX idx_transcription_status ON voice_answers(transcription_status);
