package services

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
	"xm-task/dmodels"
	"xm-task/smodels"
)

func (s *ServiceFacade) CreateCompany(company smodels.Company) (dmodels.Company, error) {
	ct, err := s.dao.GetCompanyTypeByName(company.Type)
	if err != nil {
		return dmodels.Company{}, fmt.Errorf("dao.GetCompanyTypeByName: %v", err)
	}

	createdCompany, err := s.dao.CreateCompany(dmodels.Company{
		ID:          uuid.NewV4(),
		Name:        company.Name,
		Description: company.Description,
		Employees:   company.Employees,
		Registered:  company.Registered,
		TypeID:      ct.ID,
	})
	if err != nil {
		return dmodels.Company{}, fmt.Errorf("dao.CreateCompany: %v", err)
	}

	topic := "created-companies"
	message, _ := json.Marshal(createdCompany)
	err = s.kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
	}

	go s.kafka.Flush(200)

	return createdCompany, nil
}

func (s *ServiceFacade) UpdateCompany(company smodels.Company) (dmodels.Company, error) {
	cUUID, err := uuid.FromString(company.ID)
	if err != nil {
		return dmodels.Company{}, fmt.Errorf("uuid.FromString: %v", err)
	}

	ct, err := s.dao.GetCompanyTypeByName(company.Type)
	if err != nil {
		return dmodels.Company{}, fmt.Errorf("dao.GetCompanyTypeByName: %v", err)
	}

	updatedCompany, err := s.dao.UpdateCompany(dmodels.Company{
		ID:          cUUID,
		Name:        company.Name,
		Description: company.Description,
		Employees:   company.Employees,
		Registered:  company.Registered,
		TypeID:      ct.ID,
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return dmodels.Company{}, fmt.Errorf("dao.UpdateCompany: %v", err)
	}

	topic := "updated-companies"
	message, _ := json.Marshal(updatedCompany)
	err = s.kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
	}

	return updatedCompany, nil
}

func (s *ServiceFacade) GetCompanyByID(id string) (dmodels.CompanyShow, error) {
	company, err := s.dao.GetCompanyByID(id)
	if err != nil {
		return dmodels.CompanyShow{}, fmt.Errorf("dao.GetCompanyByID: %v", err)
	}

	return company, nil
}

func (s *ServiceFacade) DeleteCompanyByID(id string) error {
	company, err := s.dao.GetCompanyByID(id)
	if err != nil {
		return fmt.Errorf("dao.GetCompanyByID: %v", err)
	}

	err = s.dao.DeleteCompanyByID(id)
	if err != nil {
		return fmt.Errorf("dao.DeleteCompanyByID: %v", err)
	}

	topic := "deleted-companies"
	message, _ := json.Marshal(company)
	err = s.kafka.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message: %s\n", err)
	}

	return nil
}
