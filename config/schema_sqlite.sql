-- SQLite uses dynamic typing; INTEGER PRIMARY KEY is implicitly an alias for rowid (autoincrement).

-- users DDL
CREATE TABLE IF NOT EXISTS users (
    id          INTEGER  NOT NULL PRIMARY KEY AUTOINCREMENT,
    username    TEXT     NOT NULL,
    nickname    TEXT,
    email       TEXT    NOT NULL,
    kind        TEXT,
    status      TEXT    NOT NULL
);

CREATE INDEX idx_users_name ON users(username);

INSERT INTO users (username,nickname,email,kind,status) VALUES
( 'Tom','Tommy','tom@example.com','student','learning'),
('Jerry','Jerrie','jerry@example.com','student','learning');
