package tasks

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"todoer/utils"
)

type Task struct {
	Id          int        `yaml:"id"`
	User        string     `yaml:"user"`
	Category    string     `yaml:"category"`
	Datetime    time.Time  `yaml:"datetime"`
	Description string     `yaml:"description"`
	Status      TaskStatus `yaml:"status"`
	ReadOnly    bool       `yaml:"read_only"`
}

func ParseValue(field TaskField, value string) (result any, err error) {
	switch field {
	case Id: /* Int */
		return strconv.Atoi(value)
	case User, Description, Category: /* String */
		return value, nil
	case Datetime: /* Time */
		return time.Parse(utils.HTMLDateFormat, value)
	case Status: /* TaskStatus */
		return ParseStatus(value)
	case ReadOnly: /* Bool */
		return strconv.ParseBool(value)
	default:
		panic(fmt.Errorf("Invalid field: %#v", field))
	}
}

func (task Task) Field(field TaskField) any {
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

func (task *Task) Put(user string, description string, readOnly bool) {
	/* TODO: Some error checking */
	task.User = user
	task.Description = description
	task.ReadOnly = readOnly

	log.Printf("Updated task #%d to %#v", task.Id, *task)
}

func (task *Task) Patch(field TaskField, value string) (updated bool, err error) {
	if parsed, err := ParseValue(field, value); err != nil {
		return false, err
	} else {
		switch field {
		case Id:
		case User:
		case Category:
		case Datetime:
		case Description:
		case Status:
			if task.Status != parsed {
				task.Status = parsed.(TaskStatus)
				updated = true
			}
		case ReadOnly:
			if task.ReadOnly != parsed {
				task.ReadOnly = parsed.(bool)
				updated = true
			}
		default:
			panic(fmt.Sprintf("Invalid field: %#v", field))
		}
		if updated {
			log.Printf("Set task #%d field %s to %#v", task.Id, field.String(), value)
		}
	}
	return
}

func (task Task) MoreThan(field TaskField, value any) bool {
	switch field {
	case Id:
		valueInt, ok := value.(int)
		if !ok {
			panic("Type assertion failed")
		}
		return task.Id > valueInt
	case Datetime:
		return task.Datetime.After(value.(time.Time))
	case Description:
		return task.Description > value.(string)
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func (task Task) LessThan(field TaskField, value any) bool {
	switch field {
	case Id:
		valueInt, ok := value.(int)
		if !ok {
			panic("Type assertion failed")
		}
		return task.Id < valueInt
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

func (task Task) Filter(field TaskField, values []string) bool {
	switch field {
	case Id:
		for _, value := range values {
			var id int
			switch idAny := any(value).(type) {
			case string:
				idInt, err := strconv.Atoi(idAny)
				if err != nil {
					panic(fmt.Sprintf("Not a valid id string: %#v", value))
				}
				id = idInt
			case int:
				id = idAny
			default:
				panic(fmt.Sprintf("Invalid type: %v", idAny))
			}
			if id == task.Id {
				return true
			}
		}
		return false
	case Description:
		for _, value := range values {
			if value == "" || strings.Contains(task.Description, value) {
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
