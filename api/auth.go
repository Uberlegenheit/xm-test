package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"xm-task/log"
)

func (api *API) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := api.services.ExtractTokenMetadata(c)
		if err != nil {
			log.Error("[api] AuthMiddleware: ExtractTokenMetadata", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		email, err := api.services.FetchAuth(ad)
		if err != nil || email == "" {
			log.Error("[api] AuthMiddleware: FetchAuth", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot fetch active user"})
			c.Abort()
			return
		}

		user, err := api.services.GetUserByEmail(email)
		if err != nil {
			log.Error("[api] AuthMiddleware: GetUserByEmail", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
