-- migrate -path ./cmd/migrate/migrations -database "postgres://user:password@localhost:5432/social?sslmode=disable" down

DROP TABLE IF EXISTS users;