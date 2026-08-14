package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Character struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Alignment string `json:"alignment"`
}

type CharacterResponse struct {
	Data []Character `json:"data"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
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
		Name: "Mr Paper",
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

func validateCharacter(character Character) error {
	if strings.TrimSpace(character.Name) == "" {
		return errors.New("name is required.")
	}

	if strings.TrimSpace(character.Description) == "" {
		return errors.New("description is required")
	}

	if strings.TrimSpace(character.Alignment) == "" {
		return errors.New("type is required")
	}

	return nil
}

func getCharacters(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	alignmentFilter := query.Get("alignment")
	nameFilter := query.Get("name")

	limit := 0
	offset := 0

	sortBy := query.Get("sort")
	order := query.Get("order")

	if sortBy != "" {
		switch sortBy {
		case "id", "name", "alignment":
			//Valid

		default:
			http.Error(w, "Invalid sort field", http.StatusBadRequest)
			return
		}
	}

	if order != "" {
		switch order {
		case "asc", "desc":
			// Valid

		default:
			http.Error(w, "Invalid sort order", http.StatusBadRequest)
			return
		}
	}

	if order != "" && sortBy == "" {
		http.Error(w, "Order requires a sort field", http.StatusBadRequest)
		return
	}

	var results []Character

	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value);

		if err != nil || parsed < 0 {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		limit = parsed
	}

	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 0 {
			http.Error(w, "Invalid Offset", http.StatusBadRequest)
			return
		}

		offset = parsed
	}

	for _, character := range characters {
		if alignmentFilter != "" && !strings.EqualFold(character.Name, alignmentFilter) {
			continue
		}

		if nameFilter != "" && !strings.EqualFold(character.Name, nameFilter) {
			continue
		}

		results = append(results, character)
	}

	descending := order == "desc"

	if sortBy != "" {
		sort.Slice(results, func(i, j int) bool {
			switch sortBy {
			case "id":
				if descending {
					return results[i].ID > results[j].ID
				}
				return results[i].ID < results[j].ID

			case "name":
				left := strings.ToLower(results[i].Name)
				right := strings.ToLower(results[j].Name)

				if descending {
					return left > right
				}

				return left < right

			case "alignment":
				left := strings.ToLower(results[i].Alignment)
				right := strings.ToLower(results[j].Alignment)

				if descending {
					return left > right
				}
				return left < right
			}

			return false
		})
	}

	if offset >= len(results) {
		results = []Character{}
	} else {
		results = results[offset:]
	}

	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	response := CharacterResponse {
		Data: results,
		Total: len(results),
		Limit: limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	err = validateCharacter(character)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	err = validateCharacter(updatedCharacter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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