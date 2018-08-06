package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
	"github.com/pmdcosta/treasure-coin/http/util"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
)

// GameHandler handles game related pages in the server.
type GameHandler struct {
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

// NewGameHandler returns a new instance of GameHandler.
func NewGameHandler(auth *middlewares.AuthMiddleware, games GameManager) *GameHandler {
	h := &GameHandler{
		logger: log.WithFields(log.Fields{"package": "http", "module": "game-handler"}),
		path:   "/games",
		auth:   auth,
		games:  games,
	}

	return h
}

// Bootstrap registers the handler routes in the server.
func (h *GameHandler) Bootstrap(router *gin.Engine) {
	h.logger.Info("Bootstrapping game handler")

	// register middleware.
	router.Use(h.auth.SetUserStatus())

	// default routes.
	h.group = router.Group(h.path)
	h.group.GET(CreateGameRoute, h.showCreatePage)
	h.group.GET(ListGameRoute, h.showListPage)
	h.group.GET(DescribeGameRoute, h.showDescribePage)
	h.group.POST(CreateGameRoute, h.performCreateGame)
	h.group.GET(DescribeTreasureRoute, h.showDescribeTreasurePage)
	h.group.GET(FoundTreasureRoute, h.performFoundTreasure)
}

// showCreatePage renders the create game page.
func (h *GameHandler) showCreatePage(c *gin.Context) {
	logedin, exists := c.Get(util.LogInCookie)
	if !exists || !logedin.(bool) {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "SignIn/SignUp to create a new game.",
		}.Render(), SignInPage)
		return
	}

	util.Render(c, gin.H{}, CreateGamePage)
}

// showListPage renders the list of available games page.
func (h *GameHandler) showListPage(c *gin.Context) {
	games := h.games.List()
	util.Render(c, gin.H{
		"games": games,
	}, ListGamePage)
}

// showDescribePage renders the describe game page.
func (h *GameHandler) showDescribePage(c *gin.Context) {
	user, exists := c.Get(util.UserCookie)
	if !exists {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Requires a logged in user.",
		}.Render(), IndexPage)
		return
	}

	g := c.Param("game")

	game, err := h.games.Find(g)
	if err != nil {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Game not found.",
		}.Render(), IndexPage)
	}
	game.ID = g

	util.Render(c, gin.H{
		"game": game,
		"user": user.(coin.User),
	}, DescribeGamePage)
}

// performCreateGame creates a new game.
func (h *GameHandler) performCreateGame(c *gin.Context) {
	user, exists := c.Get(util.UserCookie)
	if !exists {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Requires a logged in user.",
		}.Render(), IndexPage)
		return
	}

	// validate request.
	r := createGameRequest{}
	if err := r.Validate(c); err != nil {
		util.Render(c, err.Render(), CreateGamePage)
		return
	}

	// build game data.
	g := coin.Game{
		Title:       r.title,
		Description: r.description,
		StartDate:   time.Now().Truncate(time.Second),
		Creator:     user.(coin.User).Email,
		Treasures:   make(map[string]coin.Treasure),
	}

	// build treasure data.
	for _, t := range r.treasures {
		// create token.
		token, err := uuid.NewV4()
		if err != nil {
			h.logger.Info("failed to generate uuid")
			util.Render(c, util.RequestError{
				Title:   "Failed!",
				Message: "Failed to create game, please try again.",
			}.Render(), CreateGamePage)
			return
		}

		// create treasure.
		treasure := coin.Treasure{
			ID:       t.id,
			Name:     t.name,
			Hint:     t.hint,
			Location: t.location,
			Token:    token.String(),
		}
		g.Treasures[t.id] = treasure
	}

	gameID, err := h.games.Add(g)
	if err != nil {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Failed to create game, please try again.",
		}.Render(), CreateGamePage)
		return
	}

	// create qr codes for tokens.
	for _, t := range g.Treasures {
		discoveryUrl := fmt.Sprintf("http://192.168.122.1:8080/games/found/%s/%s?token=%s", gameID, t.ID, t.Token)
		codeFile := fmt.Sprintf("%s-%s.png", gameID, t.ID)
		err := qrcode.WriteFile(discoveryUrl, qrcode.Medium, 256, fmt.Sprintf("%s/%s", "public/codes", codeFile))
		if err != nil {
			h.logger.Error("failed to generate QR code file.")
		}
		t.QRCode = codeFile
		g.Treasures[t.ID] = t
	}

	// save game.
	g.ID = gameID
	h.games.Save(g)

	util.Render(c, gin.H{
		"MessageTitle":   "Success!",
		"MessageMessage": "The game has been created, check the treasures for the QR code to hide!",
		"game":           g,
		"user":           user.(coin.User),
	}, DescribeGamePage)
}

// showDescribeTreasurePage renders the describe treasure page.
func (h *GameHandler) showDescribeTreasurePage(c *gin.Context) {
	user, exists := c.Get(util.UserCookie)
	if !exists {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Requires a logged in user.",
		}.Render(), IndexPage)
		return
	}

	g := c.Param("game")
	t := c.Param("treasure")

	game, err := h.games.Find(g)
	if err != nil {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Game not found.",
		}.Render(), IndexPage)
	}

	treasure, ok := game.Treasures[t]
	if !ok {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Treasure not found.",
		}.Render(), IndexPage)
	}

	util.Render(c, gin.H{
		"game":     game,
		"treasure": treasure,
		"user":     user.(coin.User),
	}, DescribeTreasurePage)
}

// performFoundTreasure sets a treasure as found.
func (h *GameHandler) performFoundTreasure(c *gin.Context) {
	user, exists := c.Get(util.UserCookie)
	if !exists {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Requires a logged in user.",
		}.Render(), IndexPage)
		return
	}

	g := c.Param("game")
	t := c.Param("treasure")
	token := c.Query("token")

	// get game.
	game, err := h.games.Find(g)
	if err != nil {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Game not found.",
		}.Render(), IndexPage)
	}

	// get treasure.
	treasure, ok := game.Treasures[t]
	if !ok {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Treasure not found.",
		}.Render(), IndexPage)
	}

	// check if token is correct.
	if treasure.Token != token {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "Incorrect treasure token.",
		}.Render(), IndexPage)
	}

	// check if treasure was already found.
	if treasure.Found {
		util.Render(c, util.RequestError{
			Title:   "Failed!",
			Message: "This treasure was already found.",
		}.Render(), IndexPage)
	}

	// set the treasure as found.
	treasure.Found = true
	treasure.FoundUser = user.(coin.User).Email
	treasure.FoundDate = time.Now()

	game.ID = g
	game.Treasures[t] = treasure
	h.games.Save(game)

	util.Render(c, gin.H{
		"game":           game,
		"treasure":       treasure,
		"MessageTitle":   "Congratulations!",
		"MessageMessage": "You have found a lost treasure!",
	}, DescribeTreasurePage)
}

/**
 * Requests
 */

// CreateGameRequest represents the form data from a performCreateGame request.
type createGameRequest struct {
	title       string
	description string
	nTreasures  string
	treasures   []treasureRequest
}

// TreasureRequest represents the internal treasure representation of the CreateGameRequest.
type treasureRequest struct {
	id       string
	name     string
	location string
	hint     string
}

// validate validates a CreateGameRequest request.
func (r *createGameRequest) Validate(c *gin.Context) *util.RequestError {
	r.title = c.PostForm("title")
	r.description = c.PostForm("description")
	r.nTreasures = c.PostForm("treasures")

	// validate title.
	if r.title == "" {
		return &util.RequestError{
			Title:   "Failed!",
			Message: "Please provide a valid game title.",
		}
	}

	// validate description.
	if r.description == "" {
		return &util.RequestError{
			Title:   "Failed!",
			Message: "Please provide a valid game description.",
		}
	}

	// get number of treasures.
	n, err := strconv.Atoi(r.nTreasures)
	if err != nil || n == 0 {
		return &util.RequestError{
			Title:   "Failed!",
			Message: "Please provide a valid number of treasures.",
		}
	}

	// validate treasures.

	for i := 0; i < n; i++ {
		t := treasureRequest{
			name:     c.PostForm(fmt.Sprintf("treasure-name-%v", i)),
			location: c.PostForm(fmt.Sprintf("treasure-location-%v", i)),
			hint:     c.PostForm(fmt.Sprintf("treasure-hint-%v", i)),
		}
		t.id = slug.Make(t.name)
		if t.id == "" {
			return &util.RequestError{
				Title:   "Failed!",
				Message: "Please provide a valid treasure name",
			}
		}
		if t.name == "" {
			return &util.RequestError{
				Title:   "Failed!",
				Message: "Please provide a valid treasure name",
			}
		}
		if t.location == "" {
			return &util.RequestError{
				Title:   "Failed!",
				Message: "Please provide a valid treasure location",
			}
		}
		if t.hint == "" {
			return &util.RequestError{
				Title:   "Failed!",
				Message: "Please provide a valid treasure hint",
			}
		}

		r.treasures = append(r.treasures, t)
	}

	// check if there are multiple treasures with the same id.
	ids := make(map[string]bool)
	for _, t := range r.treasures {
		if ids[t.id] {
			return &util.RequestError{
				Title:   "Failed!",
				Message: "Please provide treasures with different names",
			}
		}
		ids[t.id] = true
	}

	return nil
}

// GameManager defines the interface to interact with the game persistence layer.
type GameManager interface {
	Add(game coin.Game) (string, error)
	Find(id string) (coin.Game, error)
	Save(game coin.Game) error
	Remove(game coin.Game) error
	List() map[string]coin.Game
}
