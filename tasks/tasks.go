package tasks

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"todoer/collection"
)

var list collection.Collection[TaskField]

func GetAll() collection.Collection[TaskField] {
	return list
}

func Load(newList []Task[TaskField]) {
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
	newTask := Task[TaskField]{
		Id:          getNextId(),
		User:        user,
		Description: description,
		Datetime:    now,
		Status:      InProgress,
	}
	list.Add(&newTask)
	log.Printf("New task: \"%s\"", newTask.Description)
}

func GetById[T int | string](idIntOrStr T) (*Task[TaskField], error) {
	var id int
	switch idAny := any(idIntOrStr).(type) {
	case string:
		idInt, err := strconv.Atoi(idAny)
		if err != nil {
			return nil, fmt.Errorf("Invalid id string: %s", idAny)
		}
		id = idInt
	case int:
		id = idAny
	default:
		panic(fmt.Sprintf("Invalid type: %v", idAny))
	}

	filtered := list.Filter(Id, []any{id})
	if filtered.Length() == 0 {
		return nil, fmt.Errorf("Task not found: %v", id)
	}
	if filtered.Length() != 1 {
		return nil, fmt.Errorf("More than one task found: %v", id)
	}
	result, ok := filtered.First().(*Task[TaskField])
	if !ok {
		panic("Major screwup")
	}
	return result, nil
}

func FilterAndPatch(field TaskField, filter []any, fieldToPatch TaskField, value any) uint {
	return list.FilterAndPatch(field, filter, fieldToPatch, value)
}

func DeleteOne(id int) {
	list.Delete(Id, []any{id})
}
