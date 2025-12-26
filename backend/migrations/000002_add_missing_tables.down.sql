-- Rollback migration
DROP TABLE IF EXISTS moderation_actions;
DROP TABLE IF EXISTS reports;
DROP INDEX IF EXISTS idx_bookmarks_session_id;
ALTER TABLE bookmarks DROP COLUMN IF EXISTS session_id;
ALTER TABLE bookmarks ALTER COLUMN user_id SET NOT NULL;
