package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.POST("/calculate", app.recoverPanic(app.validateJSON(app.calculateHandler)))

	return router
}
