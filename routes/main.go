package routes

import (
	"net/http"
	"todoer/config"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

type mainPageData struct {
	Title      string
	PageSizes  []int
	TotalPages int
	Tasks      []tasks.Task
	Payload    jwt.Payload
}

func Main(writer http.ResponseWriter, req *http.Request) {
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	selectedTasks, _, totalPages := tasks.GetFromPayload(*payload)
	data := mainPageData{
		Title:      "todoer",
		PageSizes:  config.PageSizes,
		Tasks:      selectedTasks,
		TotalPages: totalPages,
		Payload:    jwt.Payload(*payload),
	}
	templates.Execute(writer, "main", data)
}
