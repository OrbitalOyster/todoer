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

func executeTemplate(
	writer http.ResponseWriter,
	req *http.Request,
	selectedTasks []tasks.Task,
	totalPages int,
	page int) {
	payload := req.Context(). /* Get context from request */
					Value("token").(*token.Token[utils.Payload]). /* Get "token" field */
					GetPayload()                                  /* Load actual payload */
	pages.ExecutePartial(
		writer,
		"task-list",
		TaskListData{
			Tasks:      selectedTasks,
			TotalPages: totalPages,
			Pagination: utils.GetPagination(totalPages, page),
			Payload:    payload,
			Checkboxes: make([]bool, payload.PageSize),
		})
}

func SetPageSize(writer http.ResponseWriter, req *http.Request) {
	size, err := strconv.Atoi(req.FormValue("size"))
	/* Wrong page size, somehow */
	if err != nil || !slices.Contains(config.PageSizes, size) {
		size = config.DefaultPageSize
	}
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, size,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	payload.PageSize = size
	payload.Page = page
	token.SetPayload(payload)
	/* Done */
	executeTemplate(writer, req, selectedTasks, totalPages, page)
}

func setPage(page int, req *http.Request, writer http.ResponseWriter) {
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	payload.Page = page
	token.SetPayload(payload)
	/* Done */
	executeTemplate(writer, req, selectedTasks, totalPages, page)
}

func SetPage(writer http.ResponseWriter, req *http.Request) {
	pageStr := req.PathValue("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	setPage(page, req, writer)
}

func NextPage(writer http.ResponseWriter, req *http.Request) {
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	page := payload.Page + 1
	/* Return updated task table */
	setPage(page, req, writer)
}

func PreviousPage(writer http.ResponseWriter, req *http.Request) {
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	page := payload.Page - 1
	/* Return updated task table */
	setPage(page, req, writer)
}

func SetSortBy(writer http.ResponseWriter, req *http.Request) {
	fieldStr := req.PathValue("field")
	field := utils.ParseSortableField(fieldStr)
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
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
	payload.Page = page
	payload.SortBy = field
	payload.SortAsc = sortAsc
	token.SetPayload(payload)
	/* Done */
	executeTemplate(writer, req, selectedTasks, totalPages, page)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	searchBy := req.FormValue("searchBy")
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
	/* Get tasks */
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		searchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc)
	/* Update token */
	payload.Page = page
	payload.SearchBy = searchBy
	token.SetPayload(payload)
	/* Done */
	executeTemplate(writer, req, selectedTasks, totalPages, page)
}

func SetDate(writer http.ResponseWriter, req *http.Request) {
	fromDateStr := req.FormValue("from-date")
	toDateStr := req.FormValue("to-date")
	token := req.Context(). /* Get context from request */
				Value("token").(*token.Token[utils.Payload]) /* Get "token" field */
	payload := token.GetPayload() /* Load actual payload */
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
	payload.Page = page
	payload.FromDate = fromDateStr
	payload.ToDate = toDateStr
	token.SetPayload(payload)
	/* Update calendar elements if both dates are set */
	if req.Form.Has("from-date") && req.Form.Has("to-date") {
		pages.ExecutePartial(
			writer,
			"task-dates-oob",
			DatesOOBData{Payload: payload})
	}
	/* Done */
	executeTemplate(writer, req, selectedTasks, totalPages, page)
}
