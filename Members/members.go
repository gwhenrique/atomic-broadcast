package Members

import (
	"container/list"
	"fmt"
	"math/rand"
)

type Member struct {
	Address string
	Name    int
}

var self *Member
var l = list.New()

func SetSelf(name int) {
	self = Find(name)
}

func GetSelf() *Member {
	return self
}

func Add(memb Member) {
	l.PushFront(memb)
}

func Count() int {
	return l.Len()
}

func Find_Random() *Member {

	var num = rand.Intn(l.Len())

	for e := l.Front(); true; e = e.Next() {
		if num == 0 {
			var item = e.Value.(Member)
			return &item
		}
		num--
	}

	return nil
}

func Find(name int) *Member {
	for e := l.Front(); e != nil; e = e.Next() {
		var item = e.Value.(Member)
		if item.Name == name {
			return &item
		}
	}

	fmt.Print("Search for ", string(name), " failed!")
	return nil
}

func ForEach(fun func(member *Member)) {
	for e := l.Front(); e != nil; e = e.Next() {
		x := e.Value.(Member)
		fun(&x)
	}
}
