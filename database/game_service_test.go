package database_test

import (
	"testing"
	"time"

	"github.com/pmdcosta/treasure-coin"
	"github.com/pmdcosta/treasure-coin/database"
	"github.com/stretchr/testify/assert"
)

// default test game.
var testGame = coin.Game{
	Title:     "Pirate Golden Age",
	StartDate: time.Time{},
	Creator:   "gol@d.roger",
	Treasures: map[string]coin.Treasure{
		"treasure-1": {
			ID:        "treasure_1",
			Name:      "One Piece",
			Hint:      "Poneglyphs",
			Location:  "Raftel",
			QRCode:    "treasure_1.jpg",
			Token:     "D",
			Found:     false,
			FoundDate: time.Time{},
			FoundUser: "",
		},
	},
}

// default game ID.
var testGameID = "1"

// TestGameService_InsertRecord tests inserting a database record.
func TestGameService_InsertRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	key, err := c.GameService().Add(testGame)
	assert.Nil(t, err)
	assert.Equal(t, testGameID, key)
}

// TestGameService_LoadRecord tests retrieving a database record.
func TestGameService_LoadRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	key, err := c.GameService().Add(testGame)
	assert.Nil(t, err)
	assert.Equal(t, testGameID, key)

	game, err := c.GameService().Find(testGameID)
	assert.Nil(t, err)
	assert.Equal(t, testGame, game)
}

// TestGameService_DeleteRecords tests removing records from a database testCollection.
func TestGameService_DeleteRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	key, err := c.GameService().Add(testGame)
	assert.Nil(t, err)
	assert.Equal(t, testGameID, key)

	game, err := c.GameService().Find(testGameID)
	assert.Nil(t, err)
	assert.Equal(t, testGame, game)

	err = c.GameService().Remove(testGame)
	assert.Nil(t, err)

	game, err = c.GameService().Find(testGame.ID)
	assert.Equal(t, err, database.ErrRecordNotFound)
	assert.Equal(t, coin.Game{}, game)
}

// TestGameService_ListRecords tests listing all the records in the database.
func TestGameService_ListRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()

	key1, err := c.GameService().Add(testGame)
	assert.Nil(t, err)
	assert.Equal(t, testGameID, key1)

	newGame := testGame
	newGame.Title = "NewTestGame"
	key2, err := c.GameService().Add(newGame)
	assert.Nil(t, err)
	assert.Equal(t, "2", key2)

	games := c.GameService().List()
	assert.Equal(t, map[string]coin.Game{"1": testGame, "2": newGame}, games)
}
