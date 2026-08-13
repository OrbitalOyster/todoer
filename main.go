package main

import (
	"log"
	"todoer/config"
	"todoer/server"
	"todoer/server/pages"
	"todoer/server/routes"
	"todoer/tasks"
	"todoer/users"
	"todoer/wad"
	"todoer/collection"
)

const wadFilename = "wad.yaml"

func main() {
	/* Error handler */
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Println("Crashed:", recovered)
		}
		log.Println("Bye")
	}()
	config.Load()
	userList, taskList := wad.Load(wadFilename)
	users.Load(userList)
	tasks.Load(taskList)
	/* Pages */
	pages.Add("login", "login")
	pages.Add("main", "base")
	pages.Add("users", "base")
	pages.Add("categories", "base")
	/* Routes */
	routerMap := map[string]routes.RouterEntry{
		/* "/{$}" exactly matches root path ("/") */
		"GET /{$}":                       routes.GetMainPage,
		"GET /users":                     routes.GetUsersPage,
		"GET /categories":                routes.GetCategoriesPage,
		"GET /login":                     routes.GetLoginPage,
		"POST /login":                    routes.LoginAttempt,
		"POST /logout":                   routes.Logout,
		"GET /hx/tasks":                  routes.GetAllTasks,
		"GET /hx/tasks/{id}":             routes.GetSingleTask,
		"GET /hx/edit-task/{id}":         routes.GetEditTaskForm,
		"GET /hx/add-task":               routes.GetAddTaskForm,
		"GET /hx/clone-task/{id}":        routes.GetCloneTaskForm,
		"POST /hx/tasks":                 routes.AddTask,
		"PUT /hx/tasks/{id}":             routes.PutTask,
		"PATCH /hx/tasks/{id}/{field}":   routes.PatchTask,
		"PATCH /hx/tasks":                routes.PatchTasks,
		"DELETE /hx/tasks/{id}":          routes.DeleteTask,
		"DELETE /hx/tasks":               routes.DeleteTasks,
		"PATCH /filters/page-size":       routes.SetPageSize,
		"PATCH /filters/page/{page}":     routes.SetPage,
		"PATCH /filters/next-page":       routes.NextPage,
		"PATCH /filters/previous-page":   routes.PreviousPage,
		"PATCH /filters/sort-by/{field}": routes.SetSortBy,
		"PATCH /filters/searchBy":        routes.SetSearchBy,
		"PATCH /filters/date":            routes.SetDate,
		"GET /panic":                     routes.Panic,
		"GET /":                          routes.NotFoundHandler, /* 404 goes here */
	}
	collection.Run()
	server.Start(routerMap)
}
