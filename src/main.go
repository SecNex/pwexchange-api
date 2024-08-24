package main

import (
	"log"
	"os"

	"github.com/secnex/pwexchange/api/middlewares"
	"github.com/secnex/pwexchange/api/routes"
	"github.com/secnex/pwexchange/api/server"
	"github.com/secnex/pwexchange/storage"
	"github.com/secnex/pwexchange/utils"
)

func main() {
	var AUTH_TOKEN string
	AUTH_TOKEN = os.Getenv("AUTH_TOKEN")
	if AUTH_TOKEN == "" {
		AUTH_TOKEN = utils.NewRandom(32).String()
		log.Printf("No AUTH_TOKEN provided, using random token %s\n", AUTH_TOKEN)
	}

	var SERVER_SECRET string
	SERVER_SECRET = os.Getenv("SERVER_SECRET")
	if SERVER_SECRET == "" {
		SERVER_SECRET = utils.NewRandom(32).String()
		log.Printf("No SERVER_SECRET provided, using random secret %s\n", SERVER_SECRET)
	}
	serverSecret := []byte(SERVER_SECRET)
	vault := storage.NewVault(serverSecret)

	go vault.ExpirationRoutine()

	auth := middlewares.NewAuthentication(AUTH_TOKEN)

	server := server.NewServer("1", "pwexchange", "0.0.0.0", 8080, "/api", vault, nil, auth)
	server.AddRoute(routes.NewRoute("store/encrypt", vault.EndpointEncrypt))
	server.AddRoute(routes.NewRoute("store/decrypt", vault.EndpointDecrypt))
	server.AddRoute(routes.NewRoute("store", vault.Endpoint))
	server.RunServer()
}
