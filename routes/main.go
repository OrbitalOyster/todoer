package routes

import (
	"net/http"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	payload, _ := jwt.Get(req)
	taskList := tasks.GetAll("", 10, 0)
	data := struct {
		Title string
		Username string
		Tasks []tasks.Task
	} {
		Title: "todoer",
		Username: payload.UserID,
		Tasks: taskList,
	}
	templates.Execute(writer, "main", data)
}
