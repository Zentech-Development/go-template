package bindings

import (
	"net/http"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/Zentech-Development/go-template/public/pages"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewServerBinding(h *domain.Handlers, config *domain.ApplicationConfig) *gin.Engine {
	if config.Lifecycle == domain.LIFECYCLE_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := initializeApp()
	setupMiddleware(app)
	setupEndpoints(app, h)

	return app
}

func initializeApp() *gin.Engine {
	app := gin.Default()
	app.SetTrustedProxies(nil)

	return app
}

func setupMiddleware(app *gin.Engine) {
	app.Use(func(c *gin.Context) {
		c.Set("requestId", uuid.NewString())
	})
}

func setupEndpoints(app *gin.Engine, handlers *domain.Handlers) {
	app.Static("assets", "public/assets")
	app.StaticFile("favicon.ico", "public/assets/favicon.png")

	accountHandlers := newAccountsBinding(handlers)

	app.GET("/", func(c *gin.Context) {
		if c.GetHeader("Accept") == "application/json" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Ok",
			})
			return
		}

		if err := pages.Home().Render(c, c.Writer); err != nil {
			// TODO: not sure what I want here
		}
	})
	app.POST("/api/v1/login", accountHandlers.Login)
	app.POST("/api/v1/accounts", accountHandlers.Create)

	app.Use(requireAccessToken)

	apiV1 := app.Group("/api/v1")
	{
		accountsRouter := apiV1.Group("/accounts")
		{
			accountsRouter.POST("/logout", accountHandlers.Logout)
			accountsRouter.GET("/me", accountHandlers.GetMe)
			accountsRouter.DELETE("/:id", accountHandlers.Delete)
		}
	}
}
