package collection

type ItemId interface {
	string | int64
}

type ItemFieldValue any

type Item[T ItemId] interface {
	GetId() T
	GetValue(field string) ItemFieldValue
}

type Collection[T ItemId] interface {
	// SortBy(field string, desc bool) Collection[T]
	// FilterBy(field string, value ItemFieldValue) Collection[T]
	// Get(pageSize uint, page uint) (list []Item[T], totalPages uint, actualPage uint)
	GetOne(index T) Item[T]
	Post(item Item[T])
	GetSize() int
	// Patch(id T, field string, value ItemFieldValue)
	// Put()
}
