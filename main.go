package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("GET /characters", getCharacters)
	http.HandleFunc("GET /characters/{id}", getCharacter)

	println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}