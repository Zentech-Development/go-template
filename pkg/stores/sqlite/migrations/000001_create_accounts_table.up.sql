CREATE TABLE IF NOT EXISTS accounts (
    id       INTEGER NOT NULL PRIMARY KEY, 
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);