package tasks

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"todoer/collection"
)

var list collection.Collection[TaskFieldName]

func GetAll() collection.Collection[TaskFieldName] {
	return list
}

func Load(newList []Task[TaskFieldName]) {
	for _, task := range newList {
		list.Add(&task)
	}
}

func getNextId() int {
	/* No tasks */
	if list.Length() == 0 {
		return 1
	}
	/* Find biggest id, add 1 */
	maxId, ok := list.Max(Id).Field(Id).(int)
	if !ok {
		panic("Major screwup")
	}
	return maxId
}

func Add(user string, description string) {
	now := time.Now()
	newTask := Task[TaskFieldName]{
		Id:          getNextId(),
		User:        user,
		Description: description,
		Datetime:    now,
		Status:      InProgress,
	}
	list.Add(&newTask)
	log.Printf("New task: \"%s\"", newTask.Description)
}

func GetById[T int | string](idIntOrStr T) (*Task[TaskFieldName], error) {
	/* Check generic type */
	var id int
	switch idAny := any(idIntOrStr).(type) {
	case string:
		idInt, err := strconv.Atoi(idAny)
		/* Unparseable string */
		if err != nil {
			return nil, fmt.Errorf("Invalid id string: %s", idAny)
		}
		id = idInt
	case int:
		id = idAny
	default:
		panic(fmt.Sprintf("Invalid type: %v", idAny))
	}

	filtered := list.Filter(Id, id)
	if filtered.Length() == 0 {
		return nil, fmt.Errorf("Task not found: %v", id)
	}
	if filtered.Length() != 1 {
		return nil, fmt.Errorf("More than one task found: %v", id)
	}
	result, ok := filtered.First().(*Task[TaskFieldName])
	if !ok {
		panic("Major screwup")
	}
	return result, nil
}

func (task *Task[TaskFieldName]) SetDescription(description string) error {
	task.Description = description
	log.Printf("Set task #%d description to \"%s\"", task.Id, task.Description)
	return nil
}

func (task *Task[TaskFieldName]) SetUser(user string) error {
	task.User = user
	log.Printf("Set task #%d user to \"%s\"", task.Id, task.User)
	return nil
}

func (task *Task[TaskFieldName]) SetStatus(status TaskStatus) error {
	task.Status = status
	log.Printf("Set task #%d status to \"%s\"", task.Id, task.Status)
	return nil
}

func (task *Task[TaskFieldName]) SetReadOnly(ro bool) error {
	log.Printf("Set task #%d read only to %t", task.Id, ro)
	task.ReadOnly = ro
	return nil
}

func Delete(id int) {
	list.Delete(Id, id)
}
