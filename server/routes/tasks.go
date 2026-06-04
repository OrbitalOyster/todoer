package routes

import (
	"net/http"
	"strconv"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/server/token"
	"todoer/tasks"
	"todoer/utils"
)

func idCheck(writer http.ResponseWriter, req *http.Request) *tasks.Task {
	if task, err := tasks.GetById(req.PathValue("id")); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("Task not found: " + err.Error()))
		/* Major screwup */
		if err != nil {
			panic(err)
		}
		return nil
	} else {
		return task
	}
}

func GetSingleTask(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		pages.ExecutePartial(writer, "task", task)
	}
}

func GetAllTasks(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	token.Update(payload, "Page", page, writer)
	pages.ExecutePartial(writer, "task-list", TaskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Pagination: utils.GetPagination(totalPages, page),
		Payload:    token.Payload(*payload),
	})
}

func GetAddTaskForm(writer http.ResponseWriter, req *http.Request) {
	pages.ExecutePartial(writer, "addTaskForm", nil)
}

func GetEditTaskForm(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		pages.ExecutePartial(writer, "editTaskForm", task)
	}
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		pages.ExecutePartial(writer, "cloneTaskForm", task)
	}
}

func AddTask(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	user := payload.UserID
	description := req.FormValue("description")
	tasks.Add(user, description)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "New task", "Success")
	GetAllTasks(writer, req)
}

func PutTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task
	if task = idCheck(writer, req); task == nil {
		return
	}
	description, statusStr := req.FormValue("description"), req.FormValue("status")
	/* Status */
	status := false
	if statusStr == "on" {
		status = true
	}
	if task.Status != status {
		if err := task.SetStatus(status); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte("Unable to change task status:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}
	/* Description */
	if task.Description != description {
		if err := task.SetDescription(description); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte("Unable to set task description:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}
	/* Done */
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+strconv.Itoa(task.Id), "Success")
	GetAllTasks(writer, req)
}

func PatchTasks(writer http.ResponseWriter, req *http.Request) {
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task
	if task = idCheck(writer, req); task == nil {
		panic("Task not found")
	}
	field := req.PathValue("field")
	switch field {
	case "status":
		statusStr, status := req.FormValue("status"), false
		if statusStr == "on" || statusStr == "true" {
			status = true
		}
		if err := task.SetStatus(status); err != nil {
			_, err = writer.Write([]byte("Unable to change task status:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	default:
		_, err := writer.Write([]byte("Invalid task field: " + field))
		if err != nil {
			panic(err)
		}
	}
	toasts.Success(writer, "Task "+strconv.Itoa(task.Id), "Success")
	GetAllTasks(writer, req)
}

func DeleteTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task
	if task = idCheck(writer, req); task == nil {
		panic("Task not found")
	}
	taskId := task.Id
	tasks.Delete(taskId)
	toasts.Warning(writer, "Task "+strconv.Itoa(taskId)+" deleted", "Success")
	GetAllTasks(writer, req)
}

func DeleteTasks(writer http.ResponseWriter, req *http.Request) {
	deletedTasks := 0
	for id := range req.URL.Query() {
		task, err := tasks.GetById(id)
		if err != nil {
			panic(err)
		}
		if err := tasks.Delete(task.Id); err != nil {
			panic(err)
		}
		deletedTasks++
	}
	toasts.Warning(writer, "Deleted "+strconv.Itoa(deletedTasks)+" tasks", "Success")
	GetAllTasks(writer, req)
}
