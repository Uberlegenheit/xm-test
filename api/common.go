package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xm-task/conf"
)

func (api *API) Index(c *gin.Context) {
	c.String(http.StatusOK, "This is a service '%s'", conf.Service)
}

func (api *API) Health(c *gin.Context) {
	status := api.services.CheckDBStatus()
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
