package main

import (
	"github.com/namsral/flag"
	"github.com/pmdcosta/treasure-coin/database"
	"github.com/pmdcosta/treasure-coin/http"
	"github.com/pmdcosta/treasure-coin/http/handlers"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
)

func main() {
	var (
		dbPath = flag.String("db-path", "app.db", "Choose database path.")

		serverPort = flag.String("server-port", "8080", "Choose server port to bind to.")
	)
	flag.Parse()

	// instantiate the database client and services.
	db := database.NewClient(*dbPath)
	if err := db.Open(); err != nil {
		panic(err)
	}

	// instantiate the middlewares.
	am := middlewares.NewAuthMiddleware(db.UserService(), db.SessionService())

	// instantiate the handlers.
	dh := handlers.NewDefaultHandler(am, db.GameService())
	ah := handlers.NewAuthHandler(am, db.UserService())
	gh := handlers.NewGameHandler(am, db.GameService())

	// start the server.
	router := http.NewServer(":"+*serverPort, dh, ah, gh)
	if err := router.Open(); err != nil {
		panic(err)
	}
}
