package main

import (
	"log"
	"social/internal/auth"
	"social/internal/db"
	"social/internal/env"
	"social/internal/mailer"
	"social/internal/store"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Social API
//	@description	API for Social.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	// Initialize a new config struct
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", "postgres://user:password@localhost:5432/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("AUTH_BASIC_USERNAME", "admin"),
				password: env.GetString("AUTH_BASIC_PASSWORD", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "secret"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "social",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Initialize the database connection pool
	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Initialize the storage layer
	store := store.NewStorage(db)

	// Initialize the mailer
	// mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize the authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss, //aud
		cfg.auth.token.iss, //iss
	)
	// Initialize a new application struct
	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
	}

	// Mount the routes
	mux := app.mount()

	// Start the server
	logger.Fatal(app.run(mux))
}
