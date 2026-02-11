package main

import (
	"fmt"
	"net/http"
)


func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("status: available")
}
