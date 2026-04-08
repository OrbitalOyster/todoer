package routes

import (
	"net/http"
	"todoer/config"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	payload, _ := jwt.Get(req)
	taskList := tasks.GetAll("", payload.PageSize, 0)
	data := struct {
		Title            string
		Username         string
		Tasks            []tasks.Task
		PageSizes        []int
		SelectedPageSize int
	}{
		Title:            "todoer",
		Username:         payload.UserID,
		Tasks:            taskList,
		PageSizes:        config.PageSizes,
		SelectedPageSize: payload.PageSize,
	}
	templates.Execute(writer, "main", data)
}
