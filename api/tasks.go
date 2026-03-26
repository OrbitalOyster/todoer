package api

import (
	"net/http"
	"strconv"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
)

func checkTask(idStr string) tasks.Task {
	id, err := strconv.Atoi(idStr)
	/* User sent stoopid */
	if err != nil {
		panic(err)
	}
	task, err := tasks.Get(id)
	/* No such task */
	if err != nil {
		panic(err)
	}
	return task
}

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
	task := checkTask(req.PathValue("id"))
	templates.ExecutePartial(writer, "editTaskForm", task)
}

func CloneTask(writer http.ResponseWriter, req *http.Request) {
	task := checkTask(req.PathValue("id"))
	templates.ExecutePartial(writer, "cloneTaskForm", task)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	idStr, description := req.FormValue("id"), req.FormValue("description")
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

	writer.Header().Set("HX-Trigger", "hideModal")

	toasts.Success(writer, "Task "+idStr, "Success")
}
