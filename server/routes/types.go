package routes

import (
	"net/http"
	"todoer/server/token"
	"todoer/tasks"
)

type RouterEntry func(http.ResponseWriter, *http.Request)
type RouterMap map[string]RouterEntry

type MainPageData struct {
	Title      string
	PageSizes  []int
	TotalPages int
	Tasks      []tasks.Task
	Pagination []int
	Payload    token.Payload
}

type TaskListData struct {
	Tasks      []tasks.Task
	TotalPages int
	Pagination []int
	Payload    token.Payload
}
