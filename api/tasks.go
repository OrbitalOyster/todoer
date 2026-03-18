package api

import (
	"log"
	"net/http"
	"strconv"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
)

func Tasks(writer http.ResponseWriter, req *http.Request) {
	taskList := tasks.GetAll()
	data := struct {
		Tasks []tasks.Task
	}{
		Tasks: taskList,
	}
	templates.ExecutePartial(writer, "taskTable", data)
}

func EditTask(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.PathValue("id"))
	/* User sent stoopid */
	if err != nil {
		panic(err)
	}
	task, err := tasks.Get(id)
	/* No such task */
	if err != nil {
		panic(err)
	}
	data := struct {
		Id          int
		Description string
	}{
		Id:          id,
		Description: task.Description,
	}
	templates.ExecutePartial(writer, "editTaskModal", data)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	taskIdStr, taskDescription := req.FormValue("taskId"), req.FormValue("taskDescription")
	log.Printf("id: %s, desc: %s\n", taskIdStr, taskDescription)
	if taskDescription == "bogus" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bogus description"))
	} else {
		taskId, err := strconv.Atoi(taskIdStr)
		if err != nil {
			panic(err)
		}
		tasks.Update(taskId, taskDescription)
		toasts.Success(writer, "Task "+taskIdStr, "Success")
	}
}
