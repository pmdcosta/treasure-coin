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
		dbHost     = flag.String("db-host", "localhost", "Choose database host.")
		dbPort     = flag.String("db-port", "5432", "Choose database port.")
		dbUser     = flag.String("db-user", "postgres", "Choose database user.")
		dbPassword = flag.String("db-pwd", "postgres", "Choose database password.")
		dbDatabase = flag.String("db-database", "treasure_coin", "Choose database.")

		serverPort = flag.String("server-port", "8080", "Choose server port to bind to.")
	)
	flag.Parse()

	// instantiate the database client and services.
	db := database.NewClient(*dbHost, *dbPort, *dbUser, *dbPassword, *dbDatabase)
	if err := db.Open(); err != nil {
		panic(err)
	}

	// instantiate the middlewares.
	am := middlewares.NewAuthMiddleware(db.UserService())

	// instantiate the handlers.
	dh := handlers.NewDefaultHandler(am)
	ah := handlers.NewAuthHandler(am, db.UserService())

	// start the server.
	router := http.NewServer(":"+*serverPort, dh, ah)
	if err := router.Open(); err != nil {
		panic(err)
	}
}
