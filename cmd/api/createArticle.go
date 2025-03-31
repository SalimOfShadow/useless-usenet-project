package main

import (
	"net/http"

	"github.com/salimofshadow/usenet-client/internal/routes/store"
)


type CreateArticlePayload struct {
	Subject   string   `json:"subject" validate:"required,max=255"`
	Author    string   `json:"author" validate:"required"`
	Newsgroup string   `json:"newsgroup" validate:"required"`
	Body      string   `json:"body" validate:"required,max=10000"`
	Tags      []string `json:"tags,omitempty"`
}


func (app *application) createArticleHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateArticlePayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// user := getUserFromContext(r)

	article := &store.Article{
		ID: 0, // FIXME: Should have retrieved the user's ID from the context with a middleware by now
		Author:   payload.Author,
		Newsgroup: payload.Newsgroup,
		Body:  payload.Body,
		Tags:    payload.Tags,
	
	}

	ctx := r.Context()

	if err := app.store.Articles.Create(ctx, article); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, article); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
