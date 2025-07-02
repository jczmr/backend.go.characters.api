package ports

import "backend.go.characters.api/internal/core/domain"

type CharacterService interface {
	CreateCharacter(characterName string) (*domain.Character, error)
}
