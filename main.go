package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/characters", getCharacters)

	println("Server running on http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}