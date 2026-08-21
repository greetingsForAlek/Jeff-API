package main

import (
	"encoding/json"
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

// Function to verify the admin password that allows for admin access to the admin website
func verifyPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Password string
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	correctPassword := os.Getenv("ADMIN_PASSWORD")

	if data.Password != correctPassword {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true}`));
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
	http.HandleFunc("POST /verifyPassword", verifyPassword) // verifyPassword route

	println("Server running on http://localhost:8080") // Tell the user that the server is running

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux)) // Listen and serve!
}