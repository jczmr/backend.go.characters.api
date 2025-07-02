package domain

import "time"

type Character struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Ki        string    `json:"ki"`
	Race      string    `json:"race"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NewCharacterRequest struct {
	Name string `json:"name" binding:"required"`
}
