package postgres_test

import (
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	"backend.go.characters.api/internal/adapters/secondary/db/postgres"
	"backend.go.characters.api/internal/core/domain"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCharacterRepositorySaveCharacter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := postgres.NewCharacterRepository(db, logger)

	character := &domain.Character{
		ID:   "1",
		Name: "Goku",
		Ki:   "10000",
		Race: "Saiyan",
	}

	// Expect the INSERT or UPDATE query
	mock.ExpectExec(`INSERT INTO characters`).
		WithArgs(character.ID, character.Name, character.Ki, character.Race).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveCharacter(character)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCharacterRepositoryFindCharacterByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := postgres.NewCharacterRepository(db, logger)

	characterName := "Vegeta"
	rows := sqlmock.NewRows([]string{"id", "name", "ki", "race", "created_at", "updated_at"}).
		AddRow("2", "Vegeta", "9000", "Saiyan", time.Now(), time.Now())

	// Expect the SELECT query
	mock.ExpectQuery(`SELECT id, name, ki, race, created_at, updated_at FROM characters WHERE name ILIKE \$1`).
		WithArgs(characterName).
		WillReturnRows(rows)

	foundCharacter, err := repo.FindCharacterByName(characterName)
	assert.NoError(t, err)
	assert.NotNil(t, foundCharacter)
	assert.Equal(t, characterName, foundCharacter.Name)
	assert.NoError(t, mock.ExpectationsWereMet())

	// Test not found case
	mock.ExpectQuery(`SELECT id, name, ki, race, created_at, updated_at FROM characters WHERE name ILIKE \$1`).
		WithArgs("NonExistent").
		WillReturnError(sql.ErrNoRows)

	notFoundCharacter, err := repo.FindCharacterByName("NonExistent")
	assert.NoError(t, err)
	assert.Nil(t, notFoundCharacter)
	assert.NoError(t, mock.ExpectationsWereMet())
}
