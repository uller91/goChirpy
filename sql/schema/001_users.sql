-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email text UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
