package routes

import (
	"net/http"
	"todoer/config"
	"todoer/server/pages"
	"todoer/server/token"
	"todoer/tasks"
	"todoer/utils"
)

func GetMainPage(writer http.ResponseWriter, req *http.Request) {
	payload := token.Get(req)
	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	token.Update(payload, "Page", page, writer)
	pages.Execute(writer, "main", MainPageData{
		Title:      "todoer",
		PageSizes:  config.PageSizes,
		TotalPages: totalPages,
		Tasks:      selectedTasks,
		Pagination: utils.GetPagination(totalPages, page),
		Payload:    token.Payload(*payload),
	})
}
