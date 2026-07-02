package wad

import (
	"todoer/tasks"
	"todoer/users"
)

type WAD struct {
	Users []users.User `yaml:"users"`
	Tasks []tasks.Task `yaml:"tasks"`
}
