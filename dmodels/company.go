package dmodels

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const CompaniesTable = "companies"

type Company struct {
	ID          uuid.UUID `gorm:"column:id;PRIMARY_KEY"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Employees   uint64    `gorm:"column:employees"`
	Registered  bool      `gorm:"column:registered;default:false"`
	TypeID      uint64    `gorm:"column:type_id"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:now()"`
}

type CompanyShow struct {
	ID          uuid.UUID `gorm:"column:id;PRIMARY_KEY"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Employees   uint64    `gorm:"column:employees"`
	Registered  bool      `gorm:"column:registered;default:false"`
	Type        string    `gorm:"column:type"`
}
