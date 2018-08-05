package middlewares

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/http/util"
	log "github.com/sirupsen/logrus"
)

const TokenCookie = "token"
const LogInCookie = "is_logged_in"
const UserCookie = "user"

// AuthMiddleware represents a HTTP middleware handler for user authentication.
type AuthMiddleware struct {
	// custom logger object.
	logger *log.Entry

	// store user sessions.
	sessionMutex sync.Mutex
	sessions     map[string]uint

	// external services.
	users UserManager
}

// NewAuthMiddleware returns a new instance of the auth middleware handler.
func NewAuthMiddleware(users UserManager) *AuthMiddleware {
	m := &AuthMiddleware{
		logger:   log.WithFields(log.Fields{"package": "http", "module": "authMiddleware"}),
		sessions: make(map[string]uint),
		users:    users,
	}
	return m
}

// SetUserStatus sets whether the user is logged in or not.
func (m AuthMiddleware) SetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie(TokenCookie); err == nil && token != "" {
			m.sessionMutex.Lock()
			id := m.sessions[token]
			m.sessionMutex.Unlock()

			// get user data.
			user, err := m.users.Find(id)
			if (err == nil && user != coin.User{}) {
				c.Set(LogInCookie, true)
				c.Set(UserCookie, user)
				return
			}
		}
		c.Set(LogInCookie, false)
	}
}

// AddSession adds a new active session.
func (m *AuthMiddleware) AddSession(c *gin.Context, user uint) {
	defer m.sessionMutex.Unlock()
	m.sessionMutex.Lock()

	t := util.CreateSessionToken()
	c.SetCookie(TokenCookie, t, 3600, "", "", false, true)
	c.Set(LogInCookie, true)
	m.sessions[t] = user
}

// RemoveSession removes a new active session.
func (m *AuthMiddleware) RemoveSession(c *gin.Context) {
	defer m.sessionMutex.Unlock()
	m.sessionMutex.Lock()

	c.SetCookie(TokenCookie, "", -1, "", "", false, true)
	c.Set(LogInCookie, false)
	if token, err := c.Cookie(TokenCookie); err == nil || token != "" {
		delete(m.sessions, token)
	}
}

// UserManager defines the interface to interact with the user persistence layer.
type UserManager interface {
	Find(id uint) (coin.User, error)
}
