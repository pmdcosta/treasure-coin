package database

import (
	"github.com/pmdcosta/treasure-coin"
)

// database errors.
const (
	ErrTransaction       = coin.Error("failed to start transaction")
	ErrCreateCollection  = coin.Error("failed to create collection")
	ErrRecordNotFound    = coin.Error("record does not exist")
	ErrRecordExists      = coin.Error("record already exists")
	ErrCreateRecord      = coin.Error("failed to insert record")
	ErrDeleteRecord      = coin.Error("failed to delete record")
	ErrIterateCollection = coin.Error("failed to iterate over collection")
	ErrCreateKey         = coin.Error("failed to generate a key for the  collection")
)
