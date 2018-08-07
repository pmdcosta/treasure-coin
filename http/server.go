package http

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// templateDir describes the template asset directory.
const templateDir = "templates/*"

// staticDir describes the static asset directory.
const staticDir = "public/"

// Handler represents an http handler.
type Handler interface {
	Bootstrap(router *gin.Engine)
}

// Server represents an HTTP server.
type Server struct {
	// custom logger object.
	logger *log.Entry

	// http port to serve from.
	port   string
	cert   string
	secret string
	ssl    bool

	// router instance.
	router *gin.Engine

	// http handlers.
	handlers []Handler
}

// NewServer returns a new instance of Server.
func NewServer(port, cert, secret string, ssl bool, h ...Handler) *Server {
	// set the server to production mode.
	gin.SetMode(gin.ReleaseMode)

	s := &Server{
		router:   gin.Default(),
		logger:   log.WithFields(log.Fields{"package": "http", "module": "server"}),
		handlers: h,
		port:     port,
		cert:     cert,
		secret:   secret,
		ssl:      ssl,
	}
	return s
}

// Open starts the server.
func (c *Server) Open() error {
	// loads the templates from the disk.
	c.router.LoadHTMLGlob(templateDir)

	// serves the static assets.
	c.router.Use(static.Serve("/assets", static.LocalFile(staticDir, false)))

	// loads the http handlers of the project.
	for _, h := range c.handlers {
		h.Bootstrap(c.router)
	}

	// starts the http server.
	c.logger.Info("Starting server at port: " + c.port)

	if c.ssl {
		return c.router.RunTLS(c.port, c.cert, c.secret)
	} else {
		return c.router.Run(c.port)
	}
}
