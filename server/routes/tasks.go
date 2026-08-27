package routes

import (
	"net/http"
	"slices"
	"strconv"
	"time"
	"todoer/collection"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/server/token"
	"todoer/tasks"
	"todoer/users"
	"todoer/utils"
)

func idCheck(writer http.ResponseWriter, req *http.Request) *tasks.Task[tasks.TaskFieldName] {
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

const defaultPageSize = 10

type TaskListDataNew struct {
	Tasks      []tasks.Task[tasks.TaskFieldName]
	Page       uint
	PageSize   int
	TotalPages int
	Pagination []int
	SearchBy   string
	FromDate   time.Time
	ToDate     time.Time
	SortBy     string
	SortDesc   bool
	Checkboxes []bool
}

func getTasks(query TasksQuery[tasks.TaskFieldName]) ([]tasks.Task[tasks.TaskFieldName], uint, int) {
	result := tasks.GetAll()
	result = result.
		MoreThan(
			tasks.Datetime,
			tasks.Task[tasks.TaskFieldName]{Datetime: query.FromDate},
		).
		/* "Not after 20/03/2026" means "Not after 20/03/2026 23:59:59"  */
		LessThan(
			tasks.Datetime,
			tasks.Task[tasks.TaskFieldName]{Datetime: query.ToDate.Add(time.Hour*24 - time.Second)},
		).
		Filter(tasks.Description, query.SearchBy).
		SortBy(query.SortBy)
	if query.SortDesc {
		result.Reverse()
	}
	result, page, numberOfPages := result.GetPage(uint(query.Page), uint(query.Size))
	selectedTasks := collection.AssertType[tasks.Task[tasks.TaskFieldName]](result)
	return selectedTasks, page, int(numberOfPages)
}

func GetTasksPage(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context().
		Value("token").(*token.Token[utils.Payload]).
		GetPayload()

	query, redirect := CreateQueryFromRequest(req)

	/* Update URL */
	if redirect {
		queryStr := query.String()
		if len(queryStr) > 0 {
			queryStr = "?" + queryStr
		}
		http.Redirect(writer, req, "/tasks" + queryStr, http.StatusSeeOther)
		return
	}

	selectedTasks, page, totalPages := getTasks(query)
	checkboxedTasks := getCheckboxedTasks(req)
	checkboxes := make([]bool, len(selectedTasks))
	for i, selectedTask := range selectedTasks {
		if slices.Contains(checkboxedTasks, selectedTask.Id) {
			checkboxes[i] = true
		} else {
			checkboxes[i] = false
		}
	}

	pages.Execute(writer, "tasks", struct {
		Title     string
		Payload   utils.Payload
		PageSizes []int
		TaskListDataNew
	}{
		Title:     "todoer - tasks",
		Payload:   payload,
		PageSizes: config.PageSizes,
		TaskListDataNew: TaskListDataNew{
			Tasks:      selectedTasks,
			Page:       page,
			PageSize:   query.Size,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(int(totalPages), int(page)),
			SearchBy:   query.SearchBy,
			FromDate:   query.FromDate,
			ToDate:     query.ToDate,
			SortBy:     query.SortBy.String(),
			SortDesc:   query.SortDesc,
			Checkboxes: checkboxes,
		},
	})
}

func GetSingleTask(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		pages.ExecutePartial(writer, "task", task)
	}
}

func getCheckboxedTasks(req *http.Request) (result []int) {
	if !req.Form.Has("checked") { /* Nothing checked */
		return result
	}
	for _, checkboxStr := range req.Form["checked"] {
		n, err := strconv.Atoi(checkboxStr)
		if err != nil {
			panic(err) /* Unparseable string */
		}
		result = append(result, n)
	}
	return result
}

func GetTaskList(writer http.ResponseWriter, req *http.Request) {
	query, push := CreateQueryFromRequest(req)

	/* Get tasks */
	tasksOnCurrentPage, page, numberOfPages := getTasks(query)
	if page != uint(query.Page) {
		query.Page = int(page)
	}

	checkboxedTasks := getCheckboxedTasks(req)
	checkboxes := make([]bool, len(tasksOnCurrentPage))
	for i, selectedTask := range tasksOnCurrentPage {
		if slices.Contains(checkboxedTasks, selectedTask.Id) {
			checkboxes[i] = true
		} else {
			checkboxes[i] = false
		}
	}

	/* Update calendar elements if both dates are set */
	if req.Form.Has("from") && req.Form.Has("to") {
		pages.ExecutePartial(
			writer,
			"task-dates-oob-new",
			struct {
				FromDate time.Time
				ToDate   time.Time
			}{
				FromDate: query.FromDate,
				ToDate:   query.ToDate,
			},
		)
	}

	/* Update URL */
	if push {
		queryStr := query.String()
		if len(queryStr) > 0 {
			queryStr = "?" + queryStr
		}
		writer.Header().Add("HX-Push-Url", "/tasks"+queryStr)
	}

	/* Send actual list */
	pages.ExecutePartial(writer, "task-list-new", TaskListDataNew{
		Tasks:      tasksOnCurrentPage,
		Page:       page,
		PageSize:   query.Size,
		TotalPages: int(numberOfPages),
		Pagination: utils.GetPagination(int(numberOfPages), int(page)),
		SortBy:     query.SortBy.String(),
		SortDesc:   query.SortDesc,
		Checkboxes: checkboxes,
	})
}

func GetAddTaskForm(writer http.ResponseWriter, req *http.Request) {
	pages.ExecutePartial(writer, "addTaskForm", nil)
}

func GetEditTaskForm(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		data := EditTaskFormData{
			task,
			users.GetAllUsers(),
		}
		pages.ExecutePartial(writer, "editTaskForm", data)
	}
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	if task := idCheck(writer, req); task != nil {
		pages.ExecutePartial(writer, "cloneTaskForm", task)
	}
}

func AddTask(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()                                  /* Load actual payload */
	user := payload.UserID
	description := req.FormValue("description")
	tasks.Add(user, description)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "New task", "Success")
	// GetAllTasks(writer, req)
	GetTaskList(writer, req)
}

func PutTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task[tasks.TaskFieldName]
	if task = idCheck(writer, req); task == nil {
		return
	}
	description, user, readOnlyStr :=
		req.FormValue("description"),
		req.FormValue("user"),
		req.FormValue("read-only")
	readOnly := false
	if readOnlyStr == "true" {
		readOnly = true
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
	/* User */
	if task.User != user {
		if err := task.SetUser(user); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte("Unable to change user:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}
	/* Read only */
	if task.ReadOnly != readOnly {
		if err := task.SetReadOnly(readOnly); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, err = writer.Write([]byte("Unable to change task:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	}

	/* Done */
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "Task "+strconv.Itoa(task.Id), "Success")
	// GetAllTasks(writer, req)
	GetTaskList(writer, req)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task[tasks.TaskFieldName]
	if task = idCheck(writer, req); task == nil {
		panic("Task not found")
	}
	field := req.PathValue("field")
	switch field {
	case "status":
		status := tasks.ParseStatus(req.FormValue("status"))
		if err := task.SetStatus(status); err != nil {
			_, err = writer.Write([]byte("Unable to change task status:" + err.Error()))
			if err != nil {
				panic(err)
			}
		}
	case "read-only":
		readOnly := false
		readOnlyStr := req.FormValue("read-only")
		if readOnlyStr == "true" {
			readOnly = true
		}
		if err := task.SetReadOnly(readOnly); err != nil {
			_, err = writer.Write([]byte("Unable to change task:" + err.Error()))
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
	// GetAllTasks(writer, req)
	GetTaskList(writer, req)
}

func PatchTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)

	/* TODO: Stoopid */
	changes := make(map[string]any)
	if req.Form.Has("status") {
		changes["status"] = tasks.ParseStatus(req.FormValue("status"))
	}
	if req.Form.Has("read-only") {
		var err error
		changes["read-only"], err = strconv.ParseBool(req.FormValue("read-only"))
		if err != nil {
			panic(err)
		}
	}

	patched := 0
	for _, id := range checkboxed {
		task, err := tasks.GetById(id)
		if err != nil {
			panic(err)
		}
		/* Status */
		if changes["status"] != nil && task.Status != changes["status"] {
			if err := task.SetStatus(changes["status"].(tasks.TaskStatus)); err != nil {
				panic(err)
			}
			patched++
		}
		/* Read only */
		if changes["read-only"] != nil && task.ReadOnly != changes["read-only"] {
			if err := task.SetReadOnly(changes["read-only"].(bool)); err != nil {
				panic(err)
			}
			patched++
		}
	}
	toasts.Info(writer, "Updated "+strconv.Itoa(patched)+" tasks", "Success")

	GetTaskList(writer, req)
	// GetAllTasks(writer, req)
}

func DeleteTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task[tasks.TaskFieldName]
	if task = idCheck(writer, req); task == nil {
		panic("Task not found")
	}
	taskId := task.Id
	tasks.Delete(taskId)
	toasts.Warning(writer, "Task "+strconv.Itoa(taskId)+" deleted", "Success")
	GetTaskList(writer, req)
}

func DeleteTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)
	deletedTasks := 0
	for _, id := range checkboxed {
		task, err := tasks.GetById(id)
		if err != nil {
			panic(err)
		}
		tasks.Delete(task.Id)
		deletedTasks++
	}
	toasts.Warning(writer, "Deleted "+strconv.Itoa(deletedTasks)+" tasks", "Success")
	GetTaskList(writer, req)
}
