package util

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin"
)

// render returns either HTML or JSON based on the 'Accept' header of the request (defaults to HTML)
func Render(c *gin.Context, data gin.H, template string) {
	// check whether the user is logged in.
	if loggedIn, exists := c.Get("is_logged_in"); exists {
		data["is_logged_in"] = loggedIn.(bool)
	}
	if user, exists := c.Get("user"); exists {
		data["user"] = user.(coin.User)
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, template, data)
	}
}

// CreateSessionToken generate a new session token to store in the cookie.
func CreateSessionToken() string {
	return strconv.FormatInt(rand.Int63(), 16)
}
