package dmodels

import uuid "github.com/satori/go.uuid"

const UsersTable = "users"

type User struct {
	ID       uuid.UUID `gorm:"column:id;PRIMARY_KEY"`
	Email    string    `gorm:"column:email"`
	Password string    `gorm:"column:password"`
}
