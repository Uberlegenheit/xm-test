package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	local "xm-task/helpers/errors"
	"xm-task/log"
	"xm-task/smodels"
)

func (api *API) CreateCompany(c *gin.Context) {
	var company smodels.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		log.Error("[api] CreateCompany: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": local.BadRequest})
		return
	}

	if err := company.Validate(); err != nil {
		log.Error("[api] CreateCompany: Validate", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbCompany, err := api.services.CreateCompany(company)
	if err != nil {
		log.Error("[api] CreateCompany: CreateCompany", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	c.JSON(http.StatusOK, smodels.Company{
		ID:          dbCompany.ID.String(),
		Name:        dbCompany.Name,
		Description: dbCompany.Description,
		Employees:   dbCompany.Employees,
		Registered:  dbCompany.Registered,
		Type:        company.Type,
	})
}

func (api *API) PatchCompany(c *gin.Context) {
	companyID := c.Param("id")

	var updatedCompany smodels.Company
	if err := c.ShouldBindJSON(&updatedCompany); err != nil {
		log.Error("[api] PatchCompany: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": local.BadRequest})
		return
	}

	if err := updatedCompany.Validate(); err != nil {
		log.Error("[api] PatchCompany: Validate", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedCompany.ID = companyID
	dbCompany, err := api.services.UpdateCompany(updatedCompany)
	if err != nil {
		log.Error("[api] PatchCompany: UpdateCompany", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	c.JSON(http.StatusOK, smodels.Company{
		ID:          dbCompany.ID.String(),
		Name:        dbCompany.Name,
		Description: dbCompany.Description,
		Employees:   dbCompany.Employees,
		Registered:  dbCompany.Registered,
		Type:        updatedCompany.Type,
	})
}

func (api *API) GetCompany(c *gin.Context) {
	companyID := c.Param("id")
	company, err := api.services.GetCompanyByID(companyID)
	if err != nil {
		log.Error("[api] GetCompany: GetCompanyByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	c.JSON(http.StatusOK, smodels.Company{
		ID:          company.ID.String(),
		Name:        company.Name,
		Description: company.Description,
		Employees:   company.Employees,
		Registered:  company.Registered,
		Type:        company.Type,
	})
}

func (api *API) DeleteCompany(c *gin.Context) {
	companyID := c.Param("id")
	err := api.services.DeleteCompanyByID(companyID)
	if err != nil {
		log.Error("[api] DeleteCompany: DeleteCompanyByID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
