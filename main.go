package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
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

func validateJSON(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var req Request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.A < 0 || req.B < 0 {
			http.Error(w, `{"error":"Incorrect input"}`, http.StatusBadRequest)
			return
		}
		next(w, r, ps)
	}
}

func factorial(n int) *big.Int {
	if n == 0 {
		return big.NewInt(1)
	}
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

func calculateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req Request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Incorrect input"}`, http.StatusBadRequest)
		return
	}
	fmt.Printf("Received body: %s\n", string(body))

	err = json.Unmarshal(body, &req)
	if err != nil || req.A < 0 || req.B < 0 {
		http.Error(w, `{"error":"Incorrect input"}`, http.StatusBadRequest)
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

func main() {
	router := httprouter.New()
	router.POST("/calculate", validateJSON(calculateHandler))

	fmt.Println("Server running at http://localhost:8989")
	log.Fatal(http.ListenAndServe(":8989", router))
}
