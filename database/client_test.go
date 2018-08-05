package database_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/pmdcosta/treasure-coin/database"
	log "github.com/sirupsen/logrus"
)

// Client is a test wrapper.
type Client struct {
	*database.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)
	c := &Client{
		Client: database.NewClient("localhost", "5433", "postgres", "postgres", "treasure_coin_test"),
	}
	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Client.Open(); err == nil {
		return c
	}

	// database must be missing, try to create it and try again.
	cmd := exec.Command("docker", "run", "-d", "--name", "treasure-coin-test-postgres", "-e", "POSTGRES_PASSWORD=postgres", "-e", "POSTGRES_DB=treasure_coin_test", "-p", "5433:5432", "postgres:11-alpine")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	if err := c.Client.Open(); err != nil {
		panic(err)
	}

	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	c.UserService().PurgeAll()
	return c.Client.Close()
}

// TestClient_Create tests creating a new torrent crawling client.
func TestClient_Create(t *testing.T) {
	c := MustOpenClient()
	if c == nil {
		t.Fatal("failed to create client")
	}
}
