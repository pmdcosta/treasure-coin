package database

import (
	"encoding/json"

	"github.com/pmdcosta/treasure-coin"
)

const UserCollection = "users"

// UserService represents a service for managing user persistence.
type UserService struct {
	client *Client
}

// Add adds the record to the database if it does not exist.
func (s *UserService) Add(user coin.User) error {
	j, _ := json.Marshal(user)
	return s.client.Create(UserCollection, user.Email, j)
}

// Find retrieves a user from the database.
func (s *UserService) Find(email string) (coin.User, error) {
	j, err := s.client.Load(UserCollection, email)
	if err != nil {
		return coin.User{}, err
	}

	u := coin.User{}
	json.Unmarshal(j, &u)

	return u, nil
}

// Save upserts the user to the database.
func (s *UserService) Save(user coin.User) error {
	j, _ := json.Marshal(user)
	return s.client.Save(UserCollection, user.Email, j)
}

// Remove deletes the user from the databse.
func (s *UserService) Remove(user coin.User) error {
	return s.client.Delete(UserCollection, user.Email)
}
