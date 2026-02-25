package routes

import (
	"net/http"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	claims := jwt.Get(req)
	taskList := tasks.Get()
	data := struct {
		Title string
		Username string
		Tasks []tasks.Task
	} {
		Title: "todoer",
		Username: claims.UserID,
		Tasks: taskList,
	}
	templates.Execute(writer, "main", data)
}
