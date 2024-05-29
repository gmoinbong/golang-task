package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Request struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Response struct {
	FactorialA string `json:"a_factorial"`
	FactorialB string `json:"b_factorial"`
}

var mu sync.Mutex

func (app *application) calculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input Request

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	factorialA := factorial(input.A)
	factorialB := factorial(input.B)

	response := Response{
		FactorialA: factorialA.String(),
		FactorialB: factorialB.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding response:", err)
	}
}
