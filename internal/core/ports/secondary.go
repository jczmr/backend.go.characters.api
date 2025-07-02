package ports

import "backend.go.characters.api/internal/core/domain"

type CharacterRepository interface {
	SaveCharacter(character *domain.Character) error
	FindCharacterByName(name string) (*domain.Character, error)
}

type DragonBallAPIClient interface {
	FindCharacterByName(name string) (*domain.Character, error)
	FindCharacterByID(id string) (*domain.Character, error)
}
