package ost

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// Client represents a client to interact with the OST api.
type Client struct {
	logger *log.Entry

	url       string
	apiKey    string
	apiSecret string
}

type Transaction struct {
	Id               string `json:"id"`
	FromUserID       string `json:"from_user_id"`
	ToUserID         string `json:"to_user_id"`
	TransactionHash  string `json:"transaction_hash"`
	ActionId         int    `json:"action_id"`
	TimeStamp        int    `json:"timestamp"`
	Status           string `json:"status"`
	GasPrice         string `json:"gas_price"`
	GasUsed          string `json:"gas_used"`
	TransactionFee   string `json:"transaction_fee"`
	BlockNumber      int    `json:"block_number"`
	Amount           string `json:"amount"`
	CommissionAmount string `json:"commission_amount"`
	AirdropedAmount  string `json:"airdropped_amount"`
}

// NewClient returns a new configuration client.
func NewClient(url, apiKey, apiSecret string) *Client {
	c := &Client{
		logger:    log.WithFields(log.Fields{"package": "ost"}),
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
	return c
}

// CreateSignature creates the signature for the request.
func (c *Client) CreateSignature(resource, q string) string {
	key := []byte(c.apiSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(resource + "?" + q))

	return hex.EncodeToString(h.Sum(nil))
}

// GetUserBalance retrieves the user balance from OST.
func (c *Client) GetUserBalance(user string) (string, error) {
	r := fmt.Sprintf("/users/%s/", user)
	t := fmt.Sprintf("%d", time.Now().Unix())
	q := fmt.Sprintf("api_key=%s&id=%s&request_timestamp=%s", c.apiKey, user, t)
	s := c.CreateSignature(r, q)
	u := c.url + r + "?" + q + "&signature=" + s

	// make the request.
	response, err := http.Get(u)
	if err != nil {
		return "", err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// response struct.
	type GetUserResponse struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				Balance string `json:"token_balance"`
			} `json:"user"`
		} `json:"data"`
	}

	// unmarhal the response.
	var resp GetUserResponse
	json.Unmarshal(contents, &resp)

	return resp.Data.User.Balance, nil
}

// CreateUser creates a new user account in the OST platform.
func (c *Client) CreateUser(user string) (string, error) {
	r := fmt.Sprintf("/users/")
	t := fmt.Sprintf("%d", time.Now().Unix())
	q := fmt.Sprintf("api_key=%s&name=%s&request_timestamp=%s", c.apiKey, user, t)
	s := c.CreateSignature(r, q)
	fq := q + "&signature=" + s
	u := c.url + r + "?" + fq

	// make the request.
	response, err := http.Post(u, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(fq)))
	if err != nil {
		return "", err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// response struct.
	type CreateUserResponse struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				ID string `json:"id"`
			} `json:"user"`
		} `json:"data"`
	}

	// unmarhal the response.
	var resp CreateUserResponse
	json.Unmarshal(contents, &resp)

	return resp.Data.User.ID, nil
}

// GetUserTransactions retrieves the last 10 transactions from OST.
func (c *Client) GetUserTransactions(user string) ([]Transaction, error) {
	r := fmt.Sprintf("/ledger/%s/", user)
	t := fmt.Sprintf("%d", time.Now().Unix())
	q := fmt.Sprintf("api_key=%s&page_no=1&request_timestamp=%s", c.apiKey, t)
	s := c.CreateSignature(r, q)
	u := c.url + r + "?" + q + "&signature=" + s

	// make the request.
	response, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// response struct.
	type GetUserResponse struct {
		Success bool `json:"success"`
		Data    struct {
			Transactions []Transaction `json:"transactions"`
		} `json:"data"`
	}

	// unmarhal the response.
	var resp GetUserResponse
	json.Unmarshal(contents, &resp)

	return resp.Data.Transactions, nil
}
