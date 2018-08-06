package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
	"github.com/pmdcosta/treasure-coin/http/util"
	log "github.com/sirupsen/logrus"
)

// DefaultHandler handles miscellaneous pages in the server.
type DefaultHandler struct {
	// custom logger object.
	logger *log.Entry

	// handler path
	path string

	// router group.
	group *gin.RouterGroup

	// middleware for handling user auth.
	auth *middlewares.AuthMiddleware

	// external services.
	games GameManager
}

// NewDefaultHandler returns a new instance of DefaultHandler.
func NewDefaultHandler(auth *middlewares.AuthMiddleware, games GameManager) *DefaultHandler {
	h := &DefaultHandler{
		logger: log.WithFields(log.Fields{"package": "http", "module": "default-handler"}),
		path:   "/",
		auth:   auth,
		games:  games,
	}

	return h
}

// Bootstrap registers the handler routes in the server.
func (h *DefaultHandler) Bootstrap(router *gin.Engine) {
	h.logger.Info("Bootstrapping default handler")

	// register middleware.
	router.Use(h.auth.SetUserStatus())

	// default routes.
	h.group = router.Group(h.path)
	h.group.GET(IndexRoute, h.showIndexPage)
	h.group.GET(AboutRoute, h.showAboutPage)
	h.group.GET(ProfileRoute, h.showProfilePage)
	h.group.GET(SignInRoute, h.showSignInPage)
	h.group.GET(SignUpRoute, h.showSignUpPage)
}

// showIndexPage renders the about page.
func (h *DefaultHandler) showIndexPage(c *gin.Context) {
	games := h.games.List()
	util.Render(c, gin.H{
		"games": games,
	}, IndexPage)
}

// showAboutPage renders the about page.
func (h *DefaultHandler) showAboutPage(c *gin.Context) {
	util.Render(c, gin.H{}, AboutPage)
}

// showProfilePage renders the profile page.
func (h *DefaultHandler) showProfilePage(c *gin.Context) {
	util.Render(c, gin.H{}, ProfilePage)
}

// showSignInPage renders the about page.
func (h *DefaultHandler) showSignInPage(c *gin.Context) {
	util.Render(c, gin.H{}, SignInPage)
}

// showSignUpPage renders the about page.
func (h *DefaultHandler) showSignUpPage(c *gin.Context) {
	util.Render(c, gin.H{}, SignUpPage)
}
