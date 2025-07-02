package services

import (
	"fmt"

	"log/slog"

	"backend.go.characters.api/internal/core/domain"
	"backend.go.characters.api/internal/core/ports"
)

type characterService struct {
	characterRepository ports.CharacterRepository
	dragonBallAPIClient ports.DragonBallAPIClient
	logger              *slog.Logger
}

func NewCharacterService(
	characterRepository ports.CharacterRepository,
	dragonBallAPIClient ports.DragonBallAPIClient,
	logger *slog.Logger,
) ports.CharacterService {
	return &characterService{
		characterRepository: characterRepository,
		dragonBallAPIClient: dragonBallAPIClient,
		logger:              logger,
	}
}

func (s *characterService) CreateCharacter(characterName string) (*domain.Character, error) {
	s.logger.Info("Attempting to create or retrieve character", slog.String("character_name", characterName))

	// 1. Check if character exists in local database
	existingCharacter, err := s.characterRepository.FindCharacterByName(characterName)
	if err == nil && existingCharacter != nil {
		s.logger.Info("Character found in local database", slog.String("character_name", characterName), slog.String("character_id", existingCharacter.ID))
		return existingCharacter, nil
	}

	// 2. If not found, fetch from external API
	s.logger.Info("Character not found in local database, fetching from external API", slog.String("character_name", characterName))
	apiCharacter, err := s.dragonBallAPIClient.FindCharacterByName(characterName)
	if err != nil {
		s.logger.Error("Failed to fetch character from external API", slog.String("error", err.Error()), slog.String("character_name", characterName))
		return nil, fmt.Errorf("failed to fetch character from external API: %w", err)
	}
	if apiCharacter == nil {
		s.logger.Warn("Character not found in external API", slog.String("character_name", characterName))
		return nil, fmt.Errorf("character '%s' not found in external API", characterName)
	}

	// 3. Populate additional fields and save to database
	newCharacter := &domain.Character{
		ID:   apiCharacter.ID,
		Name: apiCharacter.Name,
		Ki:   apiCharacter.Ki,
		Race: apiCharacter.Race,
	}

	if err := s.characterRepository.SaveCharacter(newCharacter); err != nil {
		s.logger.Error("Failed to save character to database", slog.String("error", err.Error()), slog.String("character_name", newCharacter.Name))
		return nil, fmt.Errorf("failed to save character: %w", err)
	}

	s.logger.Info("Successfully fetched and saved character", slog.String("character_name", newCharacter.Name), slog.String("character_id", newCharacter.ID))
	return newCharacter, nil
}
