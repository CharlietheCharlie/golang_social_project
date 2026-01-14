package main

import (
	"net/http"
	"social/internal/store"
)


func (app *application) getUserfeedHandler(w http.ResponseWriter, r *http.Request) {
	// default values
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	// parse query
	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(1), fq) // TODO: replace 1 with authenticated user ID
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}