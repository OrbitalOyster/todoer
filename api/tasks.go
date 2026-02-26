package api

import (
	"net/http"
	"todoer/tasks"
	"todoer/templates"
)

func Tasks(writer http.ResponseWriter, req *http.Request)  {
	taskList := tasks.Get()
	data := struct {
		Tasks []tasks.Task
	} {
		Tasks: taskList,
	}
	templates.ExecutePartial(writer, "taskTable", data)
}
