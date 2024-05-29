package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) recoverPanic(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next(w, r, ps)
	}
}

func (app *application) validateJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.A < 0 || req.B < 0 {
			app.badRequestResponse(w, r, err)
			return
		}
		next(w, r, ps)
	}
}
