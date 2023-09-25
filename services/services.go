package services

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"net/http"
	"xm-task/conf"
	"xm-task/dao"
	"xm-task/dmodels"
	"xm-task/smodels"
)

type (
	Service interface {
		CheckDBStatus() bool

		SignInOrRegister(user smodels.User) (bool, error)
		GetUserByEmail(email string) (dmodels.User, error)

		CreateCompany(company smodels.Company) (dmodels.Company, error)
		UpdateCompany(company smodels.Company) (dmodels.Company, error)
		GetCompanyByID(id string) (dmodels.CompanyShow, error)
		DeleteCompanyByID(id string) error

		CreateToken(email string) (smodels.TokenDetails, error)
		CreateAuth(email string, td smodels.TokenDetails) error
		ExtractTokenMetadata(c *gin.Context) (smodels.AccessDetails, error)
		Refresh(r *http.Request) (smodels.TokenDetails, error)
		FetchAuth(authD smodels.AccessDetails) (string, error)
		DeleteAuth(UUID ...string) error
	}

	ServiceFacade struct {
		cfg   conf.Config
		dao   dao.DAO
		kafka *kafka.Producer
	}
)

func NewService(cfg conf.Config, dao dao.DAO) (*ServiceFacade, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create producer: %s\n", err)
	}
	return &ServiceFacade{
		cfg:   cfg,
		dao:   dao,
		kafka: p,
	}, nil
}
