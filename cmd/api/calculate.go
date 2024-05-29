package main

import (
	"encoding/json"
	"fmt"
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
	var input struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.A < 0 || input.B < 0 {
		app.badRequestResponse(w, r, fmt.Errorf("a and b must be non-negative"))
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
	json.NewEncoder(w).Encode(response)
}
