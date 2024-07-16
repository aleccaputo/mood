package db

import (
	"github.com/google/uuid"
	"mood/models"
)

type MemoryDb struct{}

func NewMemoryDb() *MemoryDb {
	return &MemoryDb{}
}

func (s *MemoryDb) Get(id int) *models.User {
	return &models.User{
		Id:        uuid.New(),
		FirstName: "Alec",
	}
}
