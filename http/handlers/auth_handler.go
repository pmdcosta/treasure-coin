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
	users   UserManager
	games   GameManager
	wallets WalletService
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(auth *middlewares.AuthMiddleware, users UserManager, games GameManager, wallets WalletService) *AuthHandler {
	h := &AuthHandler{
		logger:  log.WithFields(log.Fields{"package": "http", "module": "authHandler"}),
		path:    "/auth",
		auth:    auth,
		users:   users,
		games:   games,
		wallets: wallets,
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
	u, err := h.users.Find(email)
	if err != nil {
		h.logger.WithFields(log.Fields{"email": email}).Error(err)
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "It seems we messed up somehow, please try again.",
		}.Render(), SignInPage)
		return
	}

	// check if the credentials are correct.
	if !checkPasswordHash(password, u.Password) {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Invalid credentials provided.",
		}.Render(), SignInPage)
		return
	}

	// log the user in.
	h.auth.AddSession(c, u.Email)

	// redirect to home page.
	games := h.games.List()
	util.Render(c, gin.H{
		"games":          games,
		"MessageTitle":   "Success",
		"MessageMessage": "Welcome back to treasure coin " + u.Username + ".",
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
		h.logger.Error(err)
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "It seems we messed up somehow, please try again.",
		}.Render(), SignUpPage)
		return
	}

	// create a user wallet.
	w, err := h.wallets.CreateUser(username)
	if err != nil || w == "" {
		h.logger.WithFields(log.Fields{"username": username, "step": "wallet"}).Error(err)
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "It seems we messed up somehow, please try again.",
		}.Render(), SignUpPage)
		return
	}

	// airdrop the users some tokens.
	if err := h.wallets.Airdrop(w, 1.0); err != nil {
		h.logger.WithFields(log.Fields{"wallet": w, "step": "airdrop"}).Error(err)
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "It seems we messed up somehow, please try again.",
		}.Render(), SignUpPage)
		return
	}

	// build the user.
	user := coin.User{
		Email:    email,
		Username: username,
		Password: hash,
		Wallet:   w,
	}

	// store the user data.
	err = h.users.Add(user)
	if err != nil {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "An account with that email already exists.",
		}.Render(), SignUpPage)
		return
	}

	// log the user in.
	h.auth.AddSession(c, user.Email)

	// redirect to home page.
	games := h.games.List()
	util.Render(c, gin.H{
		"games":          games,
		"MessageTitle":   "Success",
		"MessageMessage": "Welcome to treasure coin " + user.Username + ".",
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
	Add(user coin.User) error
	Find(email string) (coin.User, error)
	FindByWallet(wallet string) coin.User
}

// WalletService defines the interface to interact with the blockchain wallet layer.
type WalletService interface {
	CreateUser(user string) (string, error)
	GetUserBalance(user string) (string, error)
	Airdrop(user string, amount float64) error
	GetRewarded(user string) error
	MakePayment(user string, amount int) error
	GetUserTransactions(user string) ([]coin.Transaction, error)
}
