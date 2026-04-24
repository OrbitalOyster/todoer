package routes

import (
	"net/http"
	"slices"
	"strconv"
	"time"
	"todoer/config"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/utils"
)

type taskListData struct {
	Tasks      []tasks.Task
	TotalPages int
	Payload jwt.Payload
}

func SetTaskTablePageSize(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(req.FormValue("size"))
	if err != nil || !slices.Contains(config.PageSizes, size) {
		size = config.DefaultPageSize
	}
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.PageSize = size
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetPage(writer http.ResponseWriter, req *http.Request) {
	pageStr := req.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.Page = page
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData {
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetSortBy(writer http.ResponseWriter, req *http.Request) {
	columnStr := req.PathValue("column")
	column, err := strconv.Atoi(columnStr)
	if err != nil {
		column = 0
	}
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}

	if payload.SortBy == utils.SortableColumn(column) {
		payload.SortAsc = !payload.SortAsc
	}

	payload.SortBy = utils.SortableColumn(column)
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	s := req.FormValue("searchBy")
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.SearchBy = s
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetFromDate(writer http.ResponseWriter, req *http.Request) {
	fromDateStr := req.FormValue("fromDate")
	_, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		panic(err)
	}
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.FromDate = fromDateStr
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetToDate(writer http.ResponseWriter, req *http.Request) {
	toDateStr := req.FormValue("toDate")
	_, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		panic(err)
	}
	/* Update token/cookies */
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.ToDate = toDateStr
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)
	/* Return updated task table */
	selectedTasks, totalPages := tasks.GetFromPayload(*payload)
	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload: jwt.Payload(*payload),
	}
	templates.ExecutePartial(writer, "task-list", data)
}
