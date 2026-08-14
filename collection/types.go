package collection

type ItemFieldValue any

type Item interface {
	GetId() int64
	GetValue(field string) ItemFieldValue
}

type Collection interface {
	// SortBy(field string, desc bool) Collection[T]
	// FilterBy(field string, value ItemFieldValue) Collection[T]
	// Get(pageSize uint, page uint) (list []Item[T], totalPages uint, actualPage uint)
	GetOne(index int) Item
	Post(item Item)
	GetSize() int
	// Patch(id T, field string, value ItemFieldValue)
	// Put()
}
