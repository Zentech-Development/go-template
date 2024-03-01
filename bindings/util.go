package bindings

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func requireSession(c *gin.Context) {
	session := sessions.Default(c)
	id := session.Get("id")
	email := session.Get("email")

	if id == nil || email == nil || len(email.(string)) < 1 {
		sendJSONOrRedirect(c, http.StatusUnauthorized, &gin.H{}, "/login")
	}

	c.Set("userId", id)
	c.Set("userEmail", email)

	c.Next()
}

func addSession(c *gin.Context, id string, email string) {
	session := sessions.Default(c)
	session.Set("id", id)
	session.Set("email", email)
	session.Save()
}

func sendJSONOrHTML(c *gin.Context, status int, data *gin.H, template templ.Component) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(status, &data)
		return
	}

	template.Render(c, c.Writer)
}

func sendJSONOrRedirect(c *gin.Context, status int, data *gin.H, target string) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(status, &data)
		return
	}

	c.Redirect(http.StatusMovedPermanently, target)
}
