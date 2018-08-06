package middlewares_test

import (
	"testing"

	"github.com/pmdcosta/treasure-coin/http/middlewares"
	log "github.com/sirupsen/logrus"
)

//go:generate moq -out user_service_mock.go . UserManager
//go:generate moq -out session_service_mock.go . SessionManager

// AuthMiddleware is a test wrapper.
type AuthMiddleware struct {
	AuthMiddleware *middlewares.AuthMiddleware
	Users          *middlewares.UserManagerMock
	Sessions       *middlewares.SessionManagerMock
}

// NewAuthMiddleware returns a new instance of AuthMiddleware.
func NewAuthMiddleware() *AuthMiddleware {
	log.SetLevel(log.DebugLevel)
	u := &middlewares.UserManagerMock{}
	s := &middlewares.SessionManagerMock{}
	m := &AuthMiddleware{
		AuthMiddleware: middlewares.NewAuthMiddleware(u, s),
	}
	return m
}

// TestAuthMiddleware_Create tests creating a new auth middleware.
func TestAuthMiddleware_Create(t *testing.T) {
	m := NewAuthMiddleware()
	if m == nil {
		t.Fatal("failed to create middleware")
	}
}

// TODO add more tests.
