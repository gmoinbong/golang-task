package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	var req Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Printf("Received body: %s\n", string(body))

	err = json.Unmarshal(body, &req)
	if err != nil || req.A < 0 || req.B < 0 {
		app.badRequestResponse(w, r, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	factorialA := factorial(req.A)
	factorialB := factorial(req.B)

	response := Response{
		FactorialA: factorialA.String(),
		FactorialB: factorialB.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
