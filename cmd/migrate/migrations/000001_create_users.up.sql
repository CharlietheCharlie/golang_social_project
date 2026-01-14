-- migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users
-- migrate -path ./cmd/migrate/migrations -database "postgres://user:password@localhost:5432/social?sslmode=disable" up
CREATE EXTENSION IF NOT EXISTS "citext";

CREATE TABLE
    IF NOT EXISTS users (
        id BIGSERIAL PRIMARY KEY,
        email citext UNIQUE NOT NULL, --case-insensitive text
        username VARCHAR(255) UNIQUE NOT NULL,
        password BYTEA NOT NULL,
        created_at TIMESTAMP(0)
        WITH
            TIME ZONE NOT NULL DEFAULT NOW ()
    );