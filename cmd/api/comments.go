package main

import (
	"net/http"
	"social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
}

// createComment godoc
//
//	@Summary		Create comment
//	@Description	Create comment
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		int						true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/posts/{post_id}/comments [post]
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
