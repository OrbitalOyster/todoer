package routes

import (
	"net/http"
	"slices"
	"strconv"
	"time"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/token"
	"todoer/tasks"
	"todoer/utils"
)

func executeTemplate(writer http.ResponseWriter, payload *token.Payload, selectedTasks []tasks.Task, totalPages int, page int) {
	pages.ExecutePartial(
		writer,
		"task-list",
		TaskListData{
			Tasks:      selectedTasks,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(totalPages, page),
			Payload:    token.Payload(*payload),
			Checkboxes: make([]bool, payload.PageSize),
		})
}

func SetPageSize(writer http.ResponseWriter, req *http.Request) {
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
	payload := token.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, size,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	token.Update(payload, "PageSize", size, writer)
	token.Update(payload, "Page", page, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func setPage(page int, payload *token.Payload, writer http.ResponseWriter) {
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	token.Update(payload, "Page", page, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetPage(writer http.ResponseWriter, req *http.Request) {
	pageStr := req.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	payload := token.Get(req)
	setPage(page, payload, writer)
}

func NextPage(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	page := payload.Page + 1
	/* Return updated task table */
	setPage(page, payload, writer)
}

func PreviousPage(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	page := payload.Page - 1
	/* Return updated task table */
	setPage(page, payload, writer)
}

func SetSortBy(writer http.ResponseWriter, req *http.Request) {
	fieldStr := req.PathValue("field")
	fieldInt, err := strconv.Atoi(fieldStr)
	if err != nil {
		fieldInt = 0
	}
	field := utils.SortableField(fieldInt)
	payload := token.Get(req)
	sortAsc := payload.SortAsc
	/* Reverse sort */
	if payload.SortBy == field {
		sortAsc = !sortAsc
	}
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		field, sortAsc)
	/* Update token */
	token.Update(payload, "Page", page, writer)
	token.Update(payload, "SortBy", field, writer)
	token.Update(payload, "SortAsc", sortAsc, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	searchBy := req.FormValue("searchBy")
	payload := token.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		searchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	token.Update(payload, "Page", page, writer)
	token.Update(payload, "SearchBy", searchBy, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetDate(writer http.ResponseWriter, req *http.Request) {
	fromDateStr := req.FormValue("from-date")
	toDateStr := req.FormValue("to-date")
	payload := token.Get(req)
	fromDateFallback, toDateFallback := utils.GetMonthBounds(time.Now().Year(), time.Now().Month())
	/* Setting from date? */
	if fromDateStr != "" {
		_, err := time.Parse(utils.HTMLDateFormat, fromDateStr)
		/* User sent stoopid */
		if err != nil {
			fromDateStr = fromDateFallback.Format(utils.HTMLDateFormat)
		}
	} else {
		fromDateStr = payload.FromDate
	}
	/* Setting to date? */
	if toDateStr != "" {
		_, err := time.Parse(utils.HTMLDateFormat, toDateStr)
		/* User sent stoopid */
		if err != nil {
			toDateStr = toDateFallback.Format(utils.HTMLDateFormat)
		}
	} else {
		toDateStr = payload.ToDate
	}
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		fromDateStr, toDateStr,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	token.Update(payload, "Page", page, writer)
	token.Update(payload, "FromDate", fromDateStr, writer)
	token.Update(payload, "ToDate", toDateStr, writer)
	/* Update calendar elements if both dates are set */
	if req.Form.Has("from-date") && req.Form.Has("to-date") {
		pages.ExecutePartial(
			writer,
			"task-dates-oob",
			DatesOOBData{Payload: *payload})
	}
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}
