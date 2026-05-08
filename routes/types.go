package routes

import (
	"todoer/tasks"
	"todoer/jwt"
)

type MainPageData struct {
	Title      string
	PageSizes  []int
	TotalPages int
	Tasks      []tasks.Task
	Pagination []int
	Payload    jwt.Payload
}

type TaskListData struct {
	Tasks      []tasks.Task
	TotalPages int
	Pagination []int
	Payload    jwt.Payload
}
