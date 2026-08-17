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

// The Character Struct
type Character struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Alignment string `json:"alignment"`
	Image string `json:"image"`
}

// The expected response
type CharacterResponse struct {
	Data []Character `json:"data"`
	Total int `json:"total"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
}

// Validate if character exists, if not throw an error.
// This is just a helper function
func validateCharacter(character Character) error {
	if strings.TrimSpace(character.Name) == "" {
		return errors.New("name is required.") // Error if there is no provided name
	}

	if strings.TrimSpace(character.Description) == "" {
		return errors.New("description is required") // Error if there is no provided description
	}

	if strings.TrimSpace(character.Alignment) == "" {
		return errors.New("alignment is required") // Error if there is no provided alignment
	}

	if strings.TrimSpace(character.Image) == "" {
		return errors.New("image is required")
	}

	return nil // Nil, because there are no further errors.
}

// Get all of the characters in the database.
func getCharacters(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query() // Get the query

	// Get the alignment and name
	alignmentFilter := query.Get("alignment")
	nameFilter := query.Get("name")

	// Limit and Offset, for searching.
	limit := 0
	offset := 0

	// SortBy and Order are more optional search params.
	sortBy := query.Get("sort")
	order := query.Get("order")

	// Check if the user has provided SortBy.
	if sortBy != "" {
		switch sortBy {
		case "id", "name", "alignment":
			//Valid

		default:
			http.Error(w, "Invalid sort field", http.StatusBadRequest) // Error if sort field is invalid.
			return
		}
	}

	// Check if user provided order.
	if order != "" {
		switch order {
		case "asc", "desc":
			// Valid

		default:
			http.Error(w, "Invalid sort order", http.StatusBadRequest) // Error if sort order is invalid.
			return
		}
	}

	// Check if order is not provided but sortBy is, because order needs sortBy.
	if order != "" && sortBy == "" {
		http.Error(w, "Order requires a sort field", http.StatusBadRequest) // Error if there is order, but not sortBy.
		return
	}

	var results []Character // A slice of characters

	// Check for a limit
	if value := query.Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value);

		if err != nil || parsed < 0 {
			http.Error(w, "Invalid limit", http.StatusBadRequest) // Error if limit is invalid.
			return
		}

		limit = parsed
	}

	// Check for offset
	if value := query.Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)

		if err != nil || parsed < 0 {
			http.Error(w, "Invalid Offset", http.StatusBadRequest) // Error if offset is invalid.
			return
		}

		offset = parsed
	}

	// Query the Database
	rows, err := db.Query(`
		SELECT id, name, description, alignment, image
		FROM characters
	`)

	if err != nil {
		http.Error(w, "Failed to fetch characters", http.StatusInternalServerError) // Error if characters could not be fetched.
		return
	} 

	defer rows.Close()

	for rows.Next() {
		var character Character

		// Scan the rows
		err := rows.Scan(
			&character.ID,
			&character.Name,
			&character.Description,
			&character.Alignment,
			&character.Image,
		)

		// Throw error if the Character Read failed.
		if err != nil {
			http.Error(w, "Failed to read character", http.StatusInternalServerError)
			return
		}

		// Check if the Alignment matches the provided alignment filter.
		if alignmentFilter != "" && !strings.EqualFold(character.Alignment, alignmentFilter) {
			continue
		}

		// Check if the Name matches the provided name filter.
		if nameFilter != "" && strings.EqualFold(character.Name, nameFilter) {
			continue
		}

		// Append the results to the character variable.
		results = append(results, character)
	}

	// Throw error if it failed to read characters.
	if err := rows.Err(); err != nil {
		http.Error(w, "Failed to read characters", http.StatusInternalServerError)
		return
	}

	// Check the sort order (descending/ascending)
	descending := order == "desc"

	if sortBy != "" {
		sort.Slice(results, func(i, j int) bool {
			switch sortBy { // Get all the cases of what the user can sort by.
			case "id": // If sorting by id..
				if descending { // Check descending
					return results[i].ID > results[j].ID
				}
				return results[i].ID < results[j].ID

			case "name": // If sorting by name..
				left := strings.ToLower(results[i].Name)
				right := strings.ToLower(results[j].Name)

				if descending { // Check descending
					return left > right
				}

				return left < right

			case "alignment": // If sorting by alignment..
				left := strings.ToLower(results[i].Alignment)
				right := strings.ToLower(results[j].Alignment)

				if descending { // Check descending
					return left > right
				}
				return left < right
			}

			return false // Return false if nothing matches
		})
	}

	total := len(results) // Get the total number of characters that matched the search.

	if offset >= len(results) { // Handle the offset
		results = []Character{}
	} else {
		results = results[offset:]
	}

	if limit > 0 && limit < len(results) { // Handle the limit
		results = results[:limit]
	}

	// Build the character response
	response := CharacterResponse {
		Data: results,
		Total: total,
		Limit: limit,
		Offset: offset,
	}

	w.Header().Set("Content-Type", "application/json") // Set the header
	json.NewEncoder(w).Encode(response) // Encode the response
}

// Get a single character
func getCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")); // Get the id
	if err != nil { // Check for errors
		http.Error(w, "Invalid character ID", http.StatusBadRequest) // Throw an error if the ID is invalid.
		return
	}

	var character Character

	// Query the database
	err = db.QueryRow(`
		SELECT id, name, description, alignment, image
		FROM characters
		WHERE id = ?
	`, id).Scan(
		&character.ID,
		&character.Name,
		&character.Description,
		&character.Alignment,
		&character.Image,
	)

	if err == sql.ErrNoRows { // Error if there are no rows that were found, meaning the character was not found.
		http.Error(w, "Character not found", http.StatusNotFound) // Throw an error for character was not found.
		return
	}

	if err != nil { // Check for a fetch error
		http.Error(w, "Failed to fetch character", http.StatusInternalServerError) // Throw a fetch error.
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the header
	json.NewEncoder(w).Encode(character) // Encode the response
}

// create a character
func createCharacter(w http.ResponseWriter, r *http.Request) {
	var character Character 

	err := json.NewDecoder(r.Body).Decode(&character) // New decoder

	if err != nil { // Check for invalid JSON
		http.Error(w, "Invalid JSON", http.StatusBadRequest) // Error for invalid JSON
		return
	}

	err = validateCharacter(character) // Validate the Character

	if err != nil { // If we found an error
		http.Error(w, err.Error(), http.StatusBadRequest) // Throw an error for a bad request.
		return
	}

	// Query the database to insert a new character
	result, err := db.Exec(`
		INSERT INTO characters (name, description, alignment, image)
		VALUES (?, ?, ?, ?)
	`,
		character.Name,
		character.Description,
		character.Alignment,
		character.Image,
	)

	if err != nil { // Check for an error with the database
		http.Error(w, "Failed to create character", http.StatusInternalServerError) // Throw an error if we failed to create the character
		return
	}

	id, err := result.LastInsertId() // Get the Last id

	if err != nil { // Check for errors
		http.Error(w, "Failed to get character ID", http.StatusInternalServerError) // Throw an error if it failed to get the last id.
		return
	}

	character.ID = int(id) // Give the character it's id

	w.Header().Set("Content-Type", "application/json") // Set the header
	w.WriteHeader(http.StatusCreated) // Write the header

	json.NewEncoder(w).Encode(character) // Encode character
}

// Delete a character
func deleteCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")) // Get the ID provided

	if err != nil { // If there is an error
		http.Error(w, "Invalid character ID", http.StatusBadRequest) // Error if the provided id is invalid
		return
	}

	// Query the database
	result, err := db.Exec(`
		DELETE FROM characters
		WHERE id = ?
	`, id)

	if err != nil { // If there was an error
		http.Error(w, "Failed to delete character", http.StatusInternalServerError) // Throw an error because it failed to delete the character
		return
	}

	rowsAffected, err := result.RowsAffected() // Get the affected rows

	if err != nil { // Check if there was an error while getting the affected rows
		http.Error(w, "Failed to check deletion", http.StatusInternalServerError) // Throw the error
		return
	}

	if rowsAffected == 0 { // If there was no rows that were affected, it means that the character could not be found.
		http.Error(w, "Character not found", http.StatusNotFound) // Therefore, we throw an error, because of course we do. That's all these frickin APIs do 😭
		return
	}

	w.WriteHeader(http.StatusNoContent) // Write the header
}

// Update an existing Character
func updateCharacter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")) // Get the provided ID

	if err != nil { // If the id is invalid..
		http.Error(w, "Invalid character ID", http.StatusBadRequest) // Throw an invalid ID error
		return
	}

	var updatedCharacter Character

	err = json.NewDecoder(r.Body).Decode(&updatedCharacter) // New Decoder

	if err != nil { // If there was an error decoding
		http.Error(w, "Invalid JSON", http.StatusBadRequest) // Tell the user that their JSON sucks! (jk LOL)
		return
	}

	err = validateCharacter(updatedCharacter) // Validate it

	if err != nil { // If there is a validation error
		http.Error(w, err.Error(), http.StatusBadRequest) // Throw it
		return
	}

	// Query the database
	result, err := db.Exec(`
		UPDATE characters
		SET name = ?, description = ?, alignment = ?, image = ?
		WHERE id = ?
	`,
		updatedCharacter.Name,
		updatedCharacter.Description,
		updatedCharacter.Alignment,
		updatedCharacter.Image,
		id,
	)

	if err != nil { // If the query caused an error,
		http.Error(w, "Failed to check update", http.StatusInternalServerError) // Throw it
		return
	}

	rowsAffected, err := result.RowsAffected() // Get the affected rows

	if err != nil { // If there was an error getting the affected rows
		http.Error(w, "Failed to check update", http.StatusInternalServerError) // Throw it
		return
	}

	if rowsAffected == 0 { // If there was no affected rows (character not found)
		http.Error(w, "Character not found", http.StatusNotFound) // Throw an error to complain about it
		return
	}

	updatedCharacter.ID = id // Set the id

	w.Header().Set("Content-Type", "application/json") // Set the header
	json.NewEncoder(w).Encode(updatedCharacter) // Encode
}

// My comments got A LOT MORE UNCANNY as the file went on LOL