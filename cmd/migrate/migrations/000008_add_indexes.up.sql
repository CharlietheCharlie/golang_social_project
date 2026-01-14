

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Add GIN indexes for full-text search on comments content, posts title, and posts tags
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags); -- not gin_trgm_ops as tags are not text
--- Add B-tree indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts (user_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments (post_id);