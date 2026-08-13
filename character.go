package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Character struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Alignment string `json:"alignment"`
}

var characters = []Character{
	{
		ID: 1,
		Name: "Jeff",
		Description: "The main character",
		Alignment: "Good",
	},
	{
		ID: 2,
		Name: "Mr. Paper",
		Description: "The main villain",
		Alignment: "Evil",
	},
}

func getNextId() int {
	nextID := 1

	for _, character := range characters {
		if character.ID >= nextID {
			nextID = character.ID + 1
		}
	}

	return nextID
}

func getCharacters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(characters)
}

func getCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"));
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	for _, character := range characters {
		if character.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(character)
			return
		}
	}

	http.Error(w, "Character not found", http.StatusNotFound)
}

func createCharacter(w http.ResponseWriter, r *http.Request) {
	var character Character

	err := json.NewDecoder(r.Body).Decode(&character)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	character.ID = getNextId()

	characters = append(characters, character)

	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(character)
}

func deleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	for i, character := range characters {
		if character.ID == id {
			characters = append(characters[:i], characters[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Character not found", http.StatusNotFound)
}

func updateCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	var updatedCharacter Character

	err = json.NewDecoder(r.Body).Decode(&updatedCharacter)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for i, character := range characters {
		if character.ID == id {
			updatedCharacter.ID = id
			characters[i] = updatedCharacter

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCharacter)
			return
		}
	}

	http.Error(w, "Character not found", http.StatusNotFound)
}