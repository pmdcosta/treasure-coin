package database

import (
	"github.com/pmdcosta/treasure-coin"
)

// database errors.
const (
	ErrConnect    = coin.Error("failed to connect to the database")
	ErrMigrate    = coin.Error("failed to apply database migrations")
	ErrCreate     = coin.Error("failed to create record in the database")
	ErrRetrieving = coin.Error("failed to retrieve a record from the database")
	ErrUpdate     = coin.Error("failed to update a record in the database")
	ErrDelete     = coin.Error("failed to delete a record in the database")
)
