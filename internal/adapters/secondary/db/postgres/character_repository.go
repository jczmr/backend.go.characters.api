package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"log/slog"

	"backend.go.characters.api/internal/core/domain"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type characterRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewCharacterRepository(db *sql.DB, logger *slog.Logger) *characterRepository {
	return &characterRepository{db: db, logger: logger}
}

func (r *characterRepository) SaveCharacter(character *domain.Character) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `
		INSERT INTO characters (id, name, ki, race, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, ki = EXCLUDED.ki, race = EXCLUDED.race, updated_at = NOW();
	`
	_, err := r.db.ExecContext(ctx, query, character.ID, character.Name, character.Ki, character.Race)
	if err != nil {
		r.logger.Error("Failed to save character to database", slog.String("error", err.Error()), slog.String("character_id", character.ID))
		return fmt.Errorf("failed to save character: %w", err)
	}
	r.logger.Info("Character saved successfully to database", slog.String("character_id", character.ID))
	return nil
}

func (r *characterRepository) FindCharacterByName(name string) (*domain.Character, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	query := `SELECT id, name, ki, race, created_at, updated_at FROM characters WHERE name ILIKE $1;`
	row := r.db.QueryRowContext(ctx, query, name)

	character := &domain.Character{}
	err := row.Scan(&character.ID, &character.Name, &character.Ki, &character.Race, &character.CreatedAt, &character.UpdatedAt)
	if err == sql.ErrNoRows {
		r.logger.Info("Character not found in database by name", slog.String("character_name", name))
		return nil, nil // Character not found
	}
	if err != nil {
		r.logger.Error("Failed to query character by name from database", slog.String("error", err.Error()), slog.String("character_name", name))
		return nil, fmt.Errorf("failed to find character by name: %w", err)
	}
	r.logger.Info("Character found in database by name", slog.String("character_name", name), slog.String("character_id", character.ID))
	return character, nil
}
