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
	Pagination []int
	Payload    jwt.Payload
}

func executeResult(writer http.ResponseWriter, payload *jwt.Payload) {
	/* Get data */
	selectedTasks, totalPages, page := tasks.GetFromPayload(*payload)
	jwt.HealthCheck(payload, page, writer)
	/* Set cookies */
	token := jwt.Create(*payload)
	cookies.Set(writer, token, payload.RememberMe)

	data := taskListData{
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Pagination: utils.GetPagination(totalPages, payload.Page),
		Payload:    jwt.Payload(*payload),
	}
	/* Render template */
	templates.ExecutePartial(writer, "task-list", data)
}

func SetTaskTablePageSize(writer http.ResponseWriter, req *http.Request) {
	/* Check if form is ok */
	if err := req.ParseForm(); err != nil {
		http.Error(writer, "Haxxor alert!", http.StatusBadRequest)
		return
	}
	size, err := strconv.Atoi(req.FormValue("size"))
	/* Wrong page size, somehow */
	if err != nil || !slices.Contains(config.PageSizes, size) {
		size = config.DefaultPageSize
	}
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.PageSize = size
	/* Return updated task table */
	executeResult(writer, payload)
}

func SetPage(writer http.ResponseWriter, req *http.Request) {
	pageStr := req.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}
	payload, err := jwt.Get(req)
	/* Major screw up */
	if err != nil {
		panic(err)
	}
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.GetSome(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	err = jwt.Update(payload, "Page", page, writer)
	/* Major screw up */
	if err != nil {
		panic(err)
	}
	/* Return updated task table */
	templates.ExecutePartial(
		writer,
		"task-list",
		taskListData{
			Tasks:      selectedTasks,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(totalPages, payload.Page),
			Payload:    jwt.Payload(*payload),
		})
}

func NextPage(writer http.ResponseWriter, req *http.Request) {
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	page := payload.Page + 1
	payload.Page = page
	/* Return updated task table */
	executeResult(writer, payload)
}

func PreviousPage(writer http.ResponseWriter, req *http.Request) {
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	page := payload.Page - 1
	payload.Page = page
	/* Return updated task table */
	executeResult(writer, payload)
}

func SetSortBy(writer http.ResponseWriter, req *http.Request) {
	columnStr := req.PathValue("column")
	column, err := strconv.Atoi(columnStr)
	if err != nil {
		column = 0
	}
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	/* Reverse sort */
	if payload.SortBy == utils.SortableColumn(column) {
		payload.SortAsc = !payload.SortAsc
	}
	payload.SortBy = utils.SortableColumn(column)
	/* Return updated task table */
	executeResult(writer, payload)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	s := req.FormValue("searchBy")
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.SearchBy = s
	/* Return updated task table */
	executeResult(writer, payload)
}

func SetFromDate(writer http.ResponseWriter, req *http.Request) {
	fromDateStr := req.FormValue("fromDate")
	_, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		panic(err)
	}
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.FromDate = fromDateStr
	/* Return updated task table */
	executeResult(writer, payload)
}

func SetToDate(writer http.ResponseWriter, req *http.Request) {
	toDateStr := req.FormValue("toDate")
	_, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		panic(err)
	}
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	payload.ToDate = toDateStr
	/* Return updated task table */
	executeResult(writer, payload)
}
