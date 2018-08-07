package ost

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client represents a client to interact with the OST api.
type Client struct {
	logger *log.Entry

	url       string
	apiKey    string
	apiSecret string
	companyID string
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
func NewClient(url, apiKey, apiSecret, companyID string) *Client {
	c := &Client{
		logger:    log.WithFields(log.Fields{"package": "ost"}),
		url:       url,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		companyID: companyID,
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
	// build the request.
	t := fmt.Sprintf("%d", time.Now().Unix())
	r := fmt.Sprintf("/users/%s/", user)
	query := map[string]string{
		"request_timestamp": t,
		"api_key":           c.apiKey,
		"id":                user,
	}
	u, err := c.BuildRequest(c.url, r, query)
	if err != nil {
		return "", err
	}

	// make the request.
	response, err := http.Get(u.String())
	if err != nil {
		return "", err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("Invalid status code received: " + response.Status + " | " + string(contents))
	}

	// unmarhal the response.
	type GetUserResponse struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				Balance string `json:"token_balance"`
			} `json:"user"`
		} `json:"data"`
	}
	var resp GetUserResponse
	json.Unmarshal(contents, &resp)

	return resp.Data.User.Balance, nil
}

// CreateUser creates a new user account in the OST platform.
func (c *Client) CreateUser(user string) (string, error) {
	// build the request.
	t := fmt.Sprintf("%d", time.Now().Unix())
	r := "/users/"
	query := map[string]string{
		"request_timestamp": t,
		"api_key":           c.apiKey,
		"name":              user,
	}
	u, err := c.BuildRequest(c.url, r, query)
	if err != nil {
		return "", err
	}

	// make the request.
	response, err := http.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(u.RawQuery)))
	if err != nil {
		return "", err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("Invalid status code received: " + response.Status + " | " + string(contents))
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
	// build the request.
	t := fmt.Sprintf("%d", time.Now().Unix())
	r := fmt.Sprintf("/ledger/%s/", user)
	query := map[string]string{
		"request_timestamp": t,
		"api_key":           c.apiKey,
		"page_no":           "1",
	}
	u, err := c.BuildRequest(c.url, r, query)
	if err != nil {
		return []Transaction{}, err
	}

	// make the request.
	response, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return []Transaction{}, errors.New("Invalid status code received: " + response.Status + " | " + string(contents))
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

// Airdrop adds TreasureCoins to a user's balance from OST.
func (c *Client) Airdrop(user string, amount float64) error {
	r := fmt.Sprintf("/airdrops")
	t := fmt.Sprintf("%d", time.Now().Unix())
	q := fmt.Sprintf("amount=%f&api_key=%s&request_timestamp=%s&user_ids=%s", amount, c.apiKey, t, user)
	s := c.CreateSignature(r, q)
	fq := q + "&signature=" + s
	u := c.url + r + "?" + s

	// make the request.
	response, err := http.Post(u, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(fq)))
	if err != nil {
		return err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	type GetUserResponse struct {
		Success bool `json:"success"`
		Error   struct {
			Msg string `json:"msg"`
		} `json:"err"`
	}

	var resp GetUserResponse
	json.Unmarshal(contents, &resp)

	if !resp.Success {
		return errors.New(resp.Error.Msg)
	}
	return nil
}

// GetRewarded makes a company-to-user transaction request to OST.
func (c *Client) GetRewarded(user string) error {
	// build the request.
	t := fmt.Sprintf("%d", time.Now().Unix())
	r := "/transactions/"
	query := map[string]string{
		"request_timestamp": t,
		"api_key":           c.apiKey,
		"action_id":         "39879",
		"from_user_id":      c.companyID,
		"to_user_id":        user,
	}
	u, err := c.BuildRequest(c.url, r, query)
	if err != nil {
		return err
	}

	// make the request.
	response, err := http.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(u.RawQuery)))
	if err != nil {
		return err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Invalid status code received: " + response.Status + " | " + string(contents))
	}

	return nil
}

// GetPayed makes a user-to-company transaction request to OST.
func (c *Client) MakePayment(user string, amount int) error {
	tokens := fmt.Sprintf("%f", float32(amount)*0.1)

	// build the request.
	t := fmt.Sprintf("%d", time.Now().Unix())
	r := "/transactions/"
	query := map[string]string{
		"request_timestamp": t,
		"api_key":           c.apiKey,
		"from_user_id":      user,
		"to_user_id":        c.companyID,
		"action_id":         "39876",
		"amount":            tokens,
		"currency":          "BT",
	}
	u, err := c.BuildRequest(c.url, r, query)
	if err != nil {
		return err
	}

	// make the request.
	response, err := http.Post(u.String(), "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(u.RawQuery)))
	if err != nil {
		return err
	}

	// parse the response.
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Invalid status code received: " + response.Status + " | " + string(contents))
	}

	fmt.Println(string(contents))

	return nil
}

// BuildRequest builds the OST request params.
func (c *Client) BuildRequest(host string, resource string, query map[string]string) (*url.URL, error) {
	// build url.
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	// add resource path.
	u.Path = u.Path + resource

	// add query parameters.
	q := u.Query()
	for k, v := range query {
		q.Add(k, v)
	}

	// add signature.
	q.Add("signature", c.CreateSignature(resource, q.Encode()))
	u.RawQuery = q.Encode()

	return u, nil
}
