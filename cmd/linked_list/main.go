package main

import (
	"fmt"

	"github.com/loopholelabs/common/pkg/pool"
)

func main() {
	newLinkedList := pool.NewDoubleLinkedList[string]()

	hey := newLinkedList.Insert("First")
	newLinkedList.Insert("Second")
	are := newLinkedList.Insert("Third")
	newLinkedList.Insert("Fourth")
	smile := newLinkedList.Insert("Fifth")

	fmt.Println(newLinkedList.ToArray(), newLinkedList.Len())

	newLinkedList.Delete(hey)
	newLinkedList.Delete(are)
	newLinkedList.Delete(smile)

	fmt.Println(newLinkedList.ToArray(), newLinkedList.Len())

	fmt.Println(newLinkedList.Shift(), newLinkedList.Len())

	fmt.Println(newLinkedList.Shift(), newLinkedList.Len())

	fmt.Println(newLinkedList.Shift(), newLinkedList.Len())

	fmt.Println(newLinkedList.Shift(), newLinkedList.Len())
}
