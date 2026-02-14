package main

import "net/http"

func (app *application) routes() http.Handler{
	
	r := http.NewServeMux()
	r.HandleFunc("GET /v1/healthcheck", app.healthCheckHandler)
	r.HandleFunc("POST /v1/users/register",app.registerUserHandler)
	r.HandleFunc("POST /v1/keys/identity",app.publishIdentityKeyHandler)
	return  r
	
}