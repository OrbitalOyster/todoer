package routes

import (
	"net/http"
	"strconv"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
)

func GetSingleTask(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "task", task)
}

func GetAllTasks(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	selectedTasks, _, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	jwt.Update(payload, "Page", page, writer)
	templates.ExecutePartial(writer, "task-table-body", selectedTasks)
}

func GetAddTaskForm(writer http.ResponseWriter, req *http.Request) {
	templates.ExecutePartial(writer, "addTaskForm", nil)
}

func GetEditTaskForm(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "editTaskForm", task)
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "cloneTaskForm", task)
}

func AddTask(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	user := payload.UserID
	description := req.FormValue("description")
	tasks.Add(user, description)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "New task", "Success")
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
	task := tasks.Check(idStr)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+idStr, "Success")
	templates.ExecutePartial(writer, "task-oob", task)
}
