package models

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();"`
	Name string
	Email string `gorm:"unique"`
	Password string
	RefreshToken string
	IP string
}