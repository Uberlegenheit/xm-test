package dmodels

const CompanyTypesTable = "company_types"

type CompanyType struct {
	ID   uint64 `gorm:"column:id;PRIMARY_KEY"`
	Name string `gorm:"column:name"`
}
