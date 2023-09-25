package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	local "xm-task/helpers/errors"
	"xm-task/log"
	"xm-task/smodels"
)

func (api *API) SignIn(c *gin.Context) {
	var user smodels.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("[api] SignIn: ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": local.BadRequest})
		return
	}

	if !user.Validate() {
		log.Error("[api] SignIn: Validate", zap.Error(fmt.Errorf("incorrect user data")))
		c.JSON(http.StatusBadRequest, gin.H{"error": local.BadRequest})
		return
	}

	ok, err := api.services.SignInOrRegister(user)
	if err != nil {
		log.Error("[api] SignIn: SignInOrRegister", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}
	if !ok {
		log.Error("[api] SignIn: incorrect password", zap.Error(fmt.Errorf("incorrect password")))
		c.JSON(http.StatusUnauthorized, gin.H{"error": local.UnauthorizedErr})
		return
	}

	td, err := api.services.CreateToken(user.Email)
	if err != nil {
		log.Error("[api] SignIn: CreateToken", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	err = api.services.CreateAuth(user.Email, td)
	if err != nil {
		log.Error("[api] SignIn: CreateAuth", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": local.ServiceError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"access_token":    td.AccessToken,
		"refresh_token":   td.RefreshToken,
		"access_expired":  td.AtExpires,
		"refresh_expired": td.RtExpires,
	})
}

func (api *API) Refresh(c *gin.Context) {
	td, err := api.services.Refresh(c.Request)
	if err != nil {
		log.Error("[api] Refresh: Refresh", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"access_token":    td.AccessToken,
		"refresh_token":   td.RefreshToken,
		"access_expired":  td.AtExpires,
		"refresh_expired": td.RtExpires,
	})
}

func (api *API) LogOut(c *gin.Context) {
	au, err := api.services.ExtractTokenMetadata(c)
	if err != nil {
		log.Error("[api] SignOut: ExtractTokenMetadata", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshUuid, _ := api.services.FetchAuth(smodels.AccessDetails{
		AccessUuid: fmt.Sprintf("%s_refresh", au.AccessUuid),
	})

	delErr := api.services.DeleteAuth(au.AccessUuid,
		refreshUuid,
		fmt.Sprintf("%s_refresh", au.AccessUuid),
		fmt.Sprintf("%s_access", refreshUuid),
		fmt.Sprintf("%s_td", au.Email),
	)
	if delErr != nil {
		log.Error("[api] SignOut: DeleteAuth", zap.Error(delErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": delErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
