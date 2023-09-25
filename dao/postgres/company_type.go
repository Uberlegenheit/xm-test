package postgres

import "xm-task/dmodels"

func (db *Postgres) GetCompanyTypeByName(name string) (dmodels.CompanyType, error) {
	var ct dmodels.CompanyType
	err := db.db.Table(dmodels.CompanyTypesTable).
		Select("*").
		Where("name = ?", name).
		First(&ct).Error
	return ct, err
}
