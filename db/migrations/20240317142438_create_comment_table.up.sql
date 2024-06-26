CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY NOT NULL,
    comment VARCHAR NOT NULL,
    post_id UUID REFERENCES posts(id) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);
