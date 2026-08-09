package collection

type ItemId interface {
	string | int64
}

type Item[T ItemId] interface {
	GetId() T
}

type Collection interface {
	Get()
	Post()
	Patch()
	Put()
}
