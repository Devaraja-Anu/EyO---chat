package main

import "net/http"

func (app *application) routes() http.Handler{
	
	r := http.NewServeMux()
	r.HandleFunc("/v1/healthcheck", app.healthCheckHandler)

	return  r
	
}