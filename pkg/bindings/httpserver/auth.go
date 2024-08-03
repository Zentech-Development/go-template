package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/pkg/entities"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (b *HTTPServer) handleLogin(c *gin.Context) {
	var input entities.AccountLoginInput

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing request data"})
		return
	}

	account, err := b.services.Login(c, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set(userKey, account.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

func (b *HTTPServer) handleLogout(c *gin.Context) {
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

func (b *HTTPServer) handleRegister(c *gin.Context) {
	var input entities.AccountCreateInput

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing request data"})
		return
	}

	account, err := b.services.Create(c, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to create account"})
		return
	}

	session := sessions.Default(c)
	session.Set(userKey, account.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully created user"})
}

func (b *HTTPServer) handleAuthMe(c *gin.Context) {
	userID, _ := c.Get(userKey)
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
