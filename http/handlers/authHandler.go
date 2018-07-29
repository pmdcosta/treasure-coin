package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin/http/util"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
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
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(auth *middlewares.AuthMiddleware) *AuthHandler {
	h := &AuthHandler{
		logger: log.WithFields(log.Fields{"package": "http", "module": "authHandler"}),
		path: "/auth",
		auth: auth,
	}

	return h
}

// Bootstrap registers the handler routes in the server.
func (h *AuthHandler) Bootstrap(router *gin.Engine) {
	h.logger.Info("Bootstrapping auth handler")

	// register middleware.
	//router.Use(h.auth.SetUserStatus())

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

	// check if the credentials are valid.
	if email == "pmdcosta@outlook.com" && password == "password"{
		// set the token in a cookie.
		c.SetCookie("token", util.CreateSessionToken(), 3600, "", "", false, true)
		// update context with login status.
		c.Set("is_logged_in", true)
		// redirect to home page.
		util.Render(c, gin.H{
			"MessageTitle":   "Success!",
			"MessageMessage": "welcome back to treasure coin!",
		}, IndexPage)
	} else {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "Invalid credentials provided",
		}, SignInPage)
	}
}

// performSignOut renders the about page.
func (h *AuthHandler) performSignOut(c *gin.Context) {
	// clear the cookie.
	c.SetCookie("token", "", -1, "", "", false, true)

	// update context with login status.
	c.Set("is_logged_in", false)

	// Redirect to the home page
	c.Redirect(http.StatusTemporaryRedirect, IndexRoute)
}

// performSignUp renders the about page.
func (h *AuthHandler) performSignUp(c *gin.Context) {
	// get the POSTed values.
	email := c.PostForm("email")
	username := c.PostForm("username")
	password := c.PostForm("password")

	// store the user data!!

	// check if the values are valid.
	if email != "pmdcosta@outlook.com" && username != "" && password != "" {
		// set the token in a cookie.
		c.SetCookie("token", util.CreateSessionToken(), 3600, "", "", false, true)
		// update context with login status.
		c.Set("is_logged_in", true)
		// redirect to home page.
		util.Render(c, gin.H{
			"MessageTitle":   "Success!",
			"MessageMessage": "welcome to treasure coin!",
		}, IndexPage)
	} else {
		util.Render(c, gin.H{
			"ErrorTitle":   "Failed!",
			"ErrorMessage": "An account with that email already exists!",
		}, SignUpPage)
	}
}