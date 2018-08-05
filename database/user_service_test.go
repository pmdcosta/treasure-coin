package database_test

import (
	"testing"

	"github.com/pmdcosta/treasure-coin"
	"github.com/stretchr/testify/assert"
)

var testUser = coin.User{
	ID:       1,
	Email:    "test@example.com",
	Username: "testUser",
	Password: "password",
}

// TestUserService_CreateUser tests persisting a new user to the database.
func TestUserService_CreateUser(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)
}

// TestUserService_Find tests looking up a user by their ID.
func TestUserService_Find(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)

	// find user.
	user, err := s.Find(testUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)
}

// TestUserService_FindByEmail tests looking up a user by their email.
func TestUserService_FindByEmail(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)

	// find user.
	user, err := s.FindByEmail(testUser.Email)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)
}

// TestUserService_FindByEmail tests looking up a user by their username.
func TestUserService_FindByUsername(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)

	// find user.
	user, err := s.FindByUsername(testUser.Username)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)
}

// TestUserService_Update tests updating a persisted user.
func TestUserService_Update(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)

	newUser := testUser
	newUser.Username = "newUser"

	// update user.
	user, err := s.Update(newUser)
	assert.Nil(t, err)
	assert.Equal(t, newUser, user)

	// check if user was update.
	u, err = s.Find(testUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, newUser, u)
}

// TestUserService_Delete tests deleting a persisted user.
func TestUserService_Delete(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.UserService()

	// insert user.
	u, err := s.Add(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, u)

	// delete user.
	err = s.Delete(testUser)
	assert.Nil(t, err)

	// check if user was update.
	u, err = s.Find(testUser.ID)
	assert.NotNil(t, err)
	assert.Equal(t, coin.User{}, u)
}
