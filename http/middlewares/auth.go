package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/http/util"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const TokenCookie = "token"

// AuthMiddleware represents a HTTP middleware handler for user authentication.
type AuthMiddleware struct {
	logger *log.Entry

	// external services.
	users    UserManager
	sessions SessionManager
}

// NewAuthMiddleware returns a new instance of the auth middleware handler.
func NewAuthMiddleware(users UserManager, sessions SessionManager) *AuthMiddleware {
	m := &AuthMiddleware{
		logger:   log.WithFields(log.Fields{"package": "http", "module": "auth-middleware"}),
		users:    users,
		sessions: sessions,
	}
	return m
}

// SetUserStatus sets whether the user is logged in or not.
func (m AuthMiddleware) SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie(TokenCookie); err == nil && token != "" {
			// get the user id from the session.
			s, _ := m.sessions.Find(token)

			// get user data from the db.
			user, err := m.users.Find(s)
			if (err == nil && user != coin.User{}) {
				c.Set(util.LogInCookie, true)
				c.Set(util.UserCookie, user)
				m.logger.WithFields(log.Fields{"token": token, "session": s, "user": user.Email}).Debug("current session")
				return
			}
		}
		c.Set(util.LogInCookie, false)
	}
}

// AddSession adds a new active session.
func (m *AuthMiddleware) AddSession(c *gin.Context, user string) {
	t := CreateSessionToken()
	m.logger.WithFields(log.Fields{"token": t, "user": user}).Debug("creating sessions")
	c.SetCookie(TokenCookie, t, 3600, "", "", false, true)
	c.Set(util.LogInCookie, true)
	m.sessions.Add(t, user)
	m.logger.WithFields(log.Fields{"user": user, "token": t}).Info("user signing in")
}

// RemoveSession removes a new active session.
func (m *AuthMiddleware) RemoveSession(c *gin.Context) {
	c.SetCookie(TokenCookie, "", -1, "", "", false, true)
	c.Set(util.LogInCookie, false)
	if token, err := c.Cookie(TokenCookie); err == nil || token != "" {
		m.sessions.Remove(token)
		m.logger.WithFields(log.Fields{"token": token}).Debug("removing sessions")
	}
}

// CreateSessionToken generate a new session token to store in the cookie.
func CreateSessionToken() string {
	token, _ := uuid.NewV4()
	return token.String()
}

// UserManager defines the interface to interact with the user persistence layer.
type UserManager interface {
	Find(email string) (coin.User, error)
}

// SessionManager defines the interface to interact with the session persistence layer.
type SessionManager interface {
	Add(token, session string) error
	Find(token string) (string, error)
	Remove(token string) error
}
