package collection

import (
	"log"
)

/* ========================================================================== */

type MyItem struct {
	id  int64
	foo int
	bar string
}

func (item MyItem) GetId() int64 {
	return item.id
}

func (item MyItem) GetValue(f string) ItemFieldValue {
	switch f {
	case "foo":
		return item.foo
	case "bar":
		return item.bar
	}
	panic("Oh no!")
}

/* ========================================================================== */

type MyCollection struct {
	List []Item[int64]
}

func (collection *MyCollection) Post(item Item[int64]) {
	collection.List = append(collection.List, item)
}

func (collection MyCollection) GetSize() int {
	return len(collection.List)
}

func (collection MyCollection) GetOne(index int64) Item[int64] {
	return collection.List[index]
}

func DebugCollection[T int64](collection Collection[T]) Item[T] {
	if collection.GetSize() > 0 {
		return collection.GetOne(0)
	} else {
		var emptyResult Item[T]
		return emptyResult
	}
}

func Run() {
	log.Println("Collection test")
	i1 := MyItem{1, 7, "Hello"}
	i2 := MyItem{2, 77, "World"}
	i3 := MyItem{3, 777, "boo"}
	c := &MyCollection{}
	c.Post(i1)
	c.Post(i2)
	c.Post(i3)
	log.Printf("%#v", DebugCollection(c))
}
