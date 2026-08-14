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

func (item MyItem) GetValue(field string) ItemFieldValue {
	switch field {
	case "foo":
		return item.foo
	case "bar":
		return item.bar
	}
	panic("Oh no!")
}

/* ========================================================================== */

type MyCollection struct {
	List []MyItem
}

func (collection *MyCollection) Post(item Item) {
	collection.List = append(collection.List, item.(MyItem))
}

func (collection MyCollection) GetSize() int {
	return len(collection.List)
}

func (collection MyCollection) GetOne(index int) Item {
	return collection.List[index]
}

func DebugCollection(collection Collection) Item {
	if collection.GetSize() > 0 {
		return collection.GetOne(0)
	} else {
		var emptyResult Item
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
