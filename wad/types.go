package wad

import (
	"todoer/tasks"
)

/* TODO: Not here */
type User struct {
	Name string `yaml:"name"`
	Role string `yaml:"role"`
}

type WAD struct {
	Users []User       `yaml:"users"`
	Tasks []tasks.Task `yaml:"tasks"`
}
