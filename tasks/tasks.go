package tasks

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"todoer/collection"
	"todoer/utils"
)

// var list []Task
var list collection.Collection[TaskFieldName]

func GetAll() collection.Collection[TaskFieldName]{
	return list
}

func Load(newList []Task[TaskFieldName]) {
	// list = newList
	for _, t := range newList {
		list.Items = append(list.Items, t)
	}
}

func getNextId() int {
	/* No tasks */
	// if len(list) == 0 {
	if list.Length() == 0 {
		return 1
	}
	/* Find biggest id, add 1 */
	/*
		result := slices.MaxFunc(list, func(a, b Task) int {
			return cmp.Compare(a.Id, b.Id)
		})
	*/
	max, ok := list.Max(Id).(Task[TaskFieldName])
	if !ok {
		panic("Major screwup")
	}
	return max.Id + 1
}

func Add(user string, description string) {
	now := time.Now()
	result := Task[TaskFieldName]{
		Id:          getNextId(),
		User:        user,
		Description: description,
		Datetime:    now,
		Status:      InProgress,
	}
	// list = append(list, result)
	list.Add(result)
	log.Printf("New task: \"%s\"", result.Description)
}

func Get(fromDateStr string, toDateStr string,
	searchBy string,
	page int, pageSize int,
	sortBy utils.SortableField, sortAsc bool) ([]Task[TaskFieldName], int, int) {

	result := list.Clone()

	/* Date */
	fromDate, err := time.Parse(utils.HTMLDateFormat, fromDateStr)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	toDate, err := time.Parse(utils.HTMLDateFormat, toDateStr)
	/* Should not happen */
	if err != nil {
		panic(err)
	}
	result = result.
		MoreThan(Datetime, Task[TaskFieldName]{Datetime: fromDate}).
		/* "Not after 20/03/2026" means "Not after 20/03/2026 23:59:59"  */
		LessThan(Datetime, Task[TaskFieldName]{Datetime: toDate.Add(time.Hour*24 - time.Second)}).
		Filter(Description, searchBy)
	/* Sorting */
	switch sortBy {
	case utils.Description:
		result = result.SortBy(Description)
	case utils.Datetime:
		result = result.SortBy(Datetime)
	default:
	}
	if !sortAsc {
		result.Reverse()
	}

	result, actualPage, numberOfPages := result.GetPage(uint(page), uint(pageSize))

	/*
		var actualResult []Task
		for _, i := range result.Items {
			actualResult = append(actualResult, i.(Task))
		}
	*/

	return collection.AssertType[Task[TaskFieldName]](result), int(numberOfPages), int(actualPage)
	// result := slices.Clone(list)
	/* Date */
	// fromDate, err := time.Parse(utils.HTMLDateFormat, fromDateStr)
	/* Should not happen */
	// if err != nil {
	// 	panic(err)
	// }
	// toDate, err := time.Parse(utils.HTMLDateFormat, toDateStr)
	/* Should not happen */
	// if err != nil {
	// 	panic(err)
	// }
	// result = slices.DeleteFunc(result, func(t Task) bool {
	/* "Not after 20/03/2026" means "Not after 20/03/2026 23:59:59"  */
	// 	return t.Datetime.Before(fromDate) || t.Datetime.After(toDate.Add(time.Hour*24-time.Second))
	// })
	/* Search */
	// if searchBy != "" {
	// 	result = slices.DeleteFunc(result, func(t Task) bool {
	// 		return !strings.Contains(t.Description, searchBy)
	// 	})
	// }
	/* Number of tasks after all filtering */
	// total := len(result)
	/* Nothing found - stop */
	// if total == 0 {
	// 	return nil, 1, 1
	// }
	/* Sorting */
	// switch sortBy {
	// case utils.Description:
	// 	slices.SortFunc(result, func(t1, t2 Task) int {
	// 		return cmp.Compare(t1.Description, t2.Description)
	// 	})
	// case utils.Datetime:
	// 	slices.SortFunc(result, func(t1, t2 Task) int {
	// 		return t1.Datetime.Compare(t2.Datetime)
	// 	})
	// default:
	// }
	/* On reverse order */
	// if !sortAsc {
	// 	slices.Reverse(result)
	// }
	/* Pagination */
	// totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	// if page >= totalPages {
	// 	page = totalPages
	// }
	// if page <= 0 {
	// 	page = 1
	// }
	/* Final result */
	// startInd := pageSize * (page - 1)
	// endInd := min(startInd+pageSize, total)
	// return result[startInd:endInd], totalPages, page
}

func getById(id int) (*Task[TaskFieldName], error) {
	/*
		ind := slices.IndexFunc(list, func(t Task) bool {
			return t.Id == id
		})
		if ind == -1 {
			return nil, fmt.Errorf("Task not found: %d", id)
		}
		return &list[ind], nil
	*/
	filtered := list.Filter(Id, id)
	if filtered.Length() == 0 {
		return nil, fmt.Errorf("Task not found: %d", id)
	}
	if filtered.Length() != 1 {
		return nil, fmt.Errorf("More than one task found: %d", id)
	}
	result, ok := filtered.First().(Task[TaskFieldName])
	if !ok {
		panic("Major screwup")
	}
	return &result, nil
}

/* Generic function, accepts id as int or string */
func GetById[T int | string](id T) (*Task[TaskFieldName], error) {
	switch idAny := any(id).(type) {
	case int:
		return getById(idAny)
	case string:
		idInt, err := strconv.Atoi(idAny)
		/* Unparseable string */
		if err != nil {
			return nil, fmt.Errorf("Invalid task identifier: \"%s\"", idAny)
		}
		return getById(idInt)
	default:
		/* Major screwup */
		panic("Invalid task type")
	}
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
	/*
		ind := slices.IndexFunc(list, func(t Task) bool {
			return t.Id == id
		})
		if ind == -1 {
			return fmt.Errorf("Task not found: %d", id)
		}
		list = slices.Delete(list, ind, ind+1)
		log.Printf("Deleted task #%d", id)
		return nil
	*/
	list.Delete(Id, id)
}
