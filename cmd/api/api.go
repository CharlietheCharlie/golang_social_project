package main

import (
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	store  store.Storage
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	// Create a new Chi router
	r := chi.NewRouter()
	// Add some basic middleware
	r.Use(middleware.RequestID) // assign a unique ID to each request
	r.Use(middleware.RealIP)    // get the real IP from X-Forwarded-For header
	r.Use(middleware.Logger)    // log info about the request
	r.Use(middleware.Recoverer) // recover from panics without crashing server

	r.Use(middleware.Timeout(60 * time.Second)) // set a timeout for requests

	r.Route("/v1", func(r chi.Router) { // use Router to group routes
		r.Get("/health", app.healthCheckHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware) // middleware to load post from database

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware) // middleware to load user from database

				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Delete("/follow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserfeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}
