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

func ParseTaskFieldName(field string) TaskField {
	switch strings.ToLower(field) {
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
	case "readonly", "read-only":
		return ReadOnly
	default:
		panic(fmt.Sprintf("Invalid field: %#v", field))
	}
}
