package htmxTasks

import (
	"todoer/tasks"
	"todoer/templates"
	"net/http"
)

func Get(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "task", task)
}

func GetAll(writer http.ResponseWriter, req *http.Request) {
	data := struct {
		Tasks []tasks.Task
	}{
		Tasks: tasks.GetAll(),
	}
	templates.ExecutePartial(writer, "taskTable", data)
}

func Edit(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "editTaskForm", task)
}

func Clone(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "cloneTaskForm", task)
}
