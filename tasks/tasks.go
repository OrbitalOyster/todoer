package tasks

import (
	"fmt"
	"log"
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
	return maxId + 1
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

func GetById(id string) (Task, error) {
	var emptyResult Task
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
