CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    title VARCHAR(63) NOT NULL UNIQUE
);

INSERT INTO games (title) VALUES ('clicker') ON CONFLICT (title) DO NOTHING;
INSERT INTO games (title) VALUES ('platformer') ON CONFLICT (title) DO NOTHING;