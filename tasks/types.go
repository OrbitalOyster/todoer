package tasks

import (
	"fmt"
	"strings"
	"time"
)

type TaskStatus int

const (
	InProgress TaskStatus = iota
	Done
	Failed
)

func (status TaskStatus) String() string {
	switch status {
	case InProgress:
		return "InProgress"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	default:
		/* Major screwup */
		panic("Invalid TaskStatus")
	}
}

func ParseStatus(status string) TaskStatus {
	switch strings.ToLower(status) {
	case "inprogress":
		return InProgress
	case "done":
		return Done
	case "failed":
		return Failed
	default:
		panic(fmt.Sprintf("Invalid TaskStatus: %s", status))
	}
}

type Task struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	ReadOnly    bool       `yaml:"read_only"`
}

/* Extra handler for converting status string to TaskStatus */
func (status *TaskStatus) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	*status = ParseStatus(str)
	return nil
}
