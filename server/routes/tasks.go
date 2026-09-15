package routes

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
	"todoer/collection"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/tasks"
	"todoer/users"
	"todoer/utils"
)

const defaultPageSize = 10

type TaskListData struct {
	Tasks      []collection.Item[tasks.TaskField]
	Page       uint
	PageSize   uint
	TotalPages uint
	Pagination []uint
	SearchBy   string
	FromDate   time.Time
	ToDate     time.Time
	SortBy     string
	SortDesc   bool
	Checkboxes []bool
}

func getTasks(query TasksQuery[tasks.TaskField]) (collection.Collection[tasks.TaskField], uint, uint) {
	result := tasks.GetAll()
	result = result.
		MoreThan(
			tasks.Datetime,
			query.FromDate,
		).
		/* "Not after 20/03/2026" means "Not after 20/03/2026 23:59:59"  */
		LessThan(
			tasks.Datetime,
			query.ToDate.Add(time.Hour*24-time.Second),
		).
		Filter(tasks.Description, []any{query.SearchBy}).
		SortBy(query.SortBy)
	if query.SortDesc {
		result.Reverse()
	}
	return result.GetPage(query.Page, query.Size)
}

func GetTasksPage(writer http.ResponseWriter, req *http.Request) {
	payload := utils.GetTokenPayload(req)
	query, urlUpdate := CreateQueryFromRequest(req)
	/* Update URL */
	if urlUpdate {
		queryStr := query.String()
		if len(queryStr) > 0 {
			queryStr = "?" + queryStr
		}
		http.Redirect(writer, req, "/tasks"+queryStr, http.StatusSeeOther)
		return
	}

	selectedTasks, page, totalPages := getTasks(query)
	checkboxedTasks := getCheckboxedTasks(req)
	checkboxes := make([]bool, selectedTasks.Length())
	for i, selectedTask := range selectedTasks.Items {
		checkboxes[i] = slices.Contains(checkboxedTasks, selectedTask.Field(tasks.Id).(int))
	}

	pages.Execute(writer, "tasks", struct {
		Title     string
		Payload   utils.Payload
		PageSizes []int
		TaskListData
	}{
		Title:      "todoer - tasks",
		Payload:    payload,
		PageSizes:  config.PageSizes,
		Tasks:      selectedTasks.Items,
		Page:       page,
		PageSize:   query.Size,
		TotalPages: totalPages,
		Pagination: utils.GetPagination(totalPages, page),
		SearchBy:   query.SearchBy,
		FromDate:   query.FromDate,
		ToDate:     query.ToDate,
		SortBy:     query.SortBy.String(),
		SortDesc:   query.SortDesc,
		Checkboxes: checkboxes,
	})
}

func getCheckboxedTasks(req *http.Request) (result []int) {
	if !req.Form.Has("checked") { /* Nothing checked */
		return
	}
	for _, checkboxStr := range req.Form["checked"] {
		n, err := strconv.Atoi(checkboxStr)
		if err != nil {
			panic(err) /* Unparseable string */
		}
		result = append(result, n)
	}
	return
}

func GetTaskList(writer http.ResponseWriter, req *http.Request) {
	query, urlUpdate := CreateQueryFromRequest(req)

	/* Get tasks */
	tasksOnCurrentPage, page, numberOfPages := getTasks(query)
	if page != query.Page {
		query.Page = page
	}

	checkboxedTasks := getCheckboxedTasks(req)
	checkboxes := make([]bool, tasksOnCurrentPage.Length())
	for i, selectedTask := range tasksOnCurrentPage.Items {
		checkboxes[i] = slices.Contains(checkboxedTasks, selectedTask.Field(tasks.Id).(int))
	}

	/* Update calendar elements if both dates are set */
	if req.Form.Has("from") && req.Form.Has("to") {
		pages.ExecutePartial(
			writer,
			"task-dates-oob",
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
	if urlUpdate {
		queryStr := query.String()
		if len(queryStr) > 0 {
			queryStr = "?" + queryStr
		}
		writer.Header().Add("HX-Push-Url", "/tasks"+queryStr)
	}

	/* Send actual list */
	pages.ExecutePartial(writer, "task-list", TaskListData{
		Tasks:      tasksOnCurrentPage.Items,
		Page:       page,
		PageSize:   query.Size,
		TotalPages: numberOfPages,
		Pagination: utils.GetPagination(numberOfPages, page),
		SortBy:     query.SortBy.String(),
		SortDesc:   query.SortDesc,
		Checkboxes: checkboxes,
	})
}

func GetAddTaskForm(writer http.ResponseWriter, req *http.Request) {
	pages.ExecutePartial(writer, "addTaskForm", nil)
}

func GetEditTaskForm(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		pages.ExecutePartial(writer, "taskNotFound", nil)
		return
	}
	data := struct {
		Task  *tasks.Task[tasks.TaskField]
		Users []users.User
	}{
		task,
		users.GetAllUsers(),
	}
	pages.ExecutePartial(writer, "editTaskForm", data)
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		toasts.Warning(writer, "Unable to clone", fmt.Sprintf("Task #%s not found", id))
		return
	}
	pages.ExecutePartial(writer, "cloneTaskForm", task)
}

func AddTask(writer http.ResponseWriter, req *http.Request) {
	payload := utils.GetTokenPayload(req)
	user := payload.UserID
	description := req.FormValue("description")
	tasks.Add(user, description)
	writer.Header().Set("HX-Trigger", "hideModal")
	toasts.Success(writer, "New task", "Success")
	GetTaskList(writer, req)
}

func PutTask(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		toasts.Warning(writer, "Unable to edit", fmt.Sprintf("Task #%s not found", id))
		writer.Header().Set("HX-Trigger", "hideModal")
		GetTaskList(writer, req)
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
	GetTaskList(writer, req)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	/* Check if id is parseable */
	idStr := req.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		toasts.Danger(
			writer,
			"Haxxor alert!",
			fmt.Sprintf("Not a valid task id: %#v", id),
		)
		return
	}
	/* Check if field is parseable */
	fieldStr := req.PathValue("field")
	field, err := tasks.ParseTaskField(fieldStr)
	if err != nil {
		toasts.Danger(
			writer,
			"Haxxor alert!",
			fmt.Sprintf("Not a valid field: %#v", fieldStr),
		)
		return
	}

	var patched uint = 0
	switch field {
	case tasks.Status:
		if status, err := tasks.ParseStatus(req.FormValue("status")); err != nil {
			toasts.Danger(
				writer,
				"Haxxor alert!",
				fmt.Sprintf("Invalid status: %#v", status),
			)
		} else {
			patched += tasks.FilterAndPatch(tasks.Id, []any{id}, tasks.Status, status)
		}
	case tasks.ReadOnly:
		readOnly := false
		readOnlyStr := req.FormValue("read-only")
		if readOnlyStr == "true" {
			readOnly = true
		}
		tasks.FilterAndPatch(tasks.Id, []any{id}, tasks.ReadOnly, readOnly)
	default:
		toasts.Danger(writer, "Haxxor alert!", fmt.Sprintf("Invalid task field: %#v", fieldStr))
	}
	toasts.Success(writer, "Patch task", "Success")
	GetTaskList(writer, req)
}

func PatchTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)
	/* Nothing selected */
	if len(checkboxed) == 0 {
		toasts.Warning(writer, "No tasks updated", "Nothing selected")
	} else {
		filterBy := make([]any, len(checkboxed))
		for i, c := range checkboxed {
			filterBy[i] = c
		}
		message := ""
		/* Status */
		if req.Form.Has("status") {
			statusStr := req.FormValue("status")
			status, err := tasks.ParseStatus(statusStr)
			/* Invalid status string */
			if err != nil {
				toasts.Danger(
					writer,
					"Error",
					fmt.Sprintf("Not a valid status: %#v", statusStr),
				)
			} else {
				patchedStatus := tasks.FilterAndPatch(
					tasks.Id,
					filterBy,
					tasks.Status,
					status,
				)
				if patchedStatus > 0 {
					message += fmt.Sprintf("Updated status: %d", patchedStatus)
				}
			}
		}
		/* Read only */
		if req.Form.Has("read-only") {
			readOnlyStr := req.FormValue("read-only")
			readOnly, err := strconv.ParseBool(readOnlyStr)
			if err != nil {
				toasts.Danger(
					writer,
					"Error",
					fmt.Sprintf("Not a valid boolean: %#v", readOnlyStr),
				)
			}
			patchedReadOnly := tasks.FilterAndPatch(
				tasks.Id,
				filterBy,
				tasks.ReadOnly,
				readOnly,
			)
			if patchedReadOnly > 0 {
				message += fmt.Sprintf("Updated lock: %d", patchedReadOnly)
			}
		}
		/* Send results to user */
		if message != "" {
			toasts.Info(writer, "Updated tasks", message)
		} else {
			toasts.Warning(writer, "Updated tasks", "Nothing changed")
		}
	}
	/* Done */
	GetTaskList(writer, req)
}

func DeleteTask(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		toasts.Warning(writer, "Unable to delete", fmt.Sprintf("Task #%s not found", id))
		GetTaskList(writer, req)
		return
	}
	taskId := task.Id
	tasks.DeleteOne(taskId)
	toasts.Warning(writer, "Task "+strconv.Itoa(taskId)+" deleted", "Success")
	GetTaskList(writer, req)
}

func DeleteTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)
	deletedTasks := 0
	for _, id := range checkboxed {
		task, err := tasks.GetById(id)
		if err != nil {
			toasts.Warning(writer, "Unable to delete", fmt.Sprintf("Task #%d not found", id))
			continue
		}
		tasks.DeleteOne(task.Id)
		deletedTasks++
	}
	toasts.Warning(writer, "Deleted "+strconv.Itoa(deletedTasks)+" tasks", "Success")
	GetTaskList(writer, req)
}
