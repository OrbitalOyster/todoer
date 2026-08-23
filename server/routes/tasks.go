package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/toasts"
	"todoer/server/token"
	"todoer/tasks"
	"todoer/users"
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

type TasksPageData struct {
	Title      string
	Payload    utils.Payload
	Page       int
	PageSize   int
	PageSizes  []int
	TotalPages int
	SearchBy   string
	FromDate   string
	ToDate     string
	Pagination []int
}

const defaultPageSize = 10

type TasksQuery struct {
	Page     int
	Size     int
	SearchBy string
	FromDate string
	ToDate   string
}

func defaultTaskQuery() TasksQuery {
	fromDate, toDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	return TasksQuery{
		Page:     1,
		Size:     defaultPageSize,
		SearchBy: "",
		FromDate: fromDate.Format(utils.HTMLDateFormat),
		ToDate:   toDate.Format(utils.HTMLDateFormat),
	}
}

func (taskQuery *TasksQuery) Parse(rawQuery string) {
	parsed, err := url.ParseQuery(rawQuery)
	if err != nil {
		return
	}
	if parsed.Has("page") {
		if parsedPage, err := strconv.Atoi(parsed.Get("page")); err == nil {
			taskQuery.Page = parsedPage
		}
	}
	if parsed.Has("size") {
		if parsedSize, err := strconv.Atoi(parsed.Get("size")); err == nil {
			taskQuery.Size = parsedSize
		}
	}
	if parsed.Has("searchBy") {
		taskQuery.SearchBy = parsed.Get("searchBy")
	}

	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if parsed.Has("from") {
		if _, err := time.Parse(utils.HTMLDateFormat, parsed.Get("from")); err != nil {
			taskQuery.FromDate = defaultFromDate.Format(utils.HTMLDateFormat)
		} else {
			taskQuery.FromDate = parsed.Get("from")
		}
	}
	if parsed.Has("to") {
		if _, err := time.Parse(utils.HTMLDateFormat, parsed.Get("to")); err != nil {
			taskQuery.ToDate = defaultToDate.Format(utils.HTMLDateFormat)
		} else {
			taskQuery.ToDate = parsed.Get("to")
		}
	}

}

func (taskQuery TasksQuery) String() string {
	var result []string

	if taskQuery.Page != 1 {
		result = append(result, fmt.Sprintf("page=%d", taskQuery.Page))
	}

	if taskQuery.Size != defaultPageSize {
		result = append(result, fmt.Sprintf("size=%d", taskQuery.Size))
	}

	if taskQuery.SearchBy != "" {
		result = append(result, fmt.Sprintf("searchBy=%s", taskQuery.SearchBy))
	}

	defaultFromDate, defaultToDate := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	if taskQuery.FromDate != defaultFromDate.Format(utils.HTMLDateFormat) {
		result = append(result, fmt.Sprintf("from=%s", taskQuery.FromDate))
	}

	if taskQuery.ToDate != defaultToDate.Format(utils.HTMLDateFormat) {
		result = append(result, fmt.Sprintf("to=%s", taskQuery.ToDate))
	}

	if len(result) > 0 {
		return "?" + strings.Join(result, "&")
	} else {
		return ""
	}
}

func GetTasksPage(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()                                  /* Load actual payload */

	query := defaultTaskQuery()
	query.Parse(req.URL.RawQuery)

	pages.Execute(writer, "tasks", TasksPageData{
		Title:      "todoer - tasks",
		Payload:    payload,
		Page:       query.Page,
		PageSize:   query.Size,
		PageSizes:  config.PageSizes,
		TotalPages: 5,
		SearchBy:   query.SearchBy,
		FromDate:   query.FromDate,
		ToDate:     query.ToDate,
		Pagination: []int{1, 2, 3, 4, 5},
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

func HXGetTasks(writer http.ResponseWriter, req *http.Request) {
	result := ""

	query := defaultTaskQuery()

	/* Current browser query */
	url, err := url.Parse(req.Header.Get("HX-Current-URL"))
	if err != nil {
		panic(err)
	}
	query.Parse(url.RawQuery)

	/* New stuff */
	query.Parse(req.URL.RawQuery)

	writer.Header().Add("HX-Push-Url", "/tasks"+query.String())

	result += fmt.Sprintf(
		"%d tasks. Search by \"%s\". Page %d out of ?. From %s to %s",
		query.Size,
		query.SearchBy,
		query.Page,
		query.FromDate,
		query.ToDate,
	)

	/* Update calendar elements if both dates are set */
	if req.Form.Has("from") && req.Form.Has("to") {
		pages.ExecutePartial(
			writer,
			"task-dates-oob-new",
			struct {
				FromDate string
				ToDate   string
			}{
				FromDate: query.FromDate,
				ToDate:   query.ToDate,
			},
		)
	}

	//tasks.Get()

	/*
		pages.ExecutePartial(writer, "task-list", TaskListData{
			Tasks:      selectedTasks,
			Checkboxes: checkboxes,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(totalPages, page),
			Payload:    payload,
		})
	*/
	writer.Write([]byte("<strong>" + result + "</strong>"))
}

func GetAllTasks(writer http.ResponseWriter, req *http.Request) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()                                  /* Load actual payload */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	checkboxedTasks := getCheckboxedTasks(req)
	checkboxes := make([]bool, len(selectedTasks))
	for i, selectedTask := range selectedTasks {
		if slices.Contains(checkboxedTasks, selectedTask.Id) {
			checkboxes[i] = true
		} else {
			checkboxes[i] = false
		}
	}
	pages.ExecutePartial(writer, "task-list", TaskListData{
		Tasks:      selectedTasks,
		Checkboxes: checkboxes,
		TotalPages: totalPages,
		Pagination: utils.GetPagination(totalPages, page),
		Payload:    payload,
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
	GetAllTasks(writer, req)
}

func PutTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task
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
	GetAllTasks(writer, req)
}

func PatchTask(writer http.ResponseWriter, req *http.Request) {
	var task *tasks.Task
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
	GetAllTasks(writer, req)
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
	GetAllTasks(writer, req)
}
