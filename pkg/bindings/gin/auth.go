package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/pkg/entities"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func handleLogin(c *gin.Context) {
	session := sessions.Default(c)

	var loginInput entities.AccountLoginInput

	if err := c.ShouldBind(&loginInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing request data"})
		return
	}

	if loginInput.Username != "hello" || loginInput.Password != "itsme" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	session.Set(userKey, loginInput.Username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func handleLogout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func handleAuthMe(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	c.JSON(http.StatusOK, gin.H{"userID": user})
}
