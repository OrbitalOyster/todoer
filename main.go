package main

import (
	"log"
	"todoer/config"
	"todoer/server"
	"todoer/server/pages"
	"todoer/server/routes"
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
	/* Pages */
	pages.Add("login", "login")
	pages.Add("main", "base")
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
		"PUT /htmx/tasks/{id}":           routes.PutTask,
		"PATCH /htmx/tasks/{id}/{field}": routes.PatchTask,
		"DELETE /htmx/tasks":             routes.DeleteTasks,
		"DELETE /htmx/tasks/{id}":        routes.DeleteTask,
		"PATCH /filters/page-size":       routes.SetPageSize,
		"PATCH /filters/page/{page}":     routes.SetPage,
		"PATCH /filters/next-page":       routes.NextPage,
		"PATCH /filters/previous-page":   routes.PreviousPage,
		"PATCH /filters/sortBy/{column}": routes.SetSortBy,
		"PATCH /filters/searchBy":        routes.SetSearchBy,
		"PATCH /filters/date":            routes.SetDate,
		"GET /panic":                     routes.Panic,
		"GET /":                          routes.NotFoundHandler, // 404 page
	}
	server.Start(routerMap)
}
