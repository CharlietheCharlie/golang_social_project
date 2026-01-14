package main

import (
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {
	dsn := env.GetString("DB_DSN", "postgres://user:password@localhost:5432/social?sslmode=disable")
	conn, err := db.New(dsn, 3, 3, "15m")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Initialize the storage layer
	store := store.NewStorage(conn)

	db.Seed(store)

}
