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
	result := tasks.All.
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
		checkboxes[i] = slices.Contains(
			checkboxedTasks,
			selectedTask.Field(tasks.Id).(int), /* TODO: Type assertion */
		)
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
		checkboxes[i] = slices.Contains(
			checkboxedTasks,
			selectedTask.Field(tasks.Id).(int), /* TODO: Type assertion */
		)
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
	} else {
		data := struct {
			Task  tasks.Task[tasks.TaskField]
			Users []users.User
		}{
			task,
			users.GetAllUsers(),
		}
		pages.ExecutePartial(writer, "editTaskForm", data)
	}
}

func GetCloneTaskForm(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		toasts.Danger(writer, "Error", err.Error())
	} else {
		pages.ExecutePartial(writer, "cloneTaskForm", task)
	}
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
		toasts.Danger(writer, "Error", err.Error())
	} else {
		description, user, status, readOnlyStr :=
			req.FormValue("description"),
			req.FormValue("user"),
			req.FormValue("status"),
			req.FormValue("readOnly")
		readOnly := false
		if readOnlyStr == "true" {
			readOnly = true
		}
		/* TODO: Lunacy */
		updatedDescription, _ := tasks.Patch([]int{task.Id}, tasks.Description, description)
		updatedUser, _ := tasks.Patch([]int{task.Id}, tasks.User, user)
		updatedStatus, _ := tasks.Patch([]int{task.Id}, tasks.Status, status)
		updatedReadOnly, _ :=  tasks.Patch([]int{task.Id}, tasks.ReadOnly, readOnly)
		if updatedDescription > 0 || updatedUser > 0 || updatedStatus > 0 || updatedReadOnly > 0 {
			toasts.Success(writer, "Update task", "Success")
		} else {
			toasts.Info(writer, "Update task", "Nothing changed")
		}
	}
	/* Done */
	writer.Header().Set("HX-Trigger", "hideModal")
	GetTaskList(writer, req)
}

func PatchTasks(writer http.ResponseWriter, req *http.Request) {
	checkboxed := getCheckboxedTasks(req)
	/* Nothing selected */
	if len(checkboxed) == 0 {
		toasts.Warning(writer, "No tasks updated", "Nothing selected")
	} else {
		var errors []error
		switch {
		case req.Form.Has("status"):
			patched, err := tasks.Patch(checkboxed, tasks.Status, req.FormValue("status"))
			errors = append(errors, err...)
			/* Report results */
			toasts.Success(writer, "Updated status", fmt.Sprintf("Tasks updated: %d", patched))
		case req.Form.Has("readOnly"):
			readOnly := false
			if req.FormValue("readOnly") == "true" {
				readOnly = true
			}
			patched, err := tasks.Patch(checkboxed, tasks.ReadOnly, readOnly)
			errors = append(errors, err...)
			/* Report results */
			toasts.Success(writer, "Updated lock", fmt.Sprintf("Tasks updated: %d", patched))
		}
		/* Report errors */
		for _, e := range errors {
			toasts.Danger(writer, "Error", e.Error())
		}
	}
	/* Done */
	GetTaskList(writer, req)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	task, err := tasks.GetById(id)
	if err != nil {
		toasts.Danger(writer, "Error", err.Error())
	} else {
		var (
			errors   []error
			noChange = true
		)
		switch {
		case req.Form.Has("status"):
			patched, err := tasks.Patch([]int{task.Id}, tasks.Status, req.FormValue("status"))
			if patched > 0 {
				noChange = false
				toasts.Success(writer, "Updated status", "Success")
			}
			errors = append(errors, err...)
		case req.Form.Has("readOnly"):
			readOnly := false
			if req.FormValue("readOnly") == "true" {
				readOnly = true
			}
			patched, err := tasks.Patch([]int{task.Id}, tasks.ReadOnly, readOnly)
			if patched > 0 {
				noChange = false
				toasts.Success(writer, "Updated lock", "Success")
			}
			errors = append(errors, err...)
		}
		if noChange {
			toasts.Info(writer, "Update task", "Nothing changed")
		}
		/* Report errors */
		for _, e := range errors {
			toasts.Danger(writer, "Error", e.Error())
		}
	}
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
