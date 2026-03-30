package htmxTasks

import (
	"net/http"
	"strconv"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
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

func Patch(writer http.ResponseWriter, req *http.Request) {
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
	task := tasks.Check(idStr)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+idStr, "Success")
	templates.ExecutePartial(writer, "task-oob", task)
}
