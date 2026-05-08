package routes

import (
	"net/http"
	"slices"
	"strconv"
	"time"
	"todoer/config"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/utils"
)

func executeTemplate(writer http.ResponseWriter, payload *jwt.Payload, selectedTasks []tasks.Task, totalPages int, page int) {
	templates.ExecutePartial(
		writer,
		"task-list",
		TaskListData{
			Tasks:      selectedTasks,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(totalPages, page),
			Payload:    jwt.Payload(*payload),
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
	payload := jwt.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, size,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	jwt.Update(payload, "PageSize", size, writer)
	jwt.Update(payload, "Page", page, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func setPage(page int, payload *jwt.Payload, writer http.ResponseWriter) {
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	jwt.Update(payload, "Page", page, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetPage(writer http.ResponseWriter, req *http.Request) {
	pageStr := req.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	payload := jwt.Get(req)
	setPage(page, payload, writer)
}

func NextPage(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	page := payload.Page + 1
	/* Return updated task table */
	setPage(page, payload, writer)
}

func PreviousPage(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	page := payload.Page - 1
	/* Return updated task table */
	setPage(page, payload, writer)
}

func SetSortBy(writer http.ResponseWriter, req *http.Request) {
	columnStr := req.PathValue("column")
	columnInt, err := strconv.Atoi(columnStr)
	if err != nil {
		columnInt = 0
	}
	column := utils.SortableColumn(columnInt)
	payload := jwt.Get(req)
	sortAsc := payload.SortAsc
	/* Reverse sort */
	if payload.SortBy == column {
		sortAsc = !sortAsc
	}
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		column, sortAsc)
	/* Update token */
	jwt.Update(payload, "Page", page, writer)
	jwt.Update(payload, "SortBy", column, writer)
	jwt.Update(payload, "SortAsc", sortAsc, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	searchBy := req.FormValue("searchBy")
	payload := jwt.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		searchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	jwt.Update(payload, "Page", page, writer)
	jwt.Update(payload, "SearchBy", searchBy, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetFromDate(writer http.ResponseWriter, req *http.Request) {
	fromDateStr := req.FormValue("fromDate")
	_, err := time.Parse("2006-01-02", fromDateStr)
	/* User sent stoopid */
	if err != nil {
		fromDate, _ := utils.GetMonthBounds()
		fromDateStr = fromDate.Format("2006-01-02")
	}
	payload := jwt.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		fromDateStr, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	jwt.Update(payload, "Page", page, writer)
	jwt.Update(payload, "FromDate", fromDateStr, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}

func SetToDate(writer http.ResponseWriter, req *http.Request) {
	toDateStr := req.FormValue("toDate")
	_, err := time.Parse("2006-01-02", toDateStr)
	/* User sent stoopid */
	if err != nil {
		toDate, _ := utils.GetMonthBounds()
		toDateStr = toDate.Format("2006-01-02")
	}
	payload := jwt.Get(req)
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, toDateStr,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	jwt.Update(payload, "Page", page, writer)
	jwt.Update(payload, "ToDate", toDateStr, writer)
	/* Done */
	executeTemplate(writer, payload, selectedTasks, totalPages, page)
}
