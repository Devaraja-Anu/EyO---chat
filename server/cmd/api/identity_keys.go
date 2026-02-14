package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devaraja-anu/eyo/server/internals/data"
)

type IdentityKeyRequest struct {
	Username  string `json:"username"`
	PublicKey string `json:"public_key"`
}

func (app *application) publishIdentityKeyHandler(w http.ResponseWriter, r *http.Request) {

	var input IdentityKeyRequest

	err := app.ReadJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, err.Error())
		return
	}

	decoded, err := app.decodePublicKey(input.PublicKey)

	if err != nil {
		app.badRequestResponse(w, "invalid public key  encoding")
		return
	}

	if len(decoded) != 32 {
		app.badRequestResponse(w, "invalid identity key length")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := app.models.Users.GetByUsername(ctx, input.Username)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, err)
		}

		return
	}

	idKey := &data.IdentityKey{
		UserID:    user.ID,
		PublicKey: decoded,
	}

	err = app.models.IdentityKey.Insert(ctx, idKey)

	if err != nil {
		switch {

		case errors.Is(err, data.ErrIdentityKeyExists):
			app.conflictResponse(w, "identity key already published")
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	app.logger.Info("identity key established", map[string]any{
		"user_id": user.ID,
	})

	app.writeJSON(w, http.StatusCreated, envelope{"message": "identity key stored"}, nil)

}
