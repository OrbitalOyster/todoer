package tasks

import (
	"fmt"
	"log"
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

func (task *Task[T]) Patch(field TaskField, value any) (success bool) {
	switch field {
	case User:
		valueStr, ok := value.(string)
		if !ok {
			panic("Type assertion failed")
		}
		if task.User != valueStr {
			task.User = valueStr
			log.Printf("Set task #%d user to \"%s\"", task.Id, task.User)
			return true
		} else {
			return false
		}
	case Description:
		valueString, ok := value.(string)
		if !ok {
			panic("Type assertion failed")
		}
		if task.Description != valueString {
			task.Description = valueString
			log.Printf("Set task #%d description to \"%s\"", task.Id, task.Description)
			return true
		} else {
			return false
		}
	case Status:
		valueStatus, ok := value.(TaskStatus)
		if !ok {
			panic("Type assertion failed")
		}
		if task.Status != valueStatus {
			task.Status = valueStatus
			log.Printf("Set task #%d status to \"%s\"", task.Id, task.Status)
			return true
		} else {
			return false
		}
	case ReadOnly:
		valueBool, ok := value.(bool)
		if !ok {
			panic("Type assertion failed")
		}
		if task.ReadOnly != valueBool {
			task.ReadOnly = valueBool
			log.Printf("Set task #%d read only to %t", task.Id, task.ReadOnly)
			return true
		} else {
			return false
		}
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

func (task Task[T]) Filter(field TaskField, values []any) bool {
	switch field {
	case Id:
		for _, value := range values {
			filterInt, ok := value.(int)
			if !ok {
				panic(fmt.Sprintf("Invalid filter value: %#v (expected int)", value))
			}
			if filterInt == task.Id {
				return true
			}
		}
		return false
	case Description:
		for _, value := range values {
			filterString, ok := value.(string)
			if !ok {
				panic(fmt.Sprintf("Invalid filter: %#v", values))
			}
			if filterString == "" || strings.Contains(task.Description, filterString) {
				return true
			}
		}
		return false
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

/* Extra handler for converting status string to TaskStatus */
func (status *TaskStatus) UnmarshalYAML(unmarshal func(any) error) error {
	var (
		str string
		err error
	)
	/* Not a string */
	if err = unmarshal(&str); err != nil {
		return err
	}
	/* Not a valid status */
	if *status, err = ParseStatus(str); err != nil {
		return err
	}
	return nil
}
