package http

import (
	"net/http"

	"log/slog"

	"backend.go.characters.api/internal/core/domain"
	"backend.go.characters.api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	characterService ports.CharacterService
	logger           *slog.Logger
}

func NewCharacterHandler(characterService ports.CharacterService, logger *slog.Logger) *CharacterHandler {
	return &CharacterHandler{
		characterService: characterService,
		logger:           logger,
	}
}

func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
	var req domain.NewCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request payload for CreateCharacter", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	character, err := h.characterService.CreateCharacter(req.Name)
	if err != nil {
		h.logger.Error("Failed to create/retrieve character", slog.String("error", err.Error()), slog.String("character_name", req.Name))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Character processed successfully", slog.String("character_name", character.Name), slog.String("character_id", character.ID))
	c.JSON(http.StatusCreated, character)
}
