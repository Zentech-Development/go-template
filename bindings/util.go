package bindings

import (
	"net/http"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/gin-gonic/gin"
)

func requireAccessToken(c *gin.Context) {
	token, err := c.Cookie(domain.GetConfig().TokenName)
	if err != nil {
		handleUnauthorizedRequest(c)
		return
	}

	claims, err := verifyAccessToken(token, domain.GetConfig().SecretKey, domain.GetConfig().AppName)
	if err != nil {
		handleUnauthorizedRequest(c)
		return
	}

	c.Set("userId", claims.Subject)
}

func handleUnauthorizedRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": "User is not logged in",
	})
}

func handleBadRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "Bad request",
	})
}
