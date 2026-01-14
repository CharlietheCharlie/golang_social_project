package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"

	"github.com/joho/godotenv"
)

const version = "0.0.1" 

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	// Initialize a new config struct
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			dsn: env.GetString("DB_DSN", "postgres://user:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	// Initialize the database connection pool
	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("database connection pool established")

	// Initialize the storage layer
	store := store.NewStorage(db)


	// Initialize a new application struct
	app := &application{
		config: cfg,
		store:  store,
	}

 
	// Mount the routes
	mux := app.mount()

	// Start the server
	log.Fatal(app.run(mux))
}