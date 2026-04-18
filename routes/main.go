package routes

import (
	"net/http"
	"todoer/config"
	"todoer/jwt"
	"todoer/tasks"
	"todoer/templates"
)

func Main(writer http.ResponseWriter, req *http.Request) {
	payload, err := jwt.Get(req)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	selectedTasks, page, totalPages := tasks.GetFromPayload(*payload)
	data := struct {
		Title            string
		Username         string
		PageSizes        []int
		SelectedPageSize int
		Tasks            []tasks.Task
		Page             int
		TotalPages       int
	}{
		Title:            "todoer",
		Username:         payload.UserID,
		PageSizes:        config.PageSizes,
		SelectedPageSize: payload.PageSize,
		Tasks:            selectedTasks,
		Page:             page,
		TotalPages:       totalPages,
	}
	templates.Execute(writer, "main", data)
}
