package routes

import (
	"net/http"
	"strconv"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/toasts"
	"todoer/utils"
)

func GetSingleTask(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	templates.ExecutePartial(writer, "task", task)
}

func GetAllTasks(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	jwt.Update(payload, "Page", page, writer)
	templates.ExecutePartial(writer, "task-list", TaskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Pagination: utils.GetPagination(totalPages, page),
		Payload:    jwt.Payload(*payload),
	})
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
	// GetAllTasks(writer, req)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	idStr, description := req.FormValue("id"), req.FormValue("description")
	if description == "bogus" {
		writer.WriteHeader(http.StatusBadRequest)
		_, err := writer.Write([]byte("Bogus description"))
		if err != nil {
			panic(err)
		}
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err)
	}
	err = tasks.Update(id, description)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("Unable to update task:" + err.Error()))
		if err != nil {
			panic(err)
		}
	}
	task := tasks.Check(idStr)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+idStr, "Success")
	templates.ExecutePartial(writer, "task-oob", task)
}
