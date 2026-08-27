package tasks

import (
	"fmt"
	"strings"
	"time"
	"todoer/collection"
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

type Task[T TaskFieldName] struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Category    string     `yaml:"category"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	ReadOnly    bool       `yaml:"read_only"`
}

type TaskFieldName uint

const (
	Id TaskFieldName = iota
	User
	Category
	Datetime
	Description
	Status
	ReadOnly
)

func (field TaskFieldName) String() string {
	switch field {
	case Id:
		return "Id"
	case User:
		return "User"
	case Category:
		return "Category"
	case Datetime:
		return "Datetime"
	case Description:
		return "Description"
	case Status:
		return "Status"
	case ReadOnly:
		return "ReadOnly"
	default:
		panic("Invalid type")
	}
}

func ParseTaskFieldName(s string) TaskFieldName {
	switch strings.ToLower(s) {
	case "id":
		return Id
	case "user":
		return User
	case "category":
		return Category
	case "datetime":
		return Datetime
	case "description":
		return Description
	case "status":
		return Status
	case "readonly":
		return ReadOnly
	default:
		panic("Invalid type")
	}
}

func mustBeTask[T TaskFieldName](item collection.Item[T]) Task[T] {
	result, ok := item.(Task[T])
	if !ok {
		panic("Type assertion failed")
	} else {
		return result
	}
}

func (task Task[T]) MoreThan(field T, item collection.Item[T]) bool {
	otherTask := mustBeTask(item)
	switch field {
	case T(Datetime):
		return task.Datetime.After(otherTask.Datetime)
	case T(Description):
		return task.Description > otherTask.Description
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (task Task[T]) LessThan(field T, item collection.Item[T]) bool {
	otherTask := mustBeTask(item)
	switch field {
	case T(Datetime):
		return task.Datetime.Before(otherTask.Datetime)
	case T(Description):
		return task.Description < otherTask.Description
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (task Task[T]) Filter(field T, value any) bool {
	switch field {
	case T(Id):
		filterInt, ok := value.(int)
		if !ok {
			panic(fmt.Sprintf("Invalid filter: %v", value))
		}
		return filterInt == task.Id
	case T(Description):
		filterString, ok := value.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid filter: %v", value))
		}
		if filterString == "" {
			return true
		}
		return strings.Contains(task.Description, filterString)
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
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
