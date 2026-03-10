package api

import (
	"log"
	"net/http"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
)

func Tasks(writer http.ResponseWriter, req *http.Request) {
	taskList := tasks.Get()
	data := struct {
		Tasks []tasks.Task
	} {
		Tasks: taskList,
	}
	templates.ExecutePartial(writer, "taskTable", data)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	taskId, taskDescription := req.FormValue("taskId"), req.FormValue("taskDescription")
	log.Printf("id: %s, desc: %s\n", taskId, taskDescription)
	if taskDescription == "bogus" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bogus description"))
	} else {
		toasts.Success(writer, "Task " + taskId, "Success")
	}
}
