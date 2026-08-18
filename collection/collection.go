package collection

import (
	"cmp"
	"log"
	"slices"
)

type ItemFieldEnum uint

const (
	Id ItemFieldEnum = iota
	Foo
	Bar
)

type Item interface {
	Less(Item, ItemFieldEnum) int
}

type Collection struct {
	Items []Item
}

func (collection Collection) SortBy(field ItemFieldEnum) Collection {
	result := Collection{Items: slices.Clone(collection.Items)}
	return result
}

type MyItem struct {
	Foo int64
	Bar string
}

func (myItem MyItem) Less(otherItem Item, field ItemFieldEnum) int {
	switch field {
	case Foo:
		return cmp.Compare(myItem.Foo, otherItem.(MyItem).Foo)
	case Bar:
		return cmp.Compare(myItem.Bar, otherItem.(MyItem).Bar)
	default:
		panic("Oh no!")
	}
}

func Run() {
	log.Println("Collection test")
	var _ Item = MyItem{}
	MyCollection := Collection{
		Items: []Item{
			MyItem{
				Foo: 1,
				Bar: "1",
			},
		},
	}
	MyCollection.SortBy(Foo)
}
