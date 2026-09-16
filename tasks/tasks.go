package tasks

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"todoer/collection"
)

var All collection.Collection[TaskField]

func Load(newList []Task[TaskField]) {
	for _, task := range newList {
		All.Add(&task)
	}
}

func getNextId() int {
	/* No tasks */
	if All.Length() == 0 {
		return 1
	}
	/* Find biggest id, add 1 */
	maxId, ok := All.Max(Id).Field(Id).(int)
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
	All.Add(&newTask)
	log.Printf("New task: \"%s\"", newTask.Description)
}

func GetById[T int | string](idIntOrStr T) (Task[TaskField], error) {
	var (
		id          int
		emptyResult Task[TaskField]
	)
	switch idAny := any(idIntOrStr).(type) {
	case string:
		idInt, err := strconv.Atoi(idAny)
		if err != nil {
			return emptyResult, fmt.Errorf("Invalid id string: %s", idAny)
		}
		id = idInt
	case int:
		id = idAny
	default:
		panic(fmt.Sprintf("Invalid type: %v", idAny))
	}
	filtered := All.Filter(Id, []any{id})
	if filtered.Length() == 0 {
		return emptyResult, fmt.Errorf("Task #%d not found", id)
	}
	if filtered.Length() != 1 {
		return emptyResult, fmt.Errorf("More than one task %d found", id)
	}
	result, ok := filtered.First().(*Task[TaskField])
	if !ok {
		panic("Major screwup")
	}
	return *result, nil
}

func Patch(ids []int, field TaskField, value any) (patched uint, errors []error) {
	filterBy := make([]any, len(ids))
	for i, c := range ids {
		filterBy[i] = c
	}
	switch field {
	case User:
		userStr, ok := value.(string)
		if !ok {
			panic("Major screwup")
		}
		patched += All.FilterAndPatch(Id, filterBy, User, userStr)
	case Description:
		descriptionStr, ok := value.(string)
		if !ok {
			panic("Major screwup")
		}
		patched += All.FilterAndPatch(Id, filterBy, Description, descriptionStr)
	case Status:
		statusStr, ok := value.(string)
		if !ok {
			panic("Major screwup")
		}
		status, err := ParseStatus(statusStr)
		if err != nil {
			errors = append(errors, err)
		} else {
			patched += All.FilterAndPatch(Id, filterBy, Status, status)
		}
	case ReadOnly:
		readOnly, ok := value.(bool)
		if !ok {
			panic("Major screwup")
		}
		patched += All.FilterAndPatch(Id, filterBy, ReadOnly, readOnly)
	}
	return
}

func DeleteOne(id int) {
	All.Delete(Id, []any{id})
}
