-- +goose Up
ALTER TABLE chirps 
ADD column hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
ALTER TABLE chirps
DROP column hashed_password;