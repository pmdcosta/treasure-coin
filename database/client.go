package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// Client represents a client to managing the database.
type Client struct {
	logger *log.Entry

	// database config.
	db       *gorm.DB
	host     string
	port     string
	user     string
	password string
	database string

	// application services for using the database.
	userService UserService
}

// NewClient returns a new instance of Client.
func NewClient(host, port, user, password, database string) *Client {
	c := &Client{
		logger:   log.WithFields(log.Fields{"package": "database"}),
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
	}
	c.userService.client = c
	return c
}

// Open starts the client.
func (c *Client) Open() error {
	// open database.
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.host, c.port, c.user, c.database, c.password))
	if err != nil {
		c.logger.WithFields(log.Fields{"err": err}).Error(ErrConnect)
		return ErrConnect
	}
	db.LogMode(false)
	c.db = db

	// create tables if not exists.
	if dbc := db.AutoMigrate(&User{}); dbc.Error != nil {
		c.logger.WithFields(log.Fields{"err": dbc.Error}).Error(ErrMigrate)
		return ErrMigrate
	}

	return nil
}

// Close stops the client.
func (c *Client) Close() error {
	return c.db.Close()
}

// UserService returns the service used to manage user persistence.
func (c *Client) UserService() *UserService { return &c.userService }
