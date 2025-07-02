package dragonballapi_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend.go.characters.api/internal/adapters/secondary/dragonballapi"

	"github.com/stretchr/testify/assert"
)

func TestDragonBallAPIClientFindCharacterByName(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/characters", r.URL.Path)
		if r.URL.Path == "/api/characters" {
			response := map[string][]map[string]string{
				"items": {
					{"id": "1", "name": "Goku", "ki": "9000", "race": "Saiyan"},
					{"id": "2", "name": "Vegeta", "ki": "8000", "race": "Saiyan"},
				},
			}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Override base URL for testing
	//originalBaseURL := dragonballapi.BaseURL // Assuming you made BaseURL public in dragonballapi/client.go for testing
	originalBaseURL := "https://dragonball-api.com/api"
	dragonballapi.BaseURL = server.URL + "/api"
	defer func() { dragonballapi.BaseURL = originalBaseURL }()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := dragonballapi.NewDragonBallAPIClient(logger)

	// Test case: Character found
	character, err := client.FindCharacterByName("Goku")
	assert.NoError(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, "Goku", character.Name)
	assert.Equal(t, "1", character.ID)
	assert.Equal(t, "Saiyan", character.Race)

	// Test case: Character not found
	character, err = client.FindCharacterByName("Frieza")
	assert.NoError(t, err) // No error, just character is nil
	assert.Nil(t, character)
}

func TestDragonBallAPIClientFindCharacterByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/characters/1" {
			response := map[string]string{
				"id":   "1",
				"name": "Gohan",
				"ki":   "7000",
				"race": "Saiyan",
			}
			json.NewEncoder(w).Encode(response)
		} else if r.URL.Path == "/api/characters/999" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	originalBaseURL := dragonballapi.BaseURL
	dragonballapi.BaseURL = server.URL + "/api"
	defer func() { dragonballapi.BaseURL = originalBaseURL }()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := dragonballapi.NewDragonBallAPIClient(logger)

	// Test case: Character found by ID
	character, err := client.FindCharacterByID("1")
	assert.NoError(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, "Gohan", character.Name)
	assert.Equal(t, "1", character.ID)

	// Test case: Character not found by ID
	character, err = client.FindCharacterByID("999")
	assert.NoError(t, err)
	assert.Nil(t, character)

	// Test case: API error
	_, err = client.FindCharacterByID("invalid_id") // This will hit the /api/characters/invalid_id route, causing 500
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Dragon Ball API returned status 500")
}
