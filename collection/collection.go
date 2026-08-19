package collection

import (
	"fmt"
	"log"
	"strings"
)

const (
	Id ItemFieldId = iota
	Foo
	Bar
)

type MyItem struct {
	Id  int64
	Foo int
	Bar string
}

func (myItem MyItem) MoreThan(field ItemFieldId, otherItem Item) bool {
	otherMyItem, ok := otherItem.(MyItem)
	if !ok {
		panic("Type assertion failed")
	}
	switch field {
	case Id:
		return myItem.Id > otherMyItem.Id
	case Foo:
		return myItem.Foo > otherMyItem.Foo
	case Bar:
		return myItem.Bar > otherMyItem.Bar
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (myItem MyItem) LessThan(field ItemFieldId, otherItem Item) bool {
	otherMyItem, ok := otherItem.(MyItem)
	if !ok {
		panic("Type assertion failed")
	}
	switch field {
	case Id:
		return myItem.Id < otherMyItem.Id
	case Foo:
		return myItem.Foo < otherMyItem.Foo
	case Bar:
		return myItem.Bar < otherMyItem.Bar
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func (myItem MyItem) Filter(field ItemFieldId, s string) bool {
	switch field {
	case Bar:
		return strings.Contains(myItem.Bar, s)
	default:
		panic(fmt.Sprintf("Invalid field: %d", field))
	}
}

func Run() {
	log.Println("Collection test")
	var _ Item = MyItem{}
	MyCollection := Collection{
		Items: []Item{
			MyItem{Id: 1, Foo: 3, Bar: "Lorem"},
			MyItem{Id: 2, Foo: 8, Bar: "ipsum"},
			MyItem{Id: 3, Foo: 4, Bar: "dolor"},
			MyItem{Id: 4, Foo: 9, Bar: "sit"},
			MyItem{Id: 5, Foo: 3, Bar: "amet"},
			MyItem{Id: 6, Foo: 7, Bar: "consectetur"},
			MyItem{Id: 7, Foo: 5, Bar: "adipiscing"},
			MyItem{Id: 8, Foo: 2, Bar: "elit"},
			MyItem{Id: 9, Foo: 1, Bar: "sed"},
			MyItem{Id: 10, Foo: 1, Bar: "do"},
		},
	}
	log.Printf("%#v\n", MyCollection.SortBy(Foo).MoreThan(Foo, MyItem{Foo: 5}))
	log.Printf("%#v\n", MyCollection)
}
