package collection

import (
	"math"
	"slices"
)

type FieldName uint

type Item interface {
	MoreThan(field FieldName, item Item) bool
	LessThan(field FieldName, item Item) bool
	Filter(field FieldName, s any) bool
}

type Collection struct {
	Items []Item
}

func (collection Collection) Length() int {
	return len(collection.Items)
}

func (collection Collection) Clone() Collection {
	return Collection{Items: slices.Clone(collection.Items)}
}

func (collection *Collection) Add(newItem Item) {
	collection.Items = append(collection.Items, newItem)
}

func (collection *Collection) Delete(field FieldName, filter any) {
	collection.Items = slices.DeleteFunc(collection.Items, func(item Item) bool {
		return item.Filter(field, filter)
	})
}

func (collection Collection) First() Item {
	if len(collection.Items) < 1 {
		panic("Empty array")
	}
	return collection.Items[0]
}

func compare(a Item, b Item, field FieldName) int {
	switch {
	case a.LessThan(field, b):
		return -1
	case a.MoreThan(field, b):
		return 1
	default:
		return 0
	}
}

func (collection Collection) Max(field FieldName) Item {
	return slices.MaxFunc(collection.Items, func(a, b Item) int {
		return compare(a, b, field)
	})
}

func (collection Collection) SortBy(field FieldName) Collection {
	slices.SortFunc(collection.Items, func(a Item, b Item) int {
		return compare(a, b, field)
	})
	return collection
}

func (collection Collection) Reverse() Collection {
	slices.Reverse(collection.Items)
	return collection
}

func (collection Collection) Filter(field FieldName, filter any) Collection {
	result := Collection{Items: []Item{}}
	for _, item := range collection.Items {
		if item.Filter(field, filter) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}

func (collection Collection) MoreThan(field FieldName, item Item) Collection {
	result := Collection{Items: []Item{}}
	for _, i := range collection.Items {
		if i.MoreThan(field, item) {
			result.Items = append(result.Items, i)
		}
	}
	return result
}

func (collection Collection) LessThan(field FieldName, item Item) Collection {
	result := Collection{Items: []Item{}}
	for _, i := range collection.Items {
		if i.LessThan(field, item) {
			result.Items = append(result.Items, i)
		}
	}
	return result
}

func (collection Collection) GetPage(page uint, pageSize uint) (Collection, uint, uint) {
	length := len(collection.Items)
	numberOfPages := uint(math.Ceil(float64(length) / float64(pageSize)))
	if page >= numberOfPages {
		page = numberOfPages
	}
	if page == 0 {
		page = 1
	}
	startInd := pageSize * (page - 1)
	endInd := min(startInd+pageSize, uint(length))
	result := Collection{Items: slices.Clone(collection.Items)[startInd:endInd]}
	return result, page, numberOfPages
}

func AssertType[T any](collection Collection) []T {
	var result []T
	for _, i := range collection.Items {
		t, ok := i.(T)
		if !ok {
			panic("Type assert failed")
		}
		result = append(result, t)
	}
	return result
}
