package api

import (
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
	idStr, description := req.FormValue("taskId"), req.FormValue("taskDescription")
	if description == "bogus" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("Bogus description"))
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}
	tasks.Update(id, description)
	toasts.Success(writer, "Task "+idStr, "Success")
}
