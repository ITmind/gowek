package repo

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Text   string
	UserID uint
	User   User
}

type User struct {
	gorm.Model
	Login string
	Email string
	Hash  string
}
