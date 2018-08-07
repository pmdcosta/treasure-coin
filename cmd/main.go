package main

import (
	"github.com/namsral/flag"
	"github.com/pmdcosta/treasure-coin/database"
	"github.com/pmdcosta/treasure-coin/http"
	"github.com/pmdcosta/treasure-coin/http/handlers"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
	"github.com/pmdcosta/treasure-coin/ost"
)

func main() {
	var (
		dbPath       = flag.String("db-path", "app.db", "Choose database path.")
		serverHost   = flag.String("server-host", "https://treasurecoin.powertrip.pt", "Choose server host.")
		serverPort   = flag.String("server-port", "8080", "Choose server port to bind to.")
		serverCert   = flag.String("server-cert", "ssl/certificate.pem", "Choose server certificate for ssl.")
		serverSecret = flag.String("server-secret", "ssl/secret.pem", "Choose server secret for ssl.")
		serverSSL    = flag.Bool("server-ssl", false, "Choose wheather the server should use ssl.")
		ostUrl       = flag.String("ost-url", "", "Choose the OST API base url.")
		ostKey       = flag.String("ost-key", "", "Choose the OST API key.")
		ostSecret    = flag.String("ost-secret", "", "Choose the OST API secret.")
		ostCompany   = flag.String("ost-company", "", "Choose the OST API company ID.")
	)
	flag.Parse()

	// get ost config.
	config := ost.Config{}
	config.LoadCred(".env", *ostUrl, *ostKey, *ostSecret, *ostCompany)

	// instantiate the database client and services.
	db := database.NewClient(*dbPath)
	if err := db.Open(); err != nil {
		panic(err)
	}

	// instantiate the ost client service.
	st := ost.NewClient(config)

	// instantiate the middleware.
	am := middlewares.NewAuthMiddleware(db.UserService(), db.SessionService())

	// instantiate the handlers.
	dh := handlers.NewDefaultHandler(am, db.GameService(), db.UserService(), st)
	ah := handlers.NewAuthHandler(am, db.UserService(), db.GameService(), st)
	gh := handlers.NewGameHandler(am, db.GameService(), st, *serverHost)

	// start the server.
	router := http.NewServer(":"+*serverPort, *serverCert, *serverSecret, *serverSSL, dh, ah, gh)
	if err := router.Open(); err != nil {
		panic(err)
	}
}
