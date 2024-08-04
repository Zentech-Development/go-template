package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/internals/entities"
	"github.com/Zentech-Development/go-template/pkg/logger"
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

	logger.L.Info().Interface("account", account).Msg("Account authenticated")

	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated account", "account": account})
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

	account, err := b.services.CreateAccount(c, input)
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

	logger.L.Info().Interface("account", account).Msg("Account created")

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created account", "account": account})
}

func (b *HTTPServer) handleAuthMe(c *gin.Context) {
	userID, _ := c.Get(userKey)

	account, err := b.services.GetAccountByID(c, userID.(int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Failed to retrieve current account info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account is logged in", "account": account})
}
