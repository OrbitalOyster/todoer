package routes

import (
	"net/http"
	"todoer/config"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
	"todoer/utils"
)

type mainPageData struct {
	Title      string
	PageSizes  []int
	TotalPages int
	Tasks      []tasks.Task
	Pagination []int
	Payload    jwt.Payload
}

func GetMainPage(writer http.ResponseWriter, req *http.Request) {
	payload := jwt.Get(req)
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	jwt.Update(payload, "Page", page, writer)
	data := mainPageData{
		Title:      "todoer",
		PageSizes:  config.PageSizes,
		TotalPages: totalPages,
		Tasks:      selectedTasks,
		Pagination: utils.GetPagination(totalPages, payload.Page),
		Payload:    jwt.Payload(*payload),
	}
	templates.ExecutePage(writer, "main", data)
}
