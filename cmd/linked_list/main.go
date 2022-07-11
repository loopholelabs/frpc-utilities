package main

import (
	"fmt"

	"github.com/loopholelabs/common/pkg/pool"
)

func main() {
	newLinkedList := pool.NewLinkedList[string]()

	hey := newLinkedList.Insert("Hey!")
	newLinkedList.Insert("You")
	are := newLinkedList.Insert("Are")
	newLinkedList.Insert("Cool")
	smile := newLinkedList.Insert(":)")

	fmt.Println(newLinkedList.ToArray(), newLinkedList.Len())

	newLinkedList.Delete(hey)
	newLinkedList.Delete(are)
	newLinkedList.Delete(smile)

	fmt.Println(newLinkedList.ToArray(), newLinkedList.Len())
}
