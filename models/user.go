package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt " db:"created_at"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

func NewUser(firstName, lastName, email string) *User {
	return &User{
		Id:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}
}
