package main

import (
	"log"
	"todoer/api"
	"todoer/config"
	"todoer/routes"
	"todoer/server"
)

func main() {
	/* Error handler */
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Println("Crashed:", recovered)
		}
		log.Println("Bye")
	}()
	config.Load()
	/* Routes */
	routerMap := map[string] server.RouterEntry {
		"GET /": routes.DefaultHandler,
		"GET /login": routes.Login,
		"POST /login": api.LoginAttempt,
	}
	server.Start(routerMap)
}
