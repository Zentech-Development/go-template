package bindings

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/gin-gonic/gin"
)

type AccountsBinding struct {
	Handlers *domain.Handlers
}

func newAccountsBinding(handlers *domain.Handlers) AccountsBinding {
	return AccountsBinding{
		Handlers: handlers,
	}
}

func (b AccountsBinding) Create(c *gin.Context) {
	var accountInput domain.AccountInput

	if err := c.ShouldBindJSON(&accountInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
		})
		return
	}

	account, err := b.Handlers.Accounts.Add(accountInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Unexpected error while adding account: %s", c.GetString("requestId"), err.Error()),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("[Request ID: %s]: Added account successfully", c.GetString("requestId")),
		"account": account,
	})
}

func (b AccountsBinding) Login(c *gin.Context) {
	var input domain.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
		})
		return
	}

	account, err := b.Handlers.Accounts.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Failed to login", c.GetString("requestId")),
		})
		return
	}

	token, err := generateAccessToken(
		account.Email,
		domain.GetConfig().TokenExpirationSeconds,
		domain.GetConfig().AppName,
		time.Now(),
		domain.GetConfig().SecretKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Unexpected error while logging in: %s", c.GetString("requestId"), err.Error()),
		})
		return
	}

	c.SetCookie(domain.GetConfig().TokenName, token, domain.GetConfig().TokenExpirationSeconds*1000, "/", domain.GetConfig().Host, true, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("[Request ID: %s]: Login successful", c.GetString("requestId")),
		"token":   token,
		"account": account,
	})
}

func (b AccountsBinding) Logout(c *gin.Context) {
	c.SetCookie(domain.GetConfig().TokenName, "", 1, "/", domain.GetConfig().Host, true, true)

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("[Request ID: %s]: Logout successful", c.GetString("requestId")),
	})
}

func (b AccountsBinding) GetMe(c *gin.Context) {
	userId, _ := c.Get("userId")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("[Request ID: %s]: User is logged in", c.GetString("requestId")),
		"userId":  userId,
	})
}

func (b AccountsBinding) Delete(c *gin.Context) {}
