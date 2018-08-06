package database

import (
	"encoding/json"

	"github.com/pmdcosta/treasure-coin"
)

const GameCollection = "games"

// GameService represents a service for managing game persistence.
type GameService struct {
	client *Client
}

// Add stores the game in the database.
func (s *GameService) Add(game coin.Game) (string, error) {
	j, _ := json.Marshal(game)
	return s.client.CreateIndexed(GameCollection, j)
}

// Find retrieves a game from the database.
func (s *GameService) Find(id string) (coin.Game, error) {
	j, err := s.client.Load(GameCollection, id)
	if err != nil {
		return coin.Game{}, err
	}

	var g coin.Game
	json.Unmarshal(j, &g)

	return g, nil
}

// Save upserts the game to the database.
func (s *GameService) Save(game coin.Game) error {
	j, _ := json.Marshal(game)
	return s.client.Save(GameCollection, game.ID, j)
}

// Remove removes the game from the database.
func (s *GameService) Remove(game coin.Game) error {
	return s.client.Delete(GameCollection, game.ID)
}

// List returns all the games from the database.
func (s *GameService) List() map[string]coin.Game {
	games := make(map[string]coin.Game)
	s.client.Iterate(GameCollection, func(k, v []byte) error {
		var g coin.Game
		json.Unmarshal(v, &g)

		games[string(k)] = g
		return nil
	})

	return games
}
