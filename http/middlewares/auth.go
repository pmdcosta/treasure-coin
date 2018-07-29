package middlewares

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	)

// AuthMiddleware represents a HTTP middleware handler for user authentication.
type AuthMiddleware struct {
	// custom logger object.
	logger *log.Entry
}

// NewAuthMiddleware returns a new instance of the auth middleware handler.
func NewAuthMiddleware() *AuthMiddleware {
	m := &AuthMiddleware{
		logger: log.WithFields(log.Fields{"package": "http", "module": "authMiddleware"}),
	}
	return m
}

// SetUserStatus sets whether the user is logged in or not
func (m AuthMiddleware) SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set("is_logged_in", true)
		} else {
			c.Set("is_logged_in", false)
		}


	}
}