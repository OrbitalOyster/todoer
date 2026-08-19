package collection

import (
	"slices"
)

type ItemFieldId uint

type Item interface {
	MoreThan(field ItemFieldId, otherItem Item) bool
	LessThan(field ItemFieldId, otherItem Item) bool
	Filter(field ItemFieldId, s string) bool
}

type Collection struct {
	Items []Item
}

func (collection Collection) SortBy(field ItemFieldId) Collection {
	result := Collection{Items: slices.Clone(collection.Items)}
	slices.SortFunc(result.Items, func(a Item, b Item) int {
		if a.LessThan(field, b) {
			return -1
		} else if a.MoreThan(field, b) {
			return 1
		}
		return 0
	})
	return result
}

func (collection Collection) Reverse() Collection {
	result := Collection{Items: slices.Clone(collection.Items)}
	slices.Reverse(result.Items)
	return result
}

func (collection Collection) Filter(field ItemFieldId, s string) Collection {
	result := Collection{Items: []Item{}}
	for _, item := range collection.Items {
		if item.Filter(field, s) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}

func (collection Collection) MoreThan(field ItemFieldId, otherItem Item) Collection {
	result := Collection{Items: []Item{}}
	for _, item := range collection.Items {
		if item.MoreThan(field, otherItem) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}

func (collection Collection) LessThan(field ItemFieldId, otherItem Item) Collection {
	result := Collection{Items: []Item{}}
	for _, item := range collection.Items {
		if item.LessThan(field, otherItem) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}
