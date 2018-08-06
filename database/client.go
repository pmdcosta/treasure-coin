package database

import (
	"time"

	"strconv"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// Client represents a client to managing the database.
type Client struct {
	logger *log.Entry

	// boltDB database.
	path string
	db   *bolt.DB

	// object services.
	userService    UserService
	gameService    GameService
	sessionService SessionService
}

// NewClient returns a new configuration client.
func NewClient(path string) *Client {
	c := &Client{
		logger: log.WithFields(log.Fields{"package": "database"}),
		path:   path,
	}
	c.userService.client = c
	c.gameService.client = c
	c.sessionService.client = c
	return c
}

// Open starts client handler.
func (c *Client) Open() error {
	// open database file.
	db, err := bolt.Open(c.path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	c.db = db

	// initialize top-level buckets.
	tx, err := c.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	return tx.Commit()
}

// Close terminates client.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Create persists the supplied data as a new record.
func (c *Client) Create(collection string, key string, value []byte) error {
	// start read-write transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return err
	}

	// check if the key exists.
	v := b.Get([]byte(key))
	if v != nil {
		c.logger.WithFields(log.Fields{"collection": collection, "record": key}).Debug(ErrRecordExists)
		return ErrRecordExists
	}

	// insert record.
	err = b.Put([]byte(key), value)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection, "record": key}).Error(ErrCreateRecord)
		return err
	}

	c.logger.WithFields(log.Fields{"collection": collection, "key": key, "record": string(value)}).Debug("record created")
	return tx.Commit()
}

// CreateIndexed persists the supplied data as a new record and creates a new record ID.
func (c *Client) CreateIndexed(collection string, value []byte) (string, error) {
	// start read-write transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return "", err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return "", err
	}

	// get a auto-generated key.
	id, err := b.NextSequence()
	if err != nil {
		c.logger.WithFields(log.Fields{"collection": collection}).Debug(ErrCreateKey)
		return "", ErrCreateKey
	}
	key := strconv.FormatUint(id, 10)

	// insert record.
	err = b.Put([]byte(key), value)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection, "record": key}).Error(ErrCreateRecord)
		return "", err
	}

	c.logger.WithFields(log.Fields{"collection": collection, "key": key, "record": string(value)}).Debug("record created")
	return key, tx.Commit()
}

// Load retrieves the stored data.
func (c *Client) Load(collection string, key string) ([]byte, error) {
	// start read transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return nil, err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return nil, err
	}

	// find record.
	v := b.Get([]byte(key))
	if v == nil {
		c.logger.WithFields(log.Fields{"collection": collection, "record": key}).Debug(ErrRecordNotFound)
		return nil, ErrRecordNotFound
	}

	c.logger.WithFields(log.Fields{"collection": collection, "key": key, "record": string(v)}).Debug("record loaded")
	return v, nil
}

// Save persists the supplied data.
func (c *Client) Save(collection string, key string, value []byte) error {
	// start read-write transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return err
	}

	// insert record.
	err = b.Put([]byte(key), value)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection, "record": key}).Error(ErrCreateRecord)
		return err
	}

	c.logger.WithFields(log.Fields{"collection": collection, "key": key, "record": string(value)}).Debug("record created")
	return tx.Commit()
}

// Iterate iterates over all the keys in a bucket.
func (c *Client) Iterate(collection string, executer func(k, v []byte) error) error {
	// start read-write transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return err
	}

	// iterate over all records in bucket.
	if err := b.ForEach(executer); err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrIterateCollection)
		return err
	}
	return tx.Commit()
}

// Delete removes a key from the database.
func (c *Client) Delete(collection string, keys ...string) error {
	// start read-write transaction.
	tx, err := c.db.Begin(true)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(ErrTransaction)
		return err
	}
	defer tx.Rollback()

	// create collection if it does not exist.
	b, err := tx.CreateBucketIfNotExists([]byte(collection))
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "collection": collection}).Error(ErrCreateCollection)
		return err
	}

	// delete records.
	for _, k := range keys {
		if err = b.Delete([]byte(k)); err != nil {
			c.logger.WithFields(log.Fields{"error": err, "collection": collection, "record": k}).Debug(ErrDeleteRecord)
		}
		c.logger.WithFields(log.Fields{"collection": collection, "record": k}).Debug("deleted record")
	}
	return tx.Commit()
}

// UserService returns the service used to manage user persistence.
func (c *Client) UserService() *UserService { return &c.userService }

// GameService returns the service used to manage game persistence.
func (c *Client) GameService() *GameService { return &c.gameService }

// SessionService returns the service used to manage game persistence.
func (c *Client) SessionService() *SessionService { return &c.sessionService }
