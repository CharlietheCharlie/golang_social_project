package main

import (
	"net/http"
	"social/internal/store"
)


type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate the payload
	if err := Validate.Struct(&payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	
	comment := &store.Comment{
		PostID:  post.ID,
		UserID:  1, // TODO: Replace with the authenticated user's ID
		Content: payload.Content,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}