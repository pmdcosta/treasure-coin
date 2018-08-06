package database_test

import (
	"testing"

	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/database"
	"github.com/stretchr/testify/assert"
)

var testUser = coin.User{
	Email:    "test@user.com",
	Username: "test",
	Password: "hashedPassed",
}

// TestUserService_InsertRecord tests inserting a database record.
func TestUserService_InsertRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.UserService().Add(testUser)
	assert.Nil(t, err)
}

// TestUserService_InsertRecordExists tests inserting a database record that already exists.
func TestUserService_InsertRecordExists(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.UserService().Add(testUser)
	assert.Nil(t, err)

	err = c.UserService().Add(testUser)
	assert.Equal(t, database.ErrRecordExists, err)
}

// TestUserService_LoadRecord tests retrieving a database record.
func TestUserService_LoadRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.UserService().Add(testUser)
	assert.Nil(t, err)

	user, err := c.UserService().Find(testUser.Email)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)
}

// TestUserService_UpdateRecord tests updating a database record.
func TestUserService_UpdateRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.UserService().Add(testUser)
	assert.Nil(t, err)

	user, err := c.UserService().Find(testUser.Email)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)

	user.Password = "NewHashPassword"
	err = c.UserService().Save(user)
	assert.Nil(t, err)

	nUser, err := c.UserService().Find(testUser.Email)
	assert.Nil(t, err)
	assert.Equal(t, user, nUser)

}

// TestUserService_DeleteRecords tests removing records from a database testCollection.
func TestUserService_DeleteRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.UserService().Add(testUser)
	assert.Nil(t, err)

	user, err := c.UserService().Find(testUser.Email)
	assert.Nil(t, err)
	assert.Equal(t, testUser, user)

	err = c.UserService().Remove(testUser)
	assert.Nil(t, err)

	user, err = c.UserService().Find(testUser.Email)
	assert.Equal(t, err, database.ErrRecordNotFound)
	assert.Equal(t, coin.User{}, user)

}
