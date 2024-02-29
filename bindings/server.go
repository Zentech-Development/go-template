package bindings

import (
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
	})

	if useCSRFTokens {
		app.Use(csrf.Middleware(csrf.Options{
			Secret: csrfSecret,
			ErrorFunc: func(c *gin.Context) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "CSRF token missing",
				})
			},
		}))
	}
}

func setupEndpoints(app *gin.Engine, handlers *domain.Handlers) {
	app.Static("assets", "public/assets")
	app.StaticFile("favicon.ico", "public/assets/favicon.png")

	accountHandlers := newAccountsBinding(handlers)

	app.POST("/api/v1/login", accountHandlers.Login)
	app.POST("/api/v1/accounts", accountHandlers.Create)
	app.GET("/login", accountHandlers.ViewLogin)
	app.GET("/register", accountHandlers.ViewRegister)

	app.Use(requireLogin)

	app.GET("/", func(c *gin.Context) {
		sendJSONOrHTML(
			c,
			http.StatusOK,
			&gin.H{
				"message": "Ok",
			},
			pages.Home(c.GetString("userId")),
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

func requireLogin(c *gin.Context) {
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil || len(email.(string)) < 1 {
		sendJSONOrRedirect(c, http.StatusUnauthorized, &gin.H{}, "/login")
	}

	c.Set("userId", email)

	c.Next()
}
