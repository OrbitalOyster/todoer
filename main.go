package main

import (
	"log"
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
		"GET /{$}":                       routes.Main,
		"GET /login":                     routes.Login,
		"POST /login":                    routes.LoginAttempt,
		"POST /logout":                   routes.Logout,
		"GET /htmx/tasks":                routes.GetAllTasks,
		"GET /htmx/tasks/{id}":           routes.GetSingleTask,
		"GET /htmx/edit-task/{id}":       routes.GetEditTaskForm,
		"GET /htmx/clone-task/{id}":      routes.GetCloneTaskForm,
		"PATCH /htmx/tasks":              routes.PatchTask,
		"PATCH /filters/page-size":       routes.SetPageSize,
		"PATCH /filters/page/{page}":     routes.SetPage,
		"PATCH /filters/next-page":       routes.NextPage,
		"PATCH /filters/previous-page":   routes.PreviousPage,
		"PATCH /filters/sortBy/{column}": routes.SetSortBy,
		"PATCH /filters/searchBy":        routes.SetSearchBy,
		"PATCH /filters/fromDate":        routes.SetFromDate,
		"PATCH /filters/toDate":          routes.SetToDate,
		"GET /":                          routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
