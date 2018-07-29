package main

import (
	"github.com/pmdcosta/treasure-coin/http"
	"github.com/pmdcosta/treasure-coin/http/handlers"
	"github.com/pmdcosta/treasure-coin/http/middlewares"
)

func main() {
	// instantiate the middlewares.
	am := middlewares.NewAuthMiddleware()

	// instantiate the handlers.
	dh := handlers.NewDefaultHandler(am)
	ah := handlers.NewAuthHandler(am)



	// start the server.
	router := http.NewServer(":8080", dh, ah)
	if err := router.Open(); err != nil {
		panic(err)
	}
}
