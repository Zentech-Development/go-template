package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/Zentech-Development/go-template/public/pages"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type AccountsBinding struct {
	Handlers     *domain.Handlers
	PathHome     string
	PathLogin    string
	PathRegister string
}

func newAccountsBinding(handlers *domain.Handlers) AccountsBinding {
	return AccountsBinding{
		Handlers: handlers,
	}
}

func (b AccountsBinding) Create(c *gin.Context) {
	var accountInput domain.AccountInput

	if err := c.ShouldBind(&accountInput); err != nil {
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

	addSession(c, account.ID, account.Email)

	sendJSONOrRedirect(
		c,
		http.StatusCreated,
		&gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Added account successfully", c.GetString("requestId")),
			"account": account,
		},
		URLs.Home,
	)
}

func (b AccountsBinding) Login(c *gin.Context) {
	var input domain.LoginInput

	if err := c.ShouldBind(&input); err != nil {
		sendJSONOrRedirect(
			c,
			http.StatusUnauthorized,
			&gin.H{
				"message": fmt.Sprintf("[Request ID: %s]: Failed to parse request", c.GetString("requestId")),
			},
			fmt.Sprintf("%s?login_failed=true", URLs.Login),
		)
		return
	}

	account, err := b.Handlers.Accounts.Login(input)
	if err != nil {
		sendJSONOrRedirect(
			c,
			http.StatusUnauthorized,
			&gin.H{
				"message": fmt.Sprintf("[Request ID: %s]: Failed to login", c.GetString("requestId")),
			},
			fmt.Sprintf("%s?login_failed=true", URLs.Login),
		)
		return
	}

	addSession(c, account.ID, account.Email)

	sendJSONOrRedirect(
		c,
		http.StatusCreated,
		&gin.H{
			"message": fmt.Sprintf("[Request ID: %s]: Login successful", c.GetString("requestId")),
			"account": account,
		},
		URLs.Home,
	)
}

func (b AccountsBinding) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (b AccountsBinding) ViewLogin(c *gin.Context) {
	csrfToken := csrf.GetToken(c)

	_ = pages.Login(csrfToken).Render(c, c.Writer)
}

func (b AccountsBinding) ViewRegister(c *gin.Context) {
	csrfToken := csrf.GetToken(c)

	_ = pages.Register(csrfToken).Render(c, c.Writer)
}

func (b AccountsBinding) GetMe(c *gin.Context) {
	userId, _ := c.Get("userId")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("[Request ID: %s]: User is logged in", c.GetString("requestId")),
		"userId":  userId,
	})
}

func (b AccountsBinding) Delete(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
