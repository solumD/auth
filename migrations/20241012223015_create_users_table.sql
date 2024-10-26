-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, 
    role INT DEFAULT 0,     /* 0 - UNKNOWN, 1 - USER, 2 - ADMIN */
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP 
);

-- +goose Down
DROP TABLE users;