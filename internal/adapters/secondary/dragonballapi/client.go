package dragonballapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"log/slog"

	"backend.go.characters.api/internal/core/domain"
)

var BaseURL string = "https://dragonball-api.com/api"

type apiCharacter struct {
	ID   json.Number `json:"id"`
	Name string      `json:"name"`
	Ki   string      `json:"ki"`
	Race string      `json:"race"`
}

type apiCharactersResponse struct {
	Items []apiCharacter `json:"items"`
}

type dragonBallAPIClient struct {
	httpClient *http.Client
	logger     *slog.Logger
}

func NewDragonBallAPIClient(logger *slog.Logger) *dragonBallAPIClient {
	return &dragonBallAPIClient{
		httpClient: &http.Client{},
		logger:     logger,
	}
}

func (c *dragonBallAPIClient) FindCharacterByName(name string) (*domain.Character, error) {
	// The API does not directly support lookup by name.
	// We need to fetch all characters and then filter. This is inefficient but dictated by the API.
	c.logger.Info("Fetching all characters from external API to find by name", slog.String("target_name", name))

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/characters", BaseURL))
	if err != nil {
		c.logger.Error("Failed to make request to Dragon Ball API", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Error("Dragon Ball API returned non-OK status", slog.Int("status_code", resp.StatusCode), slog.String("response_body", string(bodyBytes)))
		return nil, fmt.Errorf("the Dragon Ball API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse apiCharactersResponse
	// Enable decoding numbers into json.Number
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber() // Crucial for json.Number to work
	if err := decoder.Decode(&apiResponse); err != nil {
		c.logger.Error("Failed to decode Dragon Ball API response", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	for _, apiChar := range apiResponse.Items {
		if apiChar.Name == name {
			c.logger.Info("Character found in external API by name", slog.String("character_name", name), slog.String("character_id", apiChar.ID.String())) // Convert json.Number to string
			return &domain.Character{
				ID:   apiChar.ID.String(),
				Name: apiChar.Name,
				Ki:   apiChar.Ki,
				Race: apiChar.Race,
			}, nil
		}
	}

	c.logger.Info("Character not found in external API by name", slog.String("character_name", name))
	return nil, nil // Character not found
}

func (c *dragonBallAPIClient) FindCharacterByID(id string) (*domain.Character, error) {
	c.logger.Info("Fetching character by ID from external API", slog.String("character_id", id))

	resp, err := c.httpClient.Get(fmt.Sprintf("%s/characters/%s", BaseURL, id))
	if err != nil {
		c.logger.Error("Failed to make request to Dragon Ball API", slog.String("error", err.Error()), slog.String("character_id", id))
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.logger.Warn("Character not found in external API by ID", slog.String("character_id", id))
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Error("Dragon Ball API returned non-OK status for ID lookup", slog.Int("status_code", resp.StatusCode), slog.String("response_body", string(bodyBytes)), slog.String("character_id", id))
		return nil, fmt.Errorf("the Dragon Ball API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiChar apiCharacter
	// Enable decoding numbers into json.Number
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber() // Crucial for json.Number to work
	if err := decoder.Decode(&apiChar); err != nil {
		c.logger.Error("Failed to decode Dragon Ball API response for ID lookup", slog.String("error", err.Error()), slog.String("character_id", id))
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	c.logger.Info("Character found in external API by ID", slog.String("character_id", apiChar.ID.String()), slog.String("character_name", apiChar.Name))
	return &domain.Character{
		ID:   apiChar.ID.String(), // Convert json.Number to string
		Name: apiChar.Name,
		Ki:   apiChar.Ki,
		Race: apiChar.Race,
	}, nil
}
