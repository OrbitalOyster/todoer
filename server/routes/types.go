package routes

import (
	"net/http"
	"todoer/server/token"
	"todoer/tasks"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry

type TaskListData struct {
	Tasks      []tasks.Task
	TotalPages int
	Pagination []int
	Payload    token.Payload
	Checkboxes []bool
}

type MainPageData struct {
	Title      string
	PageSizes  []int
	TaskListData
}

type DatesOOBData struct {
	Payload token.Payload
}
