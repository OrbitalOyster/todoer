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

func GetSingleTask(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	pages.ExecutePartial(writer, "task", task)
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
	task := tasks.Check(req.PathValue("id"))
	pages.ExecutePartial(writer, "editTaskForm", task)
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	pages.ExecutePartial(writer, "cloneTaskForm", task)
}

func AddTask(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	user := payload.UserID
	description := req.FormValue("description")
	tasks.Add(user, description)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "New task", "Success")
	/* Send updated task list TODO: This chunck repeats */
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

func PutTask(writer http.ResponseWriter, req *http.Request) {
	description, doneStr := req.FormValue("description"), req.FormValue("done")
	task := tasks.Check(req.PathValue("id"))
	/* Status */
	done := false
	if doneStr == "on" {
		done = true
	}
	if task.Done != done {
		_, err := tasks.SetStatus(task.Id, done)
		if err != nil {
			_, err = writer.Write([]byte("Unable to change task status:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}
	/* Description */
	_, err := tasks.SetDescription(task.Id, description)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_, err = writer.Write([]byte("Unable to set task description:" + err.Error()))
		if err != nil {
			panic(err)
		}
	}
	/* Done */
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+strconv.Itoa(task.Id), "Success")

	payload := token.Get(req)
	/* Send updated task list */
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

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	task, field := tasks.Check(req.PathValue("id")), req.PathValue("field")
	switch field {
	case "done":
		doneStr, done := req.FormValue("done"), false
		if doneStr == "on" || doneStr == "true" {
			done = true
		}
		_, err := tasks.SetStatus(task.Id, done)
		if err != nil {
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

	payload := token.Get(req)
	/* Send updated task list */
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

func DeleteTask(writer http.ResponseWriter, req *http.Request) {
	task := tasks.Check(req.PathValue("id"))
	tasks.Delete(task.Id)
	toasts.Success(writer, "Task "+strconv.Itoa(task.Id), "Success")

	payload := token.Get(req)
	/* Send updated task list */
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
