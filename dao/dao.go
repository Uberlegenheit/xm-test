package dao

import (
	"fmt"
	"time"

	"xm-task/conf"
	"xm-task/dao/cache"
	"xm-task/dao/postgres"
	"xm-task/dmodels"
)

type (
	DAO interface {
		Postgres
		Cache
	}

	Postgres interface {
		CheckDBStatus() bool

		CreateUser(user dmodels.User) (dmodels.User, error)
		GetUserByEmail(email string) (dmodels.User, error)

		CreateCompany(company dmodels.Company) (dmodels.Company, error)
		UpdateCompany(company dmodels.Company) (dmodels.Company, error)
		GetCompanyByID(id string) (dmodels.CompanyShow, error)
		DeleteCompanyByID(id string) error

		GetCompanyTypeByName(name string) (dmodels.CompanyType, error)
	}

	Cache interface {
		AddAuthToken(key string, item interface{}, expiration time.Duration) error
		GetAuthToken(token string) (interface{}, bool, error)
		RemoveAuthToken(key string) error
	}

	daoImpl struct {
		*postgres.Postgres
		*cache.Cache
	}
)

func New(cfg conf.Config, migrate bool) (DAO, error) {
	pg, err := postgres.NewPostgres(cfg.Postgres, migrate)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres: %s", err.Error())
	}

	return daoImpl{
		Postgres: pg,
		Cache:    cache.NewCache(),
	}, nil
}
