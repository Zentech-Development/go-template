package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/pkg/service"
	"github.com/gin-gonic/gin"
)

type GinBindingOpts struct {
	DebugMode     bool
	SecretKey     string
	ListenAddr    string
	UseCSRFTokens bool
	CSRFSecret    string
}

// GinBinding represents a Gin application bound to services.
type GinBinding struct {
	opts    GinBindingOpts
	service *service.Service
	app     *gin.Engine
}

// NewBinding initializes an instance of GinBinding with the provided values.
func NewBinding(services *service.Service, opts GinBindingOpts) *GinBinding {
	ginBinding := &GinBinding{
		service: services,
		opts:    opts,
	}

	if opts.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.SetTrustedProxies(nil)

	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	app.Use(getCookieSessionMiddleware(opts.SecretKey))

	ginBinding.app = app

	ginBinding.attachHandlers()

	return ginBinding
}

// Run starts the application with the provided configuration.
func (b *GinBinding) Run() error {
	return b.app.Run(b.opts.ListenAddr)
}

// attachHandlers adds the handlers to the underlying Gin app.
func (b *GinBinding) attachHandlers() {
	b.app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	b.app.POST("/api/v1/auth/login", handleLogin)
	b.app.POST("/api/v1/auth/register")
	b.app.GET("/api/v1/auth/logout", handleLogout)

	b.app.Use(requireAuth)

	b.app.GET("/api/v1/auth/me", handleAuthMe)
}
