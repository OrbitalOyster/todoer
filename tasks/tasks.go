package tasks

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"todoer/collection"
)

var All collection.Collection[TaskField]

func Load(newList []Task) {
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
	newTask := Task{
		Id:          getNextId(),
		User:        user,
		Description: description,
		Datetime:    now,
		Status:      InProgress,
	}
	All.Add(&newTask)
	log.Printf("New task: \"%s\"", newTask.Description)
}

func GetById[T int | string](idIntOrStr T) (Task, error) {
	var (
		id          string
		emptyResult Task
	)
	switch idAny := any(idIntOrStr).(type) {
	case string:
	id = idAny
	/*
		idInt, err := strconv.Atoi(idAny)
		if err != nil {
			return emptyResult, fmt.Errorf("Invalid id string: %s", idAny)
		}
		id = idInt
		*/
	case int:
	id = strconv.Itoa(idAny)
	/*
		id = idAny
		*/
	default:
		panic(fmt.Sprintf("Invalid type: %v", idAny))
	}
	filtered := All.Filter(Id, []string{id})
	if filtered.Length() == 0 {
		return emptyResult, fmt.Errorf("Task #%s not found", id)
	}
	if filtered.Length() != 1 {
		return emptyResult, fmt.Errorf("More than one task %s found", id)
	}
	result, ok := filtered.First().(*Task)
	if !ok {
		panic("Major screwup")
	}
	return *result, nil
}

func Patch(ids []int, field TaskField, value string) (patched uint, errors []error) {
	filterBy := make([]string, len(ids))
	for i, c := range ids {
		filterBy[i] = strconv.Itoa(c)
	}
	patched += All.FilterAndPatch(Id, filterBy, field, value)
	return
}

func DeleteOne(id int) {
	All.Delete(Id, []string{strconv.Itoa(id)})
}
