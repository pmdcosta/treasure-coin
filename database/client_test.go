package database_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/database"
	"github.com/pmdcosta/yaad"
	log "github.com/sirupsen/logrus"
)

// Client is a test wrapper.
type Client struct {
	*database.Client
}

const path = "/tmp/app.db"

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)
	c := &Client{
		Client: database.NewClient(path),
	}
	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Client.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	c.Client.Close()
	return os.Remove(path)
}

// TestClient_Create tests creating a new torrent crawling client.
func TestClient_Create(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("failed to create client")
	}
}

const testCollection = "test"

var testObj = coin.User{
	Email:    "test@treasure.coin",
	Username: "test",
	Password: "hashedPassed",
}

// TestClient_InsertRecord tests inserting a database record.
func TestClient_InsertRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	r, err := json.Marshal(testObj)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Save(testCollection, testObj.Email, r)
	if err != nil {
		t.Fatal(err)
	}
}

// TestClient_LoadRecord tests retrieving a database record.
func TestClient_LoadRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	r, err := json.Marshal(testObj)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Save(testCollection, testObj.Email, r)
	if err != nil {
		t.Fatal(err)
	}

	record, err := c.Load(testCollection, testObj.Email)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(r, record) {
		t.Fatal(err)
	}
}

// TestClient_IterateCollection tests iterating over a database testCollection.
func TestClient_IterateCollection(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	// insert struct 1.
	r1, err := json.Marshal(testObj)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Save(testCollection, testObj.Email+"_1", r1)
	if err != nil {
		t.Fatal(err)
	}

	// insert struct 2.
	r2, err := json.Marshal(testObj)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Save(testCollection, testObj.Email+"_2", r2)
	if err != nil {
		t.Fatal(err)
	}

	iterated := []bool{false, false}
	// iterate over structc.
	c.Iterate(testCollection, func(k, v []byte) error {
		var anime yaad.Anime
		if err = json.Unmarshal(v, &anime); err != nil {
			t.Fatal(err)
		}

		if string(k) == testObj.Email+"_1" {
			iterated[0] = true
		} else if string(k) == testObj.Email+"_2" {
			iterated[1] = true
		}

		return nil
	})

	for k, v := range iterated {
		if !v {
			t.Fatalf("failed to iterate over %d", k)
		}
	}
}

// TestClient_DeleteRecords tests removing records from a database testCollection.
func TestClient_DeleteRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	r, err := json.Marshal(testObj)
	if err != nil {
		t.Fatal(err)
	}

	// save record.
	err = c.Save(testCollection, testObj.Email, r)
	if err != nil {
		t.Fatal(err)
	}

	// delete record.
	err = c.Delete(testCollection, testObj.Email)
	if err != nil {
		t.Fatal(err)
	}

	// load record.
	_, err = c.Load(testCollection, testObj.Email)
	if err != database.ErrRecordNotFound {
		t.Fatal(err)
	}
}

// TestClient_LoadRecord_NoRecord tests retrieving a database record that does not exist.
func TestClient_LoadRecord_NoRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	_, err := c.Load(testCollection, "FAKE")
	if err != database.ErrRecordNotFound {
		t.Fatal(err)
	}
}
