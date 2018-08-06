package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin"
)

const LogInCookie = "is_logged_in"
const UserCookie = "user"

// render returns either HTML or JSON based on the 'Accept' header of the request (defaults to HTML)
func Render(c *gin.Context, data gin.H, template string) {
	// check whether the user is logged in.
	if loggedIn, exists := c.Get(LogInCookie); exists {
		data[LogInCookie] = loggedIn.(bool)
	}
	if user, exists := c.Get(UserCookie); exists {
		data[UserCookie] = user.(coin.User)
	}

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, template, data)
	}
}

// RequestError represents a request error.
type RequestError struct {
	Title   string
	Message string
}

func (r RequestError) Render() map[string]interface{} {
	return gin.H{
		"ErrorTitle":   r.Title,
		"ErrorMessage": r.Message,
	}
}

// RequestSuccess represents a request success.
type RequestSuccess struct {
	Title   string
	Message string
}

func (r RequestSuccess) Render() map[string]interface{} {
	return gin.H{
		"MessageTitle":   r.Title,
		"MessageMessage": r.Message,
	}
}
