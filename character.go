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