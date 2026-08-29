package tasks

import (
	"fmt"
	"strings"
	"time"
	"todoer/collection"
)

type Task[T TaskFieldName] struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Category    string     `yaml:"category"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	ReadOnly    bool       `yaml:"read_only"`
}

func mustBeTask[T TaskFieldName](item collection.Item[T]) Task[T] {
	result, ok := item.(Task[T])
	if !ok {
		panic("Type assertion failed")
	} else {
		return result
	}
}

func (task Task[TaskFiledName]) Field(field TaskFiledName) any {
	switch field {
	case TaskFiledName(Id):
		return task.Id
	case TaskFiledName(User):
		return task.User
	case TaskFiledName(Category):
		return task.Category
	case TaskFiledName(Datetime):
		return task.Datetime
	case TaskFiledName(Description):
		return task.Description
	case TaskFiledName(Status):
		return task.Status
	case TaskFiledName(ReadOnly):
		return task.ReadOnly
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (task Task[TaskFiledName]) MoreThan(field TaskFiledName, value any) bool {
	switch field {
	case TaskFiledName(Datetime):
		return task.Datetime.After(value.(time.Time))
	case TaskFiledName(Description):
		return task.Description > value.(string)
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (task Task[TaskFiledName]) LessThan(field TaskFiledName, value any) bool {
	switch field {
	case TaskFiledName(Datetime):
		valueTime, ok := value.(time.Time)
		if !ok {
			panic ("Type assert failed")
		}
		return task.Datetime.Before(valueTime)
	case TaskFiledName(Description):
		valueString, ok := value.(string)
		if !ok {
			panic ("Type assert failed")
		}
		return task.Description < valueString
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (task Task[TaskFiledName]) Filter(field TaskFiledName, value any) bool {
	switch field {
	case TaskFiledName(Id):
		filterInt, ok := value.(int)
		if !ok {
			panic(fmt.Sprintf("Invalid filter: %v", value))
		}
		return filterInt == task.Id
	case TaskFiledName(Description):
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

/*
func (task *Task[T]) Patch(field T, value any) {
	switch field {
	case T(Id):
	case T(User):
	case T(Category):
	case T(Datetime):
	case T(Description):
	case T(Status):
	case T(ReadOnly):
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}
*/

/* Extra handler for converting status string to TaskStatus */
func (status *TaskStatus) UnmarshalYAML(unmarshal func(any) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	*status = ParseStatus(str)
	return nil
}
