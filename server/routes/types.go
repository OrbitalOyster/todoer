package routes

import (
	"net/http"
	"todoer/tasks"
	"todoer/users"
	"todoer/utils"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry

type TaskListData struct {
	Tasks      []tasks.Task
	TotalPages int
	Pagination []int
	Payload    utils.Payload
	Checkboxes []bool
}

type MainPageData struct {
	Title     string
	PageSizes []int
	TaskListData
}

type DatesOOBData struct {
	Payload utils.Payload
}

type EditTaskFormData struct {
	Task  *tasks.Task
	Users []users.User
}
