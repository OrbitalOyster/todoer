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
	data := struct {
		Title            string
		Username         string
		Tasks            []tasks.Task
		PageSizes        []int
		SelectedPageSize int
	}{
		Title:            "todoer",
		Username:         payload.UserID,
		Tasks:            tasks.GetFromPayload(*payload),
		PageSizes:        config.PageSizes,
		SelectedPageSize: payload.PageSize,
	}
	templates.Execute(writer, "main", data)
}
