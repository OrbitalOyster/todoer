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
	/* Pages */
	templates.AddPage("login", "login")
	templates.AddPage("main", "base")
	/* Routes */
	routerMap := map[string]routes.RouterEntry{
		"GET /{$}":                       routes.GetMainPage,
		"GET /login":                     routes.GetLoginPage,
		"POST /login":                    routes.LoginAttempt,
		"POST /logout":                   routes.Logout,
		"GET /htmx/tasks":                routes.GetAllTasks,
		"GET /htmx/tasks/{id}":           routes.GetSingleTask,
		"GET /htmx/edit-task/{id}":       routes.GetEditTaskForm,
		"GET /htmx/add-task":             routes.GetAddTaskForm,
		"GET /htmx/clone-task/{id}":      routes.GetCloneTaskForm,
		"POST /htmx/tasks":               routes.AddTask,
		"PATCH /htmx/tasks":              routes.PatchTask,
		"DELETE /htmx/tasks/{id}":        routes.DeleteTask,
		"PATCH /filters/page-size":       routes.SetPageSize,
		"PATCH /filters/page/{page}":     routes.SetPage,
		"PATCH /filters/next-page":       routes.NextPage,
		"PATCH /filters/previous-page":   routes.PreviousPage,
		"PATCH /filters/sortBy/{column}": routes.SetSortBy,
		"PATCH /filters/searchBy":        routes.SetSearchBy,
		"PATCH /filters/fromDate":        routes.SetFromDate,
		"PATCH /filters/toDate":          routes.SetToDate,
		"GET /panic":                     routes.Panic,
		"GET /":                          routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
