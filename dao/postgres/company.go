package postgres

import (
	"fmt"

	"xm-task/dmodels"
)

func (db *Postgres) CreateCompany(company dmodels.Company) (dmodels.Company, error) {
	err := db.db.Table(dmodels.CompaniesTable).Create(&company).Error
	return company, err
}

func (db *Postgres) UpdateCompany(company dmodels.Company) (dmodels.Company, error) {
	err := db.db.Table(dmodels.CompaniesTable).
		Where("id = ?", company.ID.String()).
		Updates(&company).Error
	return company, err
}

func (db *Postgres) GetCompanyByID(id string) (dmodels.CompanyShow, error) {
	var company dmodels.CompanyShow
	err := db.db.Table(fmt.Sprintf("%s c", dmodels.CompaniesTable)).
		Select("c.id, c.name, c.description, c.employees, c.registered, ct.name as type").
		Where("c.id = ?", id).
		Joins("inner join company_types ct on ct.id = c.type_id").
		Scan(&company).Error
	return company, err
}

func (db *Postgres) DeleteCompanyByID(id string) error {
	return db.db.Table(dmodels.CompaniesTable).
		Where("id = ?", id).
		Delete(&dmodels.Company{}).Error
}
