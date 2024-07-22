package db

import (
	"github.com/google/uuid"
	"mood/models"
)

type DB interface {
	GetUserById(uuid.UUID) (*models.User, error)
	InsertUser(*models.User, string) (*uuid.UUID, error)
	LoginUser(string, string) (*uuid.UUID, error)
	InsertEntry(*models.Entry) (*uuid.UUID, error)
	GetEntriesByUserId(uuid.UUID) ([]models.Entry, error)
}
