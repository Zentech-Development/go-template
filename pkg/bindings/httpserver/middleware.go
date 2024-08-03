package ginbinding

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const userKey = "userID"

func requireAuth(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Next()
}

func getCookieSessionMiddleware(secretKey string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(secretKey))
	store.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	return sessions.Sessions("APPNAME_session", store)
}
