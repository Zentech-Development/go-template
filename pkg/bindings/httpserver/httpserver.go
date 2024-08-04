package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/pkg/service"
	"github.com/gin-gonic/gin"
)

type HTTPServerOpts struct {
	DebugMode  bool
	SecretKey  string
	ListenAddr string
}

// HTTPServer represents a Gin application bound to services.
type HTTPServer struct {
	opts     HTTPServerOpts
	services *service.Service
	app      *gin.Engine
}

// NewBinding initializes an instance of GinBinding with the provided values.
func NewBinding(services *service.Service, opts HTTPServerOpts) *HTTPServer {
	ginBinding := &HTTPServer{
		services: services,
		opts:     opts,
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
func (b *HTTPServer) Run() error {
	return b.app.Run(b.opts.ListenAddr)
}

// attachHandlers adds the handlers to the underlying Gin app.
func (b *HTTPServer) attachHandlers() {
	b.app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	b.app.POST("/api/v1/auth/login", b.handleLogin)
	b.app.POST("/api/v1/auth/register", b.handleRegister)

	b.app.Use(requireAuth)

	b.app.GET("/api/v1/auth/logout", b.handleLogout)
	b.app.GET("/api/v1/auth/me", b.handleAuthMe)
}
