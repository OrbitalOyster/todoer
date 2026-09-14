package tasks

import (
	"fmt"
	"strings"
)

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

func (fieldName TaskFieldName) String() string {
	switch fieldName {
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
		panic(fmt.Sprintf("Invalid TaskFieldName: %d", fieldName))
	}
}

func ParseTaskFieldName(fieldName string) TaskFieldName {
	switch strings.ToLower(fieldName) {
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
		panic(fmt.Sprintf("Invalid TaskFieldName: %s", fieldName))
	}
}
