package main

import (
	"log"
	"todoer/api"
	"todoer/config"
	"todoer/routes"
	"todoer/server"
	"todoer/tasks"
	"todoer/templates"
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
	/* Templates */
	templates.Add("login", "base.html", "login.html")
	templates.Add("main", "base.html", "main.html")
	/* Routes */
	routerMap := map[string]server.RouterEntry{
		"GET /{$}":    routes.Main,
		"GET /login":  routes.Login,
		"POST /login": api.LoginAttempt,
		"GET /logout": api.Logout,
		"GET /":       routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
