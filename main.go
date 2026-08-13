package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("GET /characters", getCharacters)
	http.HandleFunc("GET /characters/{id}", getCharacter)
	http.HandleFunc("POST /characters", createCharacter)
	http.HandleFunc("DELETE /characters/{id}", deleteCharacter)
	http.HandleFunc("PUT /characters/{id}", updateCharacter)

	println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}