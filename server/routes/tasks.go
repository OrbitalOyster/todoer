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
		Filter(tasks.Description, query.SearchBy).
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
	fieldStr := req.PathValue("field")
	field := tasks.ParseTaskFieldName(fieldStr)
	switch field {
	case tasks.Status:
		status := tasks.ParseStatus(req.FormValue("status"))
		tasks.FilterAndPatch(tasks.Id, id, tasks.Status, status)
	case tasks.ReadOnly:
		readOnly := false
		readOnlyStr := req.FormValue("read-only")
		if readOnlyStr == "true" {
			readOnly = true
		}
		tasks.FilterAndPatch(tasks.Id, id, tasks.ReadOnly, readOnly)
	default:
		toasts.Danger(writer, "Haxxor alert!", fmt.Sprintf("Invalid task field: %#v", fieldStr))
	}
	toasts.Success(writer, "Patch task", "Success")
	GetTaskList(writer, req)
}

func PatchTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)

	/* TODO: Stoopid */
	changes := make(map[tasks.TaskField]any)
	if req.Form.Has("status") {
		changes[tasks.Status] = tasks.ParseStatus(req.FormValue("status"))
	}
	if req.Form.Has("read-only") {
		var err error
		changes[tasks.ReadOnly], err = strconv.ParseBool(req.FormValue("read-only"))
		if err != nil {
			panic(err)
		}
	}

	patched := 0
	for _, id := range checkboxed {

		if changes[tasks.Status] != nil {
			tasks.FilterAndPatch(tasks.Id, id, tasks.Status, changes[tasks.Status])
			patched++
		}
		if changes[tasks.ReadOnly] != nil {
			tasks.FilterAndPatch(tasks.Id, id, tasks.ReadOnly, changes[tasks.ReadOnly])
			patched++
		}
		/*
		task, err := tasks.GetById(id)
		if err != nil {
			toasts.Warning(writer, "Unable to patch", fmt.Sprintf("Task #%d not found", id))
			continue
		}
		if changes["status"] != nil && task.Status != changes["status"] {
			if err := task.SetStatus(changes["status"].(tasks.TaskStatus)); err != nil {
				panic(err)
			}
			patched++
		}
		if changes["read-only"] != nil && task.ReadOnly != changes["read-only"] {
			if err := task.SetReadOnly(changes["read-only"].(bool)); err != nil {
				panic(err)
			}
			patched++
		}
		*/
	}

	toasts.Info(writer, "Updated "+strconv.Itoa(patched)+" tasks", "Success")
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
			toasts.Warning(writer, "Unable to delete", fmt.Sprintf("Task #%d not found", id))
			continue
		}
		tasks.Delete(task.Id)
		deletedTasks++
	}
	toasts.Warning(writer, "Deleted "+strconv.Itoa(deletedTasks)+" tasks", "Success")
	GetTaskList(writer, req)
}
