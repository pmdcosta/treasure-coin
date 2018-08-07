package ost

import (
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

	params := q + "&signature=" + hex.EncodeToString(h.Sum(nil))

	return params
}

// GetUserBalance retrieves the user balance from OST.
func (c *Client) GetUserBalance(user string) (string, error) {
	r := fmt.Sprintf("/users/%s/", user)
	t := fmt.Sprintf("%d", time.Now().Unix())
	q := fmt.Sprintf("api_key=%s&id=%s&request_timestamp=%s", c.apiKey, user, t)
	s := c.CreateSignature(r, q)
	u := c.url + r + "?" + s

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
