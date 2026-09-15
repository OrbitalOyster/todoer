package tasks

import (
	"fmt"
	"strings"
	"time"
)

type Task[T TaskField] struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Category    string     `yaml:"category"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	ReadOnly    bool       `yaml:"read_only"`
}

func (task Task[T]) Field(field TaskField) any {
	switch field {
	case Id:
		return task.Id
	case User:
		return task.User
	case Category:
		return task.Category
	case Datetime:
		return task.Datetime
	case Description:
		return task.Description
	case Status:
		return task.Status
	case ReadOnly:
		return task.ReadOnly
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func (task *Task[T]) Patch(field TaskField, value any) {
	switch field {
	case User:
		valueStr, ok := value.(string)
		if !ok {
			panic("Type assertion failed")
		}
		task.User = valueStr
	case Description:
		valueString, ok := value.(string)
		if !ok {
			panic("Type assertion failed")
		}
		task.Description = valueString
	case Status:
		valueStatus, ok := value.(TaskStatus)
		if !ok {
			panic("Type assertion failed")
		}
		task.Status = valueStatus
	case ReadOnly:
		valueBool, ok := value.(bool)
		if !ok {
			panic("Type assertion failed")
		}
		task.ReadOnly = valueBool
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func (task Task[T]) MoreThan(field TaskField, value any) bool {
	switch field {
	case Datetime:
		return task.Datetime.After(value.(time.Time))
	case Description:
		return task.Description > value.(string)
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func (task Task[T]) LessThan(field TaskField, value any) bool {
	switch field {
	case Datetime:
		valueTime, ok := value.(time.Time)
		if !ok {
			panic("Type assertion failed")
		}
		return task.Datetime.Before(valueTime)
	case Description:
		valueString, ok := value.(string)
		if !ok {
			panic("Type assertion failed")
		}
		return task.Description < valueString
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func (task Task[T]) Filter(field TaskField, value any) bool {
	switch field {
	case Id:
		filterInt, ok := value.(int)
		if !ok {
			panic(fmt.Sprintf("Invalid filter: %#v", value))
		}
		return filterInt == task.Id
	case Description:
		filterString, ok := value.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid filter: %#v", value))
		}
		if filterString == "" {
			return true
		}
		return strings.Contains(task.Description, filterString)
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
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
