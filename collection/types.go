package collection

import (
	"math"
	"slices"
)

type FieldName interface {
	~uint
}

type Item[T FieldName] interface {
	MoreThan(field T, item Item[T]) bool
	LessThan(field T, item Item[T]) bool
	Filter(field T, s any) bool
}

type Collection[T FieldName] struct {
	Items []Item[T]
}

func (collection Collection[T]) Length() int {
	return len(collection.Items)
}

func (collection Collection[T]) Clone() Collection[T] {
	return Collection[T]{Items: slices.Clone(collection.Items)}
}

func (collection *Collection[T]) Add(newItem Item[T]) {
	collection.Items = append(collection.Items, newItem)
}

func (collection *Collection[T]) Delete(field T, filter any) {
	collection.Items = slices.DeleteFunc(collection.Items, func(item Item[T]) bool {
		return item.Filter(field, filter)
	})
}

func (collection Collection[T]) First() Item[T] {
	if len(collection.Items) < 1 {
		panic("Empty array")
	}
	return collection.Items[0]
}

func compare[T FieldName](a Item[T], b Item[T], field T) int {
	switch {
	case a.LessThan(field, b):
		return -1
	case a.MoreThan(field, b):
		return 1
	default:
		return 0
	}
}

func (collection Collection[T]) Max(field T) Item[T] {
	return slices.MaxFunc(collection.Items, func(a, b Item[T]) int {
		return compare(a, b, field)
	})
}

func (collection Collection[T]) SortBy(field T) Collection[T] {
	slices.SortFunc(collection.Items, func(a Item[T], b Item[T]) int {
		return compare(a, b, field)
	})
	return collection
}

func (collection Collection[T]) Reverse() Collection[T] {
	slices.Reverse(collection.Items)
	return collection
}

func (collection Collection[T]) Filter(field T, filter any) Collection[T] {
	result := Collection[T]{Items: []Item[T]{}}
	for _, item := range collection.Items {
		if item.Filter(field, filter) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}

func (collection Collection[T]) MoreThan(field T, item Item[T]) Collection[T] {
	result := Collection[T]{Items: []Item[T]{}}
	for _, i := range collection.Items {
		if i.MoreThan(field, item) {
			result.Items = append(result.Items, i)
		}
	}
	return result
}

func (collection Collection[T]) LessThan(field T, item Item[T]) Collection[T] {
	result := Collection[T]{Items: []Item[T]{}}
	for _, i := range collection.Items {
		if i.LessThan(field, item) {
			result.Items = append(result.Items, i)
		}
	}
	return result
}

func (collection Collection[T]) GetPage(page uint, pageSize uint) (Collection[T], uint, uint) {
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
	result := Collection[T]{Items: slices.Clone(collection.Items)[startInd:endInd]}
	return result, page, numberOfPages
}

func AssertType[T any, Y FieldName](collection Collection[Y]) []T {
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
