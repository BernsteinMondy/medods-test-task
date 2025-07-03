package entities

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Username     string
	Hash         string
	RefreshToken string
}

type RefreshToken struct {
	Token  string
	IsUsed bool
}
