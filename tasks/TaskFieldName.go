package tasks

import (
	"fmt"
	"strings"
)

type TaskField uint

const (
	Id TaskField = iota
	User
	Category
	Datetime
	Description
	Status
	ReadOnly
)

func (field TaskField) String() string {
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
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}

func ParseTaskField(field string) (TaskField, error) {
	switch strings.ToLower(field) {
	case "id":
		return Id, nil
	case "user":
		return User, nil
	case "category":
		return Category, nil
	case "datetime":
		return Datetime, nil
	case "description":
		return Description, nil
	case "status":
		return Status, nil
	case "readonly":
		return ReadOnly, nil
	default:
		return 0, fmt.Errorf("Invalid field: %#v", field)
	}
}
