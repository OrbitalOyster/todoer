package tasks

import "time"

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
		panic("Invalid TaskStatus")
	}
}

func ParseStatus(status string) TaskStatus {
	switch status {
	case "InProgress":
		return InProgress
	case "Done":
		return Done
	case "Failed":
		return Failed
	default:
		panic("Invalid TaskStatus")
	}
}

type Task struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
}

func (status *TaskStatus) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	*status = ParseStatus(str)
	return nil
}
