package main

import (
	"log"
	"todoer/api"
	"todoer/config"
	"todoer/routes"
	"todoer/server"
	"todoer/tasks"
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
	tasks.Load()
	/* Routes */
	routerMap := map[string]server.RouterEntry{
		"GET /{$}":    routes.Main,
		"GET /login":  routes.Login,
		"GET /foo":    routes.Foo,
		"POST /login": api.LoginAttempt,
		"GET /logout": api.Logout,
		"GET /":       routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
