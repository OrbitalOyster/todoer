package routes

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"todoer/config"
	"todoer/cookies"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/utils"
)

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
	selectedTasks, page, totalPages := tasks.GetFromPayload(*payload)
	data := struct {
		Tasks      []tasks.Task
		Page       int
		TotalPages int
		SortBy     utils.SortableColumn
		SortAsc    bool
	}{
		Tasks:      selectedTasks,
		Page:       page,
		TotalPages: totalPages,
		SortBy:     payload.SortBy,
		SortAsc:    payload.SortAsc,
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
	selectedTasks, page, totalPages := tasks.GetFromPayload(*payload)
	data := struct {
		Tasks      []tasks.Task
		Page       int
		TotalPages int
		SortBy     utils.SortableColumn
		SortAsc    bool
	}{
		Tasks:      selectedTasks,
		Page:       page,
		TotalPages: totalPages,
		SortBy:     payload.SortBy,
		SortAsc:    payload.SortAsc,
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
	selectedTasks, column, totalPages := tasks.GetFromPayload(*payload)
	data := struct {
		Tasks      []tasks.Task
		Page       int
		TotalPages int
		SortBy     utils.SortableColumn
		SortAsc    bool
	}{
		Tasks:      selectedTasks,
		Page:       column,
		TotalPages: totalPages,
		SortBy:     payload.SortBy,
		SortAsc:    payload.SortAsc,
	}
	templates.ExecutePartial(writer, "task-list", data)
}

func SetSearchBy(writer http.ResponseWriter, req *http.Request) {
	s := req.FormValue("searchBy")
	log.Println("Setting searchBy to", s)
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
	selectedTasks, column, totalPages := tasks.GetFromPayload(*payload)
	data := struct {
		Tasks      []tasks.Task
		Page       int
		TotalPages int
		SortBy     utils.SortableColumn
		SortAsc    bool
	}{
		Tasks:      selectedTasks,
		Page:       column,
		TotalPages: totalPages,
		SortBy:     payload.SortBy,
		SortAsc:    payload.SortAsc,
	}
	templates.ExecutePartial(writer, "task-list", data)
}
