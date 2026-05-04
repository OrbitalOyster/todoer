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

func Main(writer http.ResponseWriter, req *http.Request) {
	payload, err := jwt.Get(req)
	/* Major screw up */
	if err != nil {
		panic(err)
	}
	selectedTasks, totalPages, page := tasks.GetFromPayload(*payload)
	jwt.HealthCheck(payload, page, writer)
	data := mainPageData{
		Title:      "todoer",
		PageSizes:  config.PageSizes,
		TotalPages: totalPages,
		Tasks:      selectedTasks,
		Pagination: utils.GetPagination(totalPages, payload.Page),
		Payload:    jwt.Payload(*payload),
	}
	templates.Execute(writer, "main", data)
}
