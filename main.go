package main

import (
	"log"
	"todoer/api"
	"todoer/config"
	"todoer/routes"
	"todoer/server"
	"todoer/tasks"
	htmxTasks "todoer/tasks/htmx"
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
		"GET /{$}":                  routes.Main,
		"GET /login":                routes.Login,
		"POST /api/login":           api.LoginAttempt,
		"POST /api/logout":          api.Logout,
		"GET /htmx/tasks":           htmxTasks.GetAll,
		"GET /htmx/tasks/{id}":      htmxTasks.Get,
		"GET /htmx/edit-task/{id}":  htmxTasks.Edit,
		"GET /htmx/clone-task/{id}": htmxTasks.Clone,
		"PATCH /htmx/tasks":         htmxTasks.Patch,
		"PATCH /filters/page-size":  api.SetTaskTablePageSize,
		"GET /":                     routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
