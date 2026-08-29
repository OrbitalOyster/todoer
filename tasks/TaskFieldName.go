package tasks

import "strings"

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
