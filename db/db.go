package db

import (
	"github.com/google/uuid"
	"mood/models"
)

type DB interface {
	Get(uuid.UUID) (*models.User, error)
	CreateUser(*models.User, string) (*uuid.UUID, error)
	LoginUser(string, string) (*uuid.UUID, error)
}
