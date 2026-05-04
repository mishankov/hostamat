-- +migrate ID: init
-- +migrate Up
CREATE TABLE IF NOT EXISTS feedback (
	id VARCHAR(255) PRIMARY KEY,
	user_id VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	url VARCHAR(2048) NOT NULL,
	message TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_feedback_created_at ON feedback(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedback_user_id ON feedback(user_id);
CREATE INDEX IF NOT EXISTS idx_feedback_user_created_at ON feedback(user_id, created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS feedback;
