package ginbinding

import (
	"net/http"

	"github.com/Zentech-Development/go-template/pkg/service"
	"github.com/gin-gonic/gin"
)

// GinBinding represents a Gin application bound to services.
type GinBinding struct {
	service    *service.Service
	debugMode  bool
	secretKey  string
	listenAddr string
	app        *gin.Engine
}

// NewBinding initializes an instance of GinBinding with the provided values.
func NewBinding(services *service.Service, listenAddr string, debugMode bool, secretKey string) *GinBinding {
	ginBinding := &GinBinding{
		service:    services,
		debugMode:  debugMode,
		secretKey:  secretKey,
		listenAddr: listenAddr,
	}

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.SetTrustedProxies(nil)
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	ginBinding.app = app

	ginBinding.attachHandlers()

	return ginBinding
}

// Run starts the application with the provided configuration.
func (b *GinBinding) Run() error {
	return b.app.Run(b.listenAddr)
}

// attachHandlers adds the handlers to the underlying Gin app.
func (b *GinBinding) attachHandlers() {
	b.app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
}
