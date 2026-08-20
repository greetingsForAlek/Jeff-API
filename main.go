package main

import (
	"log"
	"net/http"
	"os"
)

func enableCORS(next http.Handler) http.Handler { // enable CORS
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r);
	})
}

// The MAIN function
func main() {
	err := initDB() // Initialise the database

	if err != nil { // If there was an error
		log.Fatal(err) // Throw a FATAL ERROR 💀
	}

	http.HandleFunc("GET /characters", getCharacters) // getCharacters route
	http.HandleFunc("GET /characters/{id}", getCharacter) // getCharacter route
	http.HandleFunc("POST /characters", createCharacter) // createCharacter route
	http.HandleFunc("DELETE /characters/{id}", deleteCharacter) // deleteCharacter route
	http.HandleFunc("PUT /characters/{id}", updateCharacter) // updateCharacter route

	println("Server running on http://localhost:8080") // Tell the user that the server is running

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux)) // Listen and serve!
}