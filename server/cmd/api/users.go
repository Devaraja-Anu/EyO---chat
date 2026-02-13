package main

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/devaraja-anu/eyo/server/internals/data"
)

type registerUserRequest struct  {
	Username string `json:"username"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input registerUserRequest

	err := app.ReadJSON(w,r,&input)
	
	if err != nil {
		app.logger.Error("invalid register request",map[string]any{
			"error":err.Error(),
			"ip":r.RemoteAddr,
		})
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}

	input.Username = strings.TrimSpace(input.Username)

	if input.Username == ""  {
		http.Error(w,"username must not be empty",http.StatusBadRequest)
		return
	}

	if len(input.Username) > 32 {
		http.Error(w,"username is too long",http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	user := &data.User{Username: input.Username}

	err = app.models.Users.Insert(ctx,user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			app.logger.Info("duplicate username attempt",map[string]any{
				"username":input.Username,
				"ip":r.RemoteAddr,
			})
			http.Error(w,"username already taken",http.StatusConflict)
			return
		default:
			app.logger.Info("User Insert failed",map[string]any{
				"username":input.Username,
				"ip":r.RemoteAddr,
			})
			http.Error(w,"Server Error",http.StatusInternalServerError)
			return
		}
	}

	app.logger.Info("user registered",map[string]any{
		"user_ID": user.ID,
		"username": user.Username,
	})

	err  = app.writeJSON(w, http.StatusCreated,envelope{"user":user},nil)

	if err != nil {
		app.logger.Error("response write failed",map[string]any{"error":err.Error()})
	}


}