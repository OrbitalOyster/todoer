package routes

import (
	"net/http"
	"todoer/tasks"
	"todoer/users"
	"todoer/utils"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry

type MainPageData struct {
	Title     string
	PageSizes []int
	Payload   utils.Payload
	TaskListData
}

type DatesOOBData struct {
	Payload utils.Payload
}

type EditTaskFormData struct {
	Task  *tasks.Task[tasks.TaskFieldName]
	Users []users.User
}
