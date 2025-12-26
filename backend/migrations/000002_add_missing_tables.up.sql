-- Add session_id to bookmarks for anonymous users
ALTER TABLE bookmarks ADD COLUMN IF NOT EXISTS session_id VARCHAR(255);
ALTER TABLE bookmarks ALTER COLUMN user_id DROP NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bookmarks_session_id ON bookmarks(session_id);

-- Reports table for content moderation
CREATE TABLE IF NOT EXISTS reports (
    id SERIAL PRIMARY KEY,
    reportable_type VARCHAR(50) NOT NULL CHECK (reportable_type IN ('tool', 'review')),
    reportable_id INT NOT NULL,
    reporter_user_id INT,
    reason VARCHAR(50) NOT NULL CHECK (reason IN ('spam', 'abuse', 'misinformation', 'other')),
    comment TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'reviewed', 'dismissed')) DEFAULT 'pending',
    reviewed_by INT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reporter_user_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_type ON reports(reportable_type);

-- Moderation actions audit log
CREATE TABLE IF NOT EXISTS moderation_actions (
    id SERIAL PRIMARY KEY,
    review_id INT NOT NULL,
    moderator_id INT NOT NULL,
    action_type VARCHAR(50) NOT NULL CHECK (action_type IN ('approve', 'hide', 'remove', 'restore')),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_moderation_actions_review ON moderation_actions(review_id);
