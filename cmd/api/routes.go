package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.POST("/calculate", app.validateJSON(app.calculateHandler))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		router.ServeHTTP(w, r)
	})
}
