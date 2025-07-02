package services_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"backend.go.characters.api/internal/core/domain"
	"backend.go.characters.api/internal/core/services"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

// Mock for CharacterRepository
type MockCharacterRepository struct {
	mock.Mock
}

func (m *MockCharacterRepository) SaveCharacter(character *domain.Character) error {
	args := m.Called(character)
	return args.Error(0)
}

func (m *MockCharacterRepository) FindCharacterByName(name string) (*domain.Character, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Character), args.Error(1)
}

// Mock for DragonBallAPIClient
type MockDragonBallAPIClient struct {
	mock.Mock
}

func (m *MockDragonBallAPIClient) FindCharacterByName(name string) (*domain.Character, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Character), args.Error(1)
}

func (m *MockDragonBallAPIClient) FindCharacterByID(id string) (*domain.Character, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Character), args.Error(1)
}

func TestCharacterService_CreateCharacter_FromDB(t *testing.T) {
	mockRepo := new(MockCharacterRepository)
	mockAPIClient := new(MockDragonBallAPIClient)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	charService := services.NewCharacterService(mockRepo, mockAPIClient, logger)

	expectedCharacter := &domain.Character{
		ID:   "123",
		Name: "Goku",
	}

	// Expect FindCharacterByName to return an existing character
	mockRepo.On("FindCharacterByName", "Goku").Return(expectedCharacter, nil).Once()

	character, err := charService.CreateCharacter("Goku")
	assert.NoError(t, err)
	assert.Equal(t, expectedCharacter, character)
	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertNotCalled(t, "FindCharacterByName") // Should not call API if found in DB
}

func TestCharacterService_CreateCharacter_FromAPI(t *testing.T) {
	mockRepo := new(MockCharacterRepository)
	mockAPIClient := new(MockDragonBallAPIClient)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	charService := services.NewCharacterService(mockRepo, mockAPIClient, logger)

	apiCharacter := &domain.Character{
		ID:   "456",
		Name: "Vegeta",
		Ki:   "8000",
		Race: "Saiyan",
	}

	// Expect FindCharacterByName from DB to return nil (not found)
	mockRepo.On("FindCharacterByName", "Vegeta").Return(nil, nil).Once()
	// Expect FindCharacterByName from API to return a character
	mockAPIClient.On("FindCharacterByName", "Vegeta").Return(apiCharacter, nil).Once()
	// Expect SaveCharacter to be called
	mockRepo.On("SaveCharacter", mock.AnythingOfType("*domain.Character")).Return(nil).Once()

	character, err := charService.CreateCharacter("Vegeta")
	assert.NoError(t, err)
	assert.NotNil(t, character)
	assert.Equal(t, apiCharacter.ID, character.ID)
	assert.Equal(t, apiCharacter.Name, character.Name)
	assert.Equal(t, apiCharacter.Ki, character.Ki)
	assert.Equal(t, apiCharacter.Race, character.Race)

	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertExpectations(t)
}

func TestCharacterService_CreateCharacter_APIError(t *testing.T) {
	mockRepo := new(MockCharacterRepository)
	mockAPIClient := new(MockDragonBallAPIClient)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	charService := services.NewCharacterService(mockRepo, mockAPIClient, logger)

	// Expect FindCharacterByName from DB to return nil (not found)
	mockRepo.On("FindCharacterByName", "Krillin").Return(nil, nil).Once()
	// Expect FindCharacterByName from API to return an error
	mockAPIClient.On("FindCharacterByName", "Krillin").Return(nil, errors.New("API error")).Once()

	character, err := charService.CreateCharacter("Krillin")
	assert.Error(t, err)
	assert.Nil(t, character)
	assert.Contains(t, err.Error(), "failed to fetch character from external API")
	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "SaveCharacter")
}

func TestCharacterService_CreateCharacter_SaveError(t *testing.T) {
	mockRepo := new(MockCharacterRepository)
	mockAPIClient := new(MockDragonBallAPIClient)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	charService := services.NewCharacterService(mockRepo, mockAPIClient, logger)

	apiCharacter := &domain.Character{
		ID:   "789",
		Name: "Piccolo",
	}

	// Expect FindCharacterByName from DB to return nil (not found)
	mockRepo.On("FindCharacterByName", "Piccolo").Return(nil, nil).Once()
	// Expect FindCharacterByName from API to return a character
	mockAPIClient.On("FindCharacterByName", "Piccolo").Return(apiCharacter, nil).Once()
	// Expect SaveCharacter to return an error
	mockRepo.On("SaveCharacter", mock.AnythingOfType("*domain.Character")).Return(errors.New("DB save error")).Once()

	character, err := charService.CreateCharacter("Piccolo")
	assert.Error(t, err)
	assert.Nil(t, character)
	assert.Contains(t, err.Error(), "failed to save character")
	mockRepo.AssertExpectations(t)
	mockAPIClient.AssertExpectations(t)
}
