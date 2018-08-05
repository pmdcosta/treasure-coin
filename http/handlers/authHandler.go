package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
	"github.com/pmdcosta/treasure-coin/http/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles the authentication routes in the server.
type AuthHandler struct {
	// custom logger object.
	logger *log.Entry

	// handler path
	path string

	// router group.
	group *gin.RouterGroup

	// middleware for handling user auth.
	auth *middlewares.AuthMiddleware

	// external services.
	users UserManager
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(auth *middlewares.AuthMiddleware, users UserManager) *AuthHandler {
	h := &AuthHandler{
		logger: log.WithFields(log.Fields{"package": "http", "module": "authHandler"}),
		path:   "/auth",
		auth:   auth,
		users:  users,
	}

	return h
}

// Bootstrap registers the handler routes in the server.
func (h *AuthHandler) Bootstrap(router *gin.Engine) {
	h.logger.Info("Bootstrapping auth handler")

	// register middleware.
	router.Use(h.auth.SetUserStatus())

	// auth routes.
	h.group = router.Group(h.path)
	h.group.POST(SignInRoute, h.performSignIn)
	h.group.POST(SignUpRoute, h.performSignUp)
	h.group.GET(SignOutRoute, h.performSignOut)
}

// performSignIn logs the user in.
func (h *AuthHandler) performSignIn(c *gin.Context) {
	// get the POSTed values.
	email := c.PostForm("email")
	password := c.PostForm("password")

	// get user from the database.
	u, err := h.users.FindByEmail(email)
	if err != nil {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "It seems we messed up somehow, please try again.",
		}, SignInPage)
		return
	}

	// check if the credentials are correct.
	if !checkPasswordHash(password, u.Password) {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "Invalid credentials provided",
		}, SignInPage)
		return
	}

	// log the user in.
	h.auth.AddSession(c, u.ID)

	// redirect to home page.
	util.Render(c, gin.H{
		"MessageTitle":   "Success!",
		"MessageMessage": "welcome back to treasure coin!",
	}, IndexPage)
}

// performSignOut renders the about page.
func (h *AuthHandler) performSignOut(c *gin.Context) {
	h.auth.RemoveSession(c)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, IndexRoute)
}

// performSignUp renders the about page.
func (h *AuthHandler) performSignUp(c *gin.Context) {
	// get the POSTed values.
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	// hash the supplied password.
	hash, err := hashPassword(password)
	if err != nil {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "It seems we messed up somehow, please try again.",
		}, SignUpPage)
		return
	}
	user := coin.User{
		Email:    email,
		Username: username,
		Password: hash,
	}

	// store the user data.
	user, err = h.users.Add(user)
	if err != nil {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "An account with that email already exists.",
		}, SignUpPage)
		return
	}

	// log the user in.
	h.auth.AddSession(c, user.ID)

	// redirect to home page.
	util.Render(c, gin.H{
		"MessageTitle":   "Success!",
		"MessageMessage": "welcome to treasure coin!",
	}, IndexPage)
}

// hashPassword generates an hash based on the supplied string.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// checkPasswordHash checks if the hash was created from the string.
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// UserManager defines the interface to interact with the user persistence layer.
type UserManager interface {
	Add(user coin.User) (coin.User, error)
	Find(id uint) (coin.User, error)
	FindByEmail(email string) (coin.User, error)
	FindByUsername(username string) (coin.User, error)
	Update(user coin.User) (coin.User, error)
	Delete(user coin.User) error
}
