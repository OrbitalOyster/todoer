package wad

import (
	"todoer/tasks"
	"todoer/users"
)

type WAD struct {
	Users      []users.User `yaml:"users"`
	Categories []string     `yaml:"categories"`
	Tasks      []tasks.Task[tasks.TaskFieldName] `yaml:"tasks"`
}
