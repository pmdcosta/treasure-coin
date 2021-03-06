package coin

import "time"

// User represents the domain user structure.
type User struct {
	Email    string
	Username string
	Password string
	Wallet   string
}

// Game represents the domain game structure.
type Game struct {
	ID          string
	Title       string
	Description string
	StartDate   time.Time
	Creator     string
	Treasures   map[string]Treasure
}

// Treasure represents the domain treasure structure.
type Treasure struct {
	ID        string
	Name      string
	Hint      string
	Location  string
	QRCode    string
	Token     string
	Found     bool
	FoundDate time.Time
	FoundUser string
}

// Transaction represents the domain coin transfer event.
type Transaction struct {
	FromWallet string
	ToWallet   string
	Event      string
	Date       time.Time
	Amount     string
}
