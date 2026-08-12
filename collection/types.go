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
	SortBy(field string, desc bool) Collection[T]
	FilterBy(field string, value ItemFieldValue) Collection[T]
	Get(pageSize uint, page uint) (list []Item[T], totalPages uint, actualPage uint)
	GetOne(id T) Item[T]
	Post(item Item[T])
	Patch(id T, field string, value ItemFieldValue)
	Put()
}
