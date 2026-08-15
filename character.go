package main

import (
	"database/sql"
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

func validateCharacter(character Character) error {
	if strings.TrimSpace(character.Name) == "" {
		return errors.New("name is required.")
	}

	if strings.TrimSpace(character.Description) == "" {
		return errors.New("description is required")
	}

	if strings.TrimSpace(character.Alignment) == "" {
		return errors.New("alignment is required")
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

	rows, err := db.Query(`
		SELECT id, name, description, alignment
		FROM characters
	`)

	if err != nil {
		http.Error(w, "Failed to fetch characters", http.StatusInternalServerError)
		return
	} 

	defer rows.Close()

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Description,
			&character.Alignment,
		)

		if err != nil {
			http.Error(w, "Failed to read character", http.StatusInternalServerError)
			return
		}

		if alignmentFilter != "" && !strings.EqualFold(character.Alignment, alignmentFilter) {
			continue
		}

		if nameFilter != "" && strings.EqualFold(character.Name, nameFilter) {
			continue
		}

		results = append(results, character)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read characters", http.StatusInternalServerError)
		return
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

	total := len(results)

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
		Total: total,
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

	var character Character

	err = db.QueryRow(`
		SELECT id, name, description, alignment
		FROM characters
		WHERE id = ?
	`, id).Scan(
		&character.ID,
		&character.Name,
		&character.Description,
		&character.Alignment,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Failed to fetch character", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
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

	result, err := db.Exec(`
		INSERT INTO characters (name, description, alignment)
		VALUES (?, ?, ?)
	`,
		character.Name,
		character.Description,
		character.Alignment,
	)

	if err != nil {
		http.Error(w, "Failed to create character", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		http.Error(w, "Failed to get character ID", http.StatusInternalServerError)
		return
	}

	character.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(character)
}

func deleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		DELETE FROM characters
		WHERE id = ?
	`, id)

	if err != nil {
		http.Error(w, "Failed to delete character", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		http.Error(w, "Failed to check deletion", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	result, err := db.Exec(`
		UPDATE characters
		SET name = ?, description = ?, alignment = ?
		WHERE id = ?
	`,
		updatedCharacter.Name,
		updatedCharacter.Description,
		updatedCharacter.Alignment,
		id,
	)

	if err != nil {
		http.Error(w, "Failed to check update", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		http.Error(w, "Failed to check update", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	updatedCharacter.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCharacter)
}