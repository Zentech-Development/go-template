package bindings

import (
	"fmt"
	"net/http"

	"github.com/Zentech-Development/go-template/domain"
	"github.com/Zentech-Development/go-template/public/pages"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	csrf "github.com/utrack/gin-csrf"
)

func NewServerBinding(h *domain.Handlers, config *domain.ApplicationConfig) *gin.Engine {
	if config.Lifecycle == domain.LIFECYCLE_PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app := initializeApp(config.SecretKey, config.TokenName)
	setupMiddleware(app, config.UseCSRFTokens, config.CSRFSecret)
	setupEndpoints(app, h)

	return app
}

func initializeApp(secretKey string, cookieName string) *gin.Engine {
	app := gin.Default()
	app.SetTrustedProxies(nil)

	store := cookie.NewStore([]byte(secretKey))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	app.Use(sessions.Sessions(cookieName, store))

	return app
}

func setupMiddleware(app *gin.Engine, useCSRFTokens bool, csrfSecret string) {
	app.Use(func(c *gin.Context) {
		c.Set("requestId", uuid.NewString())
		c.Next()
	})

	csrfIgnoreMethods := []string{"GET", "HEAD", "OPTIONS"}

	if !useCSRFTokens {
		csrfIgnoreMethods = append(csrfIgnoreMethods, "POST")
	}

	app.Use(csrf.Middleware(csrf.Options{
		IgnoreMethods: csrfIgnoreMethods,
		Secret:        csrfSecret,
		ErrorFunc: func(c *gin.Context) {
			sendJSONOrRedirect(
				c,
				http.StatusBadRequest,
				&gin.H{
					"message": fmt.Sprintf("[Request ID: %s]: CSRF token missing", c.GetString("requestId")),
				},
				URLs.LoginPage,
			)
		},
	}))
}

func setupEndpoints(app *gin.Engine, handlers *domain.Handlers) {
	app.Static("assets", "public/assets")
	app.StaticFile("favicon.ico", "public/assets/favicon.png")

	accountHandlers := newAccountsBinding(handlers)

	app.GET(URLs.LandingPage, func(c *gin.Context) {
		_ = pages.Landing().Render(c, c.Writer)
	})
	app.GET(URLs.LoginPage, accountHandlers.ViewLogin)
	app.GET(URLs.RegisterPage, accountHandlers.ViewRegister)

	app.POST("/api/v1/login", accountHandlers.Login)
	app.POST("/api/v1/accounts", accountHandlers.Create)

	app.Use(requireSession)

	app.GET(URLs.HomePage, func(c *gin.Context) {
		sendJSONOrHTML(
			c,
			http.StatusOK,
			&gin.H{
				"message": "Ok",
			},
			pages.Home(c.GetString("userEmail")),
		)
	})

	apiV1 := app.Group("/api/v1")
	{
		accountsRouter := apiV1.Group("/accounts")
		{
			accountsRouter.GET("/logout", accountHandlers.Logout)
			accountsRouter.GET("/me", accountHandlers.GetMe)
			accountsRouter.DELETE("/:id", accountHandlers.Delete)
		}
	}
}
