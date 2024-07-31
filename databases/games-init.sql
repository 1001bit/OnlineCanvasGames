CREATE TABLE IF NOT EXISTS games (
    title VARCHAR(63) PRIMARY KEY
);

INSERT INTO games (title) VALUES ('clicker') ON CONFLICT (title) DO NOTHING;
INSERT INTO games (title) VALUES ('platformer') ON CONFLICT (title) DO NOTHING;