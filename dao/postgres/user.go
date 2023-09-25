package postgres

import "xm-task/dmodels"

func (db *Postgres) CreateUser(user dmodels.User) (dmodels.User, error) {
	err := db.db.Table(dmodels.UsersTable).Create(&user).Error
	return user, err
}

func (db *Postgres) GetUserByEmail(email string) (dmodels.User, error) {
	var user dmodels.User
	err := db.db.Table(dmodels.UsersTable).
		Select("*").
		Where("email = ?", email).
		First(&user).Error
	return user, err
}
