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
	// payload := token.Get(req)

	funkyToken := token.CreateFunkyToken[token.Payload](
		req,
		&writer,
		config.CookieName,
		config.JWTSecret,
		config.CookieShortLifetime,
	)
	/* Should not happen */
	if err := funkyToken.Load(); err != nil {
		panic(err)
	}
	payload := funkyToken.GetPayload()

	selectedTasks, totalPages, page := tasks.Get(
		payload.FromDate, payload.ToDate,
		payload.SearchBy,
		payload.Page, payload.PageSize,
		payload.SortBy, payload.SortAsc,
	)
	pages.Execute(writer, "main", MainPageData{
		Title:     "todoer",
		PageSizes: config.PageSizes,
		TaskListData: TaskListData{
			TotalPages: totalPages,
			Tasks:      selectedTasks,
			Pagination: utils.GetPagination(totalPages, page),
			Payload:    payload,
			Checkboxes: make([]bool, payload.PageSize),
		},
	})
}
