package ost_test

import (
	"encoding/json"
	"fmt"
	"github.com/pmdcosta/treasure-coin/ost"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Client is a test wrapper.
type Client struct {
	*ost.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)

	// get credentials from 'env' file.
	type Credentials struct {
		Key     string
		Secret  string
		Url     string
		Company string
	}
	file, err := os.Open("env")
	if err != nil {
		panic(err)
	}
	var cred Credentials
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cred)
	if err != nil {
		panic(err)
	}

	c := &Client{
		Client: ost.NewClient(cred.Url, cred.Key, cred.Secret, cred.Company),
	}
	return c
}

// TestClient_CreateSignature tests creating a signature.
func TestClient_CreateSignature(t *testing.T) {
	c := NewClient()
	sig := c.CreateSignature("/users/ID/", "api_key=KEY&id=ID&request_timestamp=TIME")
	assert.Equal(t, "28e9035850612343fdd46a38d5c35f451e0035680509572e56cd4f984987ebc9", sig)
}

// TestClient_GetUserBalance tests getting the user balance from OST.
func TestClient_GetUserBalance(t *testing.T) {
	c := NewClient()
	b, err := c.GetUserBalance("5190fed7-dbfb-4687-b2c8-b5cd57002198")
	assert.Nil(t, err)
	assert.Equal(t, "0", b)
}

// TestClient_CreateUser tests creating a new user using the API.
func TestClient_CreateUser(t *testing.T) {
	c := NewClient()
	u, err := c.CreateUser("Luffy")
	assert.Nil(t, err)
	fmt.Println(u)
}

// TestClient_GetUserTransactions tests getting user transactions using the API.
func TestClient_GetUserTransactions(t *testing.T) {
	c := NewClient()
	b, err := c.GetUserTransactions("87e9132d-0586-4beb-9600-ffa050966bc8")
	assert.Nil(t, err)
	assert.Equal(t, len(b), 1)
	assert.Equal(t, b[0].Amount, "0.1")
	assert.Equal(t, b[0].ActionId, 39879)
	assert.Equal(t, b[0].AirdropedAmount, "0")
	assert.Equal(t, b[0].FromUserID, "9bc1ee0d-084b-4d75-b798-fda6a270adcc")
	assert.Equal(t, b[0].ToUserID, "87e9132d-0586-4beb-9600-ffa050966bc8")
}

// TestClient_GetRewarded tests making a company to user transaction.
func TestClient_GetRewarded(t *testing.T) {
	c := NewClient()
	err := c.GetRewarded("5190fed7-dbfb-4687-b2c8-b5cd57002198")
	assert.Nil(t, err)
}

// TestClient_MakePayment tests making a user to company transaction.
func TestClient_MakePayment(t *testing.T) {
	c := NewClient()
	err := c.MakePayment("5190fed7-dbfb-4687-b2c8-b5cd57002198", 2)
	assert.Nil(t, err)
}
