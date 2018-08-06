package database_test

import (
	"testing"

	"github.com/pmdcosta/treasure-coin/database"
	"github.com/stretchr/testify/assert"
)

var testToken = "token"
var testSession = "session"

// TestSessionService_InsertRecord tests inserting a database record.
func TestSessionService_InsertRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.SessionService().Add(testToken, testSession)
	assert.Nil(t, err)
}

// TestSessionService_LoadRecord tests retrieving a database record.
func TestSessionService_LoadRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.SessionService().Add(testToken, testSession)
	assert.Nil(t, err)

	session, err := c.SessionService().Find(testToken)
	assert.Nil(t, err)
	assert.Equal(t, testSession, session)
}

// TestSessionService_DeleteRecords tests removing records from a database testCollection.
func TestSessionService_DeleteRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	err := c.SessionService().Add(testToken, testSession)
	assert.Nil(t, err)

	session, err := c.SessionService().Find(testToken)
	assert.Nil(t, err)
	assert.Equal(t, testSession, session)

	err = c.SessionService().Remove(testToken)
	assert.Nil(t, err)

	session, err = c.SessionService().Find(testToken)
	assert.Equal(t, err, database.ErrRecordNotFound)
	assert.Equal(t, "", session)

}
