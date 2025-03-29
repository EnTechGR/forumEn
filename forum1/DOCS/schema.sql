-- Users table
CREATE TABLE IF NOT EXISTS user (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE 
        CHECK (length(username) >= 3 AND length(username) <= 15)
        CHECK (username GLOB '[a-zA-Z0-9_]*')
        CHECK (username NOT GLOB '*[^a-zA-Z0-9_]*'),
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User authentication table
CREATE TABLE IF NOT EXISTS user_auth (
    user_id TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL 
    CHECK (length(password_hash) = 60),
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
);

-- Sessions table (one-to-one with user)
CREATE TABLE sessions (
    user_id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL UNIQUE,
    ip_address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
);

-- Categories table
CREATE TABLE categories (
    category_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Posts table
CREATE TABLE posts (
    post_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

-- Comments table
CREATE TABLE comments (
    comment_id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE
);

-- Reactions table
CREATE TABLE IF NOT EXISTS reactions (
    reaction_id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    reaction_type INTEGER NOT NULL, -- 1 for like, 2 for dislike
    comment_id TEXT,
    post_id TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY (user_id) REFERENCES user(user_id) ON DELETE CASCADE,
    KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE,
    CHECK ((post_id IS NULL AND comment_id IS NOT NULL) OR (post_id IS NOT NULL AND comment_id IS NULL))
);

-- Create necessary indexes
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);,
CREATE INDEX IF NOT EXISTS idx_posts_category_id ON posts(category_id);,
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON comments(post_id);,
CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);,
CREATE INDEX IF NOT EXISTS idx_reactions_user_id ON reactions(user_id);,
CREATE INDEX IF NOT EXISTS idx_reactions_post_id ON reactions(post_id);,
CREATE INDEX IF NOT EXISTS idx_reactions_comment_id ON reactions(comment_id);,
CREATE INDEX IF NOT EXISTS idx_sessions_session_id ON sessions(session_id);,
CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_user_post ON reactions(user_id, post_id) WHERE post_id IS NOT NULL;,
CREATE UNIQUE INDEX IF NOT EXISTS idx_reactions_user_comment ON reactions(user_id, comment_id) WHERE comment_id IS NOT NULL;,