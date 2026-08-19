package collection

import "slices"

type ItemFieldId uint

type Item interface {
	Compare(otherItem Item, field ItemFieldId) int
	Filter(s string, field ItemFieldId) bool
}

type Collection struct {
	Items []Item
}

func (collection Collection) SortBy(field ItemFieldId) Collection {
	result := Collection{Items: slices.Clone(collection.Items)}
	slices.SortFunc(result.Items, func(a Item, b Item) int { return a.Compare(b, field) })
	return result
}

func (collection Collection) Reverse() Collection {
	result := Collection{Items: slices.Clone(collection.Items)}
	slices.Reverse(result.Items)
	return result
}

func (collection Collection) Filter(s string, field ItemFieldId) Collection {
	result := Collection{Items: []Item{}}
	for _, item := range collection.Items {
		if item.Filter(s, field) {
			result.Items = append(result.Items, item)
		}
	}
	return result
}
