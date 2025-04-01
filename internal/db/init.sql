CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    subject VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    newsgroup VARCHAR(255) NOT NULL,
    body TEXT NOT NULL CHECK (LENGTH(body) <= 10000),
    tags TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);