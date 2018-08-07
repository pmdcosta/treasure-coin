package ost_test

import (
	"fmt"
	"github.com/pmdcosta/treasure-coin/ost"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Client is a test wrapper.
type Client struct {
	*ost.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)

	cred := ost.Config{}
	cred.LoadCred("../.env", "", "", "", "")

	c := &Client{
		Client: ost.NewClient(cred),
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
	assert.Equal(t, b[0].FromWallet, "9bc1ee0d-084b-4d75-b798-fda6a270adcc")
	assert.Equal(t, b[0].ToWallet, "87e9132d-0586-4beb-9600-ffa050966bc8")
	fmt.Println(b[0].Date)
}

// TestClient_Airdrop tests incrementing user's balance the API.
func TestClient_Airdrop(t *testing.T) {
	c := NewClient()
	err := c.Airdrop("1bc46b40-2d76-4bfa-a806-b9ce1983ae8f", 0.1)
	assert.Nil(t, err)
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

// TestClient_RemoveTokens tests removing all tokens from a client transaction.
func TestClient_RemoveTokens(t *testing.T) {
	c := NewClient()
	err := c.RemoveTokens("5190fed7-dbfb-4687-b2c8-b5cd57002198")
	assert.Nil(t, err)
}
